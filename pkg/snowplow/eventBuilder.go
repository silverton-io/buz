package snowplow

import (
	b64 "encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/tidwall/gjson"
)

func getStringParam(params map[string]interface{}, k string) *string {
	v := params[k]
	if v == nil {
		return nil
	} else {
		val := v.(string)
		return &val
	}
}

func getInt64Param(params map[string]interface{}, k string) *int64 {
	v := params[k]
	if v == nil {
		return nil
	} else {
		val := v.(string)
		iVal, _ := strconv.ParseInt(val, 10, 64)
		return &iVal
	}
}

func getFloat64Param(params map[string]interface{}, k string) *float64 {
	v := params[k]
	if v == nil {
		return nil
	} else {
		val := v.(string)
		fVal, _ := strconv.ParseFloat(val, 64)
		return &fVal
	}
}

func getTimeParam(params map[string]interface{}, k string) *time.Time {
	v := params[k]
	if v == nil {
		return nil
	} else {
		val := v.(string)
		iVal, _ := strconv.ParseInt(val, 10, 64)
		t := time.UnixMilli(iVal)
		return &t
	}
}

func getBoolParam(params map[string]interface{}, k string) *bool {
	v := params[k]
	if v == nil {
		return nil
	} else {
		val := v.(string)
		bVal, _ := strconv.ParseBool(val)
		return &bVal
	}
}

func getDimensions(dimensionString string) (Dimension, error) {
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

func getContexts(b64encodedContexts *string) *[]event.SelfDescribingContext {
	var contexts []event.SelfDescribingContext
	payload, err := b64.RawStdEncoding.DecodeString(*b64encodedContexts)
	if err != nil {
		log.Error().Err(err).Msg("could not decode b64 encoded contexts")
	}
	contextPayload := gjson.ParseBytes(payload)
	for _, pl := range contextPayload.Get("data").Array() {
		context := event.SelfDescribingContext{
			Schema: pl.Get("schema").String(),
			Data:   pl.Get("data").Value().(map[string]interface{}),
		}
		contexts = append(contexts, context)
	}
	return &contexts
}

func getSdEvent(b64SelfDescribingEvent *string) *event.SelfDescribingEvent {
	event := event.SelfDescribingEvent{}
	payload, err := b64.RawStdEncoding.DecodeString(*b64SelfDescribingEvent)
	if err != nil {
		log.Error().Err(err).Msg("could not decode b64 encoded self describing event")
	}
	schema := gjson.GetBytes(payload, "data.schema").Value().(string)
	data := gjson.GetBytes(payload, "data.data").Value().(map[string]interface{})
	event.Payload.Schema = schema
	event.Payload.Data = data
	return &event
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

func setTsFields(e *SnowplowEvent, params map[string]interface{}) {
	e.Dvce_created_tstamp = *getTimeParam(params, "dtm")
	e.Dvce_sent_tstamp = *getTimeParam(params, "stm")
	e.True_tstamp = getTimeParam(params, "ttm")
	e.Collector_tstamp = time.Now().UTC()
	e.Etl_tstamp = time.Now().UTC()
	timeOnDevice := e.Dvce_sent_tstamp.Sub(e.Dvce_created_tstamp)
	e.Derived_tstamp = e.Collector_tstamp.Add(-timeOnDevice)
}

func setMetadataFields(e *SnowplowEvent, params map[string]interface{}, conf config.Config) {
	fingerprint := uuid.New()
	eName := getStringParam(params, "e")
	e.Name_tracker = *getStringParam(params, "tna")
	e.App_id = *getStringParam(params, "aid")
	e.Platform = *getStringParam(params, "p")
	e.Txn_id = getStringParam(params, "tid")
	e.Event_id = getStringParam(params, "eid")
	e.Event_fingerprint = fingerprint
	e.Os_timezone = getStringParam(params, "tz")
	e.Tracker_version = getStringParam(params, "tv")
	e.Etl_version = &conf.Version
	e.Collector_version = &conf.Version
	e.Event = getEventType(*eName)
}

func setUserFields(c *gin.Context, e *SnowplowEvent, params map[string]interface{}) {
	ip, _ := c.RemoteIP()
	sIp := ip.String() // FIXME! Should incorporate other IP sources
	useragent := c.Request.UserAgent()
	nuid := c.GetString("identity")
	e.Domain_userid = getStringParam(params, "duid")
	fmt.Println(nuid)
	e.Network_userid = &nuid
	e.Userid = getStringParam(params, "uid")
	e.Domain_sessionidx = getInt64Param(params, "vid")
	e.Domain_sessionid = getStringParam(params, "sid")
	e.User_ipaddress = &sIp
	e.Useragent = &useragent
	e.Mac_address = getStringParam(params, "mac")
}

func setBrowserFeatures(e *SnowplowEvent, params map[string]interface{}) {
	e.Br_cookies = getBoolParam(params, "cookie")
	e.Br_lang = getStringParam(params, "lang")
	e.Br_features_pdf = getBoolParam(params, "f_pdf")
	e.Br_features_quicktime = getBoolParam(params, "f_qt")
	e.Br_features_realplayer = getBoolParam(params, "f_realp")
	e.Br_features_windowsmedia = getBoolParam(params, "f_wma")
	e.Br_features_director = getBoolParam(params, "f_dir")
	e.Br_features_flash = getBoolParam(params, "f_fla")
	e.Br_features_java = getBoolParam(params, "f_java")
	e.Br_features_gears = getBoolParam(params, "f_gears")
	e.Br_features_silverlight = getBoolParam(params, "f_ag")
	e.Br_colordepth = getInt64Param(params, "cd")
	e.Doc_charset = getStringParam(params, "cs")
}

func setDimensionFields(e *SnowplowEvent, params map[string]interface{}) {
	e.Doc_size = getStringParam(params, "ds")
	e.Viewport_size = getStringParam(params, "vp")
	e.Monitor_resolution = getStringParam(params, "res")
	if e.Doc_size != nil {
		docDimension, _ := getDimensions(*e.Doc_size)
		e.Doc_width, e.Doc_height = &docDimension.width, &docDimension.height
	}
	if e.Viewport_size != nil {
		vpDimension, _ := getDimensions(*e.Viewport_size)
		e.Br_viewwidth, e.Br_viewheight = &vpDimension.width, &vpDimension.height
	}
	if e.Monitor_resolution != nil {
		monDimension, _ := getDimensions(*e.Monitor_resolution)
		e.Dvce_screenwidth, e.Dvce_screenheight = &monDimension.width, &monDimension.height
	}
}

func setPageFields(e *SnowplowEvent, params map[string]interface{}) {
	e.Page_url = getStringParam(params, "url")
	e.Page_title = getStringParam(params, "page")
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

func setReferrerFields(e *SnowplowEvent, params map[string]interface{}) {
	e.Page_referrer = getStringParam(params, "refr")
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

func setPagePingFields(e *SnowplowEvent, params map[string]interface{}) {
	e.Pp_xoffset_min = getInt64Param(params, "pp_mix")
	e.Pp_xoffset_max = getInt64Param(params, "pp_max")
	e.Pp_yoffset_min = getInt64Param(params, "pp_miy")
	e.Pp_yoffset_max = getInt64Param(params, "pp_may")
}

func setStructFields(e *SnowplowEvent, params map[string]interface{}) {
	e.Se_category = getStringParam(params, "se_ca")
	e.Se_action = getStringParam(params, "se_ac")
	e.Se_label = getStringParam(params, "se_la")
	e.Se_property = getStringParam(params, "se_pr")
	e.Se_value = getFloat64Param(params, "se_va")
}

func setTransactionFields(e *SnowplowEvent, params map[string]interface{}) {
	e.Tr_orderid = getStringParam(params, "tr_id")
	e.Tr_affiliation = getStringParam(params, "tr_af")
	e.Tr_total = getFloat64Param(params, "tr_tt")
	e.Tr_tax = getFloat64Param(params, "tr_tx")
	e.Tr_shipping = getFloat64Param(params, "tr_sh")
	e.Tr_city = getStringParam(params, "tr_ci")
	e.Tr_state = getStringParam(params, "tr_st")
	e.Tr_country = getStringParam(params, "tr_co")
	e.Tr_currency = getStringParam(params, "tr_cu")
}

func setTransactionItemFields(e *SnowplowEvent, params map[string]interface{}) {
	e.Ti_orderid = getStringParam(params, "ti_id")
	e.Ti_sku = getStringParam(params, "ti_sk")
	e.Ti_name = getStringParam(params, "ti_nm")
	e.Ti_category = getStringParam(params, "ti_ca")
	e.Ti_price = getFloat64Param(params, "ti_pr")
	e.Ti_quantity = getInt64Param(params, "ti_qu")
	e.Ti_currency = getStringParam(params, "ti_cu")
}

func setContexts(e *SnowplowEvent, params map[string]interface{}) {
	b64encodedContexts := getStringParam(params, "cx")
	e.Contexts = getContexts(b64encodedContexts)
}

func setSelfDescribingFields(e *SnowplowEvent, params map[string]interface{}) {
	b64encodedEvent := getStringParam(params, "ue_px")
	e.Self_describing_event = getSdEvent(b64encodedEvent)
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

func BuildEventFromMappedParams(c *gin.Context, params map[string]interface{}, conf config.Config) SnowplowEvent {
	event := SnowplowEvent{}
	setTsFields(&event, params)
	setMetadataFields(&event, params, conf)
	setUserFields(c, &event, params)
	setBrowserFeatures(&event, params)
	setDimensionFields(&event, params)
	setPageFields(&event, params)
	setReferrerFields(&event, params)
	anonymizeFields(&event, conf.Snowplow)
	setContexts(&event, params)
	if event.Event == PAGE_PING {
		setPagePingFields(&event, params)
	}
	if event.Event == STRUCT_EVENT {
		setStructFields(&event, params)
	}
	if event.Event == TRANSACTION {
		setTransactionFields(&event, params)
	}
	if event.Event == TRANSACTION_ITEM {
		setTransactionItemFields(&event, params)
	}
	if event.Event == SELF_DESCRIBING_EVENT {
		setSelfDescribingFields(&event, params)
	}
	util.Pprint(event)
	return event
}
