package snowplow

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func setEventCollectorMetadataFields(c *gin.Context, e *Event) {}

func setEventFieldsFromRequest(c *gin.Context, e *Event) {
	nuid := c.GetString("identity")
	ip, _ := c.RemoteIP()
	useragent := c.Request.UserAgent()
	e.Network_userid = &nuid
	// NOTE!! Ignore the query-param-based ip and useragent. FIXME!
	e.User_ipaddress = ip.String()
	e.Useragent = &useragent
	e.Collector_tstamp = time.Now()
}

func setEventWidthHeightFields(e *Event) {
	// Doc
	if e.Doc_size != nil {
		docDimension, _ := parseWidthHeight(*e.Doc_size)
		e.Doc_width, e.Doc_height = &docDimension.width, &docDimension.height
	}
	// Viewport
	if e.Viewport_size != nil {
		vpDimension, _ := parseWidthHeight(*e.Viewport_size)
		e.Br_viewwidth, e.Br_viewheight = &vpDimension.width, &vpDimension.height
	}
	// Screen
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
	return PageFields{
		scheme: parsedUrl.Scheme,
		host:   parsedUrl.Host,
		// port: FIXME!!!
		path:     parsedUrl.Path,
		query:    unescapedQry,
		fragment: parsedUrl.Fragment,
		medium:   medium,
		source:   source,
		term:     term,
		content:  content,
		campaign: campaign,
	}, nil
}

func setPageFields(e *Event) {
	if e.Page_url != nil {
		pageFields, err := getPageFieldsFromUrl(*e.Page_url)
		if err != nil {
			fmt.Printf("error setting page fields %s\n", err)
		}
		e.Page_urlscheme = &pageFields.scheme
		e.Page_urlhost = &pageFields.host
		e.Page_urlpath = &pageFields.path
		e.Page_urlquery = &pageFields.query
		e.Page_urlfragment = &pageFields.fragment
		e.Mkt_medium = &pageFields.medium
		e.Mkt_source = &pageFields.source
		e.Mkt_term = &pageFields.term
		e.Mkt_content = &pageFields.content
		e.Mkt_campaign = &pageFields.campaign
	}

}

func setReferrerFields(e *Event) {
	if e.Page_referrer != nil {
		pageFields, err := getPageFieldsFromUrl(*e.Page_referrer)
		if err != nil {
			fmt.Printf("error setting page fields %s\n", err)
		}
		e.Refr_urlscheme = &pageFields.scheme
		e.Refr_urlhost = &pageFields.host
		e.Refr_urlpath = &pageFields.path
		e.Refr_urlquery = &pageFields.query
		e.Refr_urlfragment = &pageFields.fragment
		e.Mkt_medium = &pageFields.medium
		e.Mkt_source = &pageFields.source
		e.Mkt_term = &pageFields.term
		e.Mkt_content = &pageFields.content
		e.Mkt_campaign = &pageFields.campaign
	}
}

func anonymizeFields(e *Event) {
	// TODO! Add field-level anonymization
}

func BuildEventFromMappedParams(c *gin.Context, e map[string]interface{}) Event {

	body, err := json.Marshal(e)
	if err != nil {
		fmt.Println(err)
	}
	shortenedEvent := ShortenedEvent{}
	err = json.Unmarshal(body, &shortenedEvent)
	if err != nil {
		fmt.Printf("error unmarshalling to shortened event %s", err)
	}
	event := Event(shortenedEvent)
	setEventCollectorMetadataFields(c, &event)
	setEventFieldsFromRequest(c, &event)
	setEventWidthHeightFields(&event)
	setPageFields(&event)
	setReferrerFields(&event)
	anonymizeFields(&event) // TODO, as necessary.
	return event
}
