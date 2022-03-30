package snowplow

import (
	b64 "encoding/base64"
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
		bVal, err := strconv.ParseBool(val)
		if err != nil {
			log.Error().Err(err).Msg("could not parse bool param")
		}
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
		return nil
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

func getSdPayload(b64EncodedPayload *string) *event.SelfDescribingPayload {
	payload, err := b64.RawStdEncoding.DecodeString(*b64EncodedPayload)
	if err != nil {
		log.Error().Err(err).Msg("could not decode b64 encoded self describing payload")
		return nil
	}
	p := event.SelfDescribingPayload{
		Schema: gjson.GetBytes(payload, "data.schema").String(),
		Data:   gjson.GetBytes(payload, "data.data").Value().(map[string]interface{}),
	}
	return &p
}

func getQueryParam(u url.URL, k string) *string {
	param := u.Query()
	val := param.Get(k)
	if val == "" {
		return nil
	} else {
		return &val
	}
}

func getPageFieldsFromUrl(rawUrl string) (PageFields, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		log.Error().Err(err).Interface("url", rawUrl).Msg("could not parse url")
		return PageFields{}, err
	}
	unescapedQry, err := url.QueryUnescape(parsedUrl.RawQuery)
	if err != nil {
		log.Error().Err(err).Interface("query", parsedUrl.RawQuery).Msg("could not unescape query params")
		return PageFields{}, err
	}
	pageFields := PageFields{
		scheme:   parsedUrl.Scheme,
		host:     parsedUrl.Host,
		path:     parsedUrl.Path,
		query:    &unescapedQry,
		medium:   getQueryParam(*parsedUrl, "utm_medium"),
		source:   getQueryParam(*parsedUrl, "utm_source"),
		term:     getQueryParam(*parsedUrl, "utm_term"),
		content:  getQueryParam(*parsedUrl, "utm_content"),
		campaign: getQueryParam(*parsedUrl, "utm_campaign"),
	}
	if parsedUrl.Fragment != "" {
		pageFields.fragment = &parsedUrl.Fragment
	}
	return pageFields, nil
}

func setTsFields(e *SnowplowEvent, params map[string]interface{}) {
	e.DvceCreatedTstamp = *getTimeParam(params, "dtm")
	e.DvceSentTstamp = *getTimeParam(params, "stm")
	e.TrueTstamp = getTimeParam(params, "ttm")
	e.CollectorTstamp = time.Now().UTC()
	e.EtlTstamp = time.Now().UTC()
	timeOnDevice := e.DvceSentTstamp.Sub(e.DvceCreatedTstamp)
	e.DerivedTstamp = e.CollectorTstamp.Add(-timeOnDevice)
}

func setMetadataFields(e *SnowplowEvent, params map[string]interface{}, conf config.Config) {
	fingerprint := uuid.New()
	eName := getStringParam(params, "e")
	e.NameTracker = *getStringParam(params, "tna")
	e.AppId = *getStringParam(params, "aid")
	e.Platform = *getStringParam(params, "p")
	e.TxnId = getStringParam(params, "tid")
	e.EventId = getStringParam(params, "eid")
	e.EventFingerprint = fingerprint
	e.OsTimezone = getStringParam(params, "tz")
	e.TrackerVersion = getStringParam(params, "tv")
	e.EtlVersion = &conf.App.Version
	e.CollectorVersion = &conf.App.Version
	e.Event = getEventType(*eName)
}

func setUserFields(c *gin.Context, e *SnowplowEvent, params map[string]interface{}) {
	ip, _ := c.RemoteIP()
	sIp := ip.String() // FIXME! Should incorporate other IP sources
	useragent := c.Request.UserAgent()
	nuid := c.GetString("identity")
	e.DomainUserid = getStringParam(params, "duid")
	e.NetworkUserid = &nuid
	e.Userid = getStringParam(params, "uid")
	e.DomainSessionIdx = getInt64Param(params, "vid")
	e.DomainSessionId = getStringParam(params, "sid")
	e.UserIpAddress = &sIp
	e.Useragent = &useragent
	e.MacAddress = getStringParam(params, "mac")
}

func setBrowserFeatures(e *SnowplowEvent, params map[string]interface{}) {
	e.BrCookies = getBoolParam(params, "cookie")
	e.BrLang = getStringParam(params, "lang")
	e.BrFeaturesPdf = getBoolParam(params, "f_pdf")
	e.BrFeaturesQuicktime = getBoolParam(params, "f_qt")
	e.BrFeaturesRealplayer = getBoolParam(params, "f_realp")
	e.BrFeaturesWindowsmedia = getBoolParam(params, "f_wma")
	e.BrFeaturesDirector = getBoolParam(params, "f_dir")
	e.BrFeaturesFlash = getBoolParam(params, "f_fla")
	e.BrFeaturesJava = getBoolParam(params, "f_java")
	e.BrFeaturesGears = getBoolParam(params, "f_gears")
	e.BrFeaturesSilverlight = getBoolParam(params, "f_ag")
	e.BrColordepth = getInt64Param(params, "cd")
	e.DocCharset = getStringParam(params, "cs")
}

func setDimensionFields(e *SnowplowEvent, params map[string]interface{}) {
	e.DocSize = getStringParam(params, "ds")
	e.ViewportSize = getStringParam(params, "vp")
	e.MonitorResolution = getStringParam(params, "res")
	if e.DocSize != nil {
		docDimension, _ := getDimensions(*e.DocSize)
		e.DocWidth, e.DocHeight = &docDimension.width, &docDimension.height
	}
	if e.ViewportSize != nil {
		vpDimension, _ := getDimensions(*e.ViewportSize)
		e.BrViewWidth, e.BrViewHeight = &vpDimension.width, &vpDimension.height
	}
	if e.MonitorResolution != nil {
		monDimension, _ := getDimensions(*e.MonitorResolution)
		e.DvceScreenWidth, e.DvceScreenHeight = &monDimension.width, &monDimension.height
	}
}

func setPageFields(e *SnowplowEvent, params map[string]interface{}) {
	e.PageUrl = getStringParam(params, "url")
	e.PageTitle = getStringParam(params, "page")
	if e.PageUrl != nil {
		pageFields, err := getPageFieldsFromUrl(*e.PageUrl)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error getting page fields from url")
		}
		e.PageUrlScheme = &pageFields.scheme
		e.PageUrlHost = &pageFields.host
		e.PageUrlPath = &pageFields.path
		e.PageUrlQuery = pageFields.query
		e.PageUrlFragment = pageFields.fragment
		e.MktMedium = pageFields.medium
		e.MktSource = pageFields.source
		e.MktTerm = pageFields.term
		e.MktContent = pageFields.content
		e.MktCampaign = pageFields.campaign
	}
}

func setReferrerFields(e *SnowplowEvent, params map[string]interface{}) {
	e.PageReferrer = getStringParam(params, "refr")
	if e.PageReferrer != nil {
		pageFields, err := getPageFieldsFromUrl(*e.PageReferrer)
		if err != nil {
			log.Error().Err(err).Msg("error setting page fields")
		}
		e.RefrUrlScheme = &pageFields.scheme
		e.RefrUrlHost = &pageFields.host
		e.RefrUrlPath = &pageFields.path
		e.RefrUrlQuery = pageFields.query
		e.RefrUrlFragment = pageFields.fragment
		e.RefrMedium = pageFields.medium
		e.RefrSource = pageFields.source
		e.RefrTerm = pageFields.term
		e.RefrContent = pageFields.content
		e.RefrCampaign = pageFields.campaign
	}
}

func anonymizeFields(e *SnowplowEvent, conf config.Snowplow) {
	if conf.Anonymize.Ip && e.UserIpAddress != nil {
		hashedIp := util.Md5(*e.UserIpAddress)
		e.UserIpAddress = &hashedIp
	}
	if conf.Anonymize.UserId && e.Userid != nil {
		hashedUserId := util.Md5(*e.Userid)
		e.Userid = &hashedUserId
	}
}

func setPagePingFields(e *SnowplowEvent, params map[string]interface{}) {
	e.PpXOffsetMin = getInt64Param(params, "pp_mix")
	e.PpXOffsetMax = getInt64Param(params, "pp_max")
	e.PpYOffsetMin = getInt64Param(params, "pp_miy")
	e.PpYOffsetMax = getInt64Param(params, "pp_may")
}

func setStructFields(e *SnowplowEvent, params map[string]interface{}) {
	e.SeCategory = getStringParam(params, "se_ca")
	e.SeAction = getStringParam(params, "se_ac")
	e.SeLabel = getStringParam(params, "se_la")
	e.SeProperty = getStringParam(params, "se_pr")
	e.SeValue = getFloat64Param(params, "se_va")
}

func setTransactionFields(e *SnowplowEvent, params map[string]interface{}) {
	e.TrOrderId = getStringParam(params, "tr_id")
	e.TrAffiliation = getStringParam(params, "tr_af")
	e.TrTotal = getFloat64Param(params, "tr_tt")
	e.TrTax = getFloat64Param(params, "tr_tx")
	e.TrShipping = getFloat64Param(params, "tr_sh")
	e.TrCity = getStringParam(params, "tr_ci")
	e.TrState = getStringParam(params, "tr_st")
	e.TrCountry = getStringParam(params, "tr_co")
	e.TrCurrency = getStringParam(params, "tr_cu")
}

func setTransactionItemFields(e *SnowplowEvent, params map[string]interface{}) {
	e.TiOrderId = getStringParam(params, "ti_id")
	e.TiSku = getStringParam(params, "ti_sk")
	e.TiName = getStringParam(params, "ti_nm")
	e.TiCategory = getStringParam(params, "ti_ca")
	e.TiPrice = getFloat64Param(params, "ti_pr")
	e.TiQuantity = getInt64Param(params, "ti_qu")
	e.TiCurrency = getStringParam(params, "ti_cu")
}

func setContexts(e *SnowplowEvent, params map[string]interface{}) {
	b64encodedContexts := getStringParam(params, "cx")
	e.Contexts = getContexts(b64encodedContexts)
}

func setSelfDescribingFields(e *SnowplowEvent, params map[string]interface{}) {
	b64EncodedPayload := getStringParam(params, "ue_px")
	e.SelfDescribingEvent = getSdPayload(b64EncodedPayload)
}

func setEventMetadataFields(e *SnowplowEvent, schema []byte) {
	schemaContents := gjson.ParseBytes(schema)
	vendor := schemaContents.Get("self.vendor").String()
	name := schemaContents.Get("self.name").String()
	format := schemaContents.Get("self.format").String()
	version := schemaContents.Get("self.version").String()
	if vendor != "" {
		e.EventVendor = &vendor
	}
	if name != "" {
		e.EventName = &name
	}
	if format != "" {
		e.EventFormat = &format
	}
	if version != "" {
		e.EventVersion = &version
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
	return event
}
