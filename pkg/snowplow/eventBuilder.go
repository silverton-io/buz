package snowplow

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/tidwall/gjson"
)

func parseWidthHeight(dimensionString string) (Dimension, error) {
	dim := strings.Split(dimensionString, "x")
	wS, hS := dim[0], dim[1]
	width, err := strconv.Atoi(wS)
	height, err := strconv.Atoi(hS)
	if err != nil {
		return Dimension{}, err
	}
	return Dimension{
		width:  width,
		height: height,
	}, nil
}

func setEventMetadataFields(e *SnowplowEvent, schema []byte) {
	schemaContents := gjson.ParseBytes(schema)
	vendor := schemaContents.Get("self.vendor").String()
	name := schemaContents.Get("self.name").String()
	format := schemaContents.Get("self.format").String()
	version := schemaContents.Get("self.version").String()
	if vendor != "" {
		e.Event_vendor = &vendor
	}
	if name != "" {
		e.Event_name = &name
	}
	if format != "" {
		e.Event_format = &format
	}
	if version != "" {
		e.Event_version = &version
	}
}

func setEventFieldsFromRequest(c *gin.Context, e *SnowplowEvent, conf *config.App) {
	nuid := c.GetString("identity")
	ip, _ := c.RemoteIP()
	sIp := ip.String() // FIXME! Should incorporate other IP sources
	useragent := c.Request.UserAgent()
	e.Network_userid = &nuid
	e.User_ipaddress = &sIp
	e.Useragent = &useragent
	e.Collector_tstamp = time.Now()
	e.Etl_tstamp = time.Now()
	timeOnDevice := e.Dvce_sent_tstamp.Time.Sub(e.Dvce_created_tstamp.Time)
	e.Derived_tstamp = e.Collector_tstamp.Add(-timeOnDevice)
	e.Collector_version = &conf.Version
	e.Etl_version = &conf.Version
}

func setEventWidthHeightFields(e *SnowplowEvent) {
	if e.Doc_size != nil {
		docDimension, _ := parseWidthHeight(*e.Doc_size)
		e.Doc_width, e.Doc_height = &docDimension.width, &docDimension.height
	}
	if e.Viewport_size != nil {
		vpDimension, _ := parseWidthHeight(*e.Viewport_size)
		e.Br_viewwidth, e.Br_viewheight = &vpDimension.width, &vpDimension.height
	}
	if e.Monitor_resolution != nil {
		monDimension, _ := parseWidthHeight(*e.Monitor_resolution)
		e.Dvce_screenwidth, e.Dvce_screenheight = &monDimension.width, &monDimension.height
	}
}

func getPageFieldsFromUrl(rawUrl string) (PageFields, error) {
	parsedUrl, err := url.Parse(rawUrl)
	queryParams := parsedUrl.Query()
	unescapedQry, err := url.QueryUnescape(parsedUrl.RawQuery)
	if err != nil {
		return PageFields{}, err
	}
	medium := queryParams.Get("utm_medium")
	source := queryParams.Get("utm_source")
	term := queryParams.Get("utm_term")
	content := queryParams.Get("utm_content")
	campaign := queryParams.Get("utm_campaign")
	pageFields := PageFields{
		scheme: parsedUrl.Scheme,
		host:   parsedUrl.Host,
		path:   parsedUrl.Path,
	}
	if unescapedQry != "" { // FIXME! Has to be a better way
		pageFields.query = &unescapedQry
	}
	if medium != "" {
		pageFields.medium = &medium
	}
	if source != "" {
		pageFields.source = &source
	}
	if term != "" {
		pageFields.term = &term
	}
	if content != "" {
		pageFields.content = &content
	}
	if campaign != "" {
		pageFields.campaign = &campaign
	}
	if parsedUrl.Fragment != "" {
		pageFields.fragment = &parsedUrl.Fragment
	}
	return pageFields, nil
}

func setPageFields(e *SnowplowEvent) {
	if e.Page_url != nil {
		pageFields, err := getPageFieldsFromUrl(*e.Page_url)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error getting page fields from url")
		}
		e.Page_urlscheme = &pageFields.scheme // FIXME! Has to be a better way
		e.Page_urlhost = &pageFields.host
		e.Page_urlpath = &pageFields.path
		e.Page_urlquery = pageFields.query
		e.Page_urlfragment = pageFields.fragment
		e.Mkt_medium = pageFields.medium
		e.Mkt_source = pageFields.source
		e.Mkt_term = pageFields.term
		e.Mkt_content = pageFields.content
		e.Mkt_campaign = pageFields.campaign
	}
}

func setReferrerFields(e *SnowplowEvent) {
	if e.Page_referrer != nil {
		pageFields, err := getPageFieldsFromUrl(*e.Page_referrer)
		if err != nil {
			log.Error().Err(err).Msg("error setting page fields")
		}
		e.Refr_urlscheme = &pageFields.scheme // FIXME! Has to be a better way
		e.Refr_urlhost = &pageFields.host
		e.Refr_urlpath = &pageFields.path
		e.Refr_urlquery = pageFields.query
		e.Refr_urlfragment = pageFields.fragment
		e.Refr_medium = pageFields.medium
		e.Refr_source = pageFields.source
		e.Refr_term = pageFields.term
		e.Refr_content = pageFields.content
		e.Refr_campaign = pageFields.campaign
	}
}

func anonymizeFields(e *SnowplowEvent, conf config.Snowplow) {
	if conf.Anonymize.Ip && e.User_ipaddress != nil {
		hashedIp := util.Md5(*e.User_ipaddress)
		e.User_ipaddress = &hashedIp
	}
	if conf.Anonymize.UserId && e.Userid != nil {
		hashedUserId := util.Md5(*e.Userid)
		e.Userid = &hashedUserId
	}
}

func BuildEventFromMappedParams(c *gin.Context, params map[string]interface{}, conf config.Config) SnowplowEvent {
	body, err := json.Marshal(params)
	if err != nil {
		log.Error().Err(err).Msg("error when marshaling params")
	}
	shortenedEvent := ShortenedSnowplowEvent{}
	err = json.Unmarshal(body, &shortenedEvent)
	if err != nil {
		log.Error().Err(err).Msg("error when unmarshaling to shortened event")
	}
	event := SnowplowEvent(shortenedEvent)
	setEventFieldsFromRequest(c, &event, &conf.App)
	setEventWidthHeightFields(&event)
	setPageFields(&event)
	setReferrerFields(&event)
	anonymizeFields(&event, conf.Snowplow)
	return event
}
