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

func setPage(e *SnowplowEvent, params map[string]interface{}) {
	page := getStringParam(params, "page")
}

func setReferrer(e *SnowplowEvent, params map[string]interface{}) {
	referrer := getStringParam(params, "refr")
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
	evnt := PagePingEvent{
		PpXOffsetMin: getInt64Param(params, "pp_mix"),
		PpXOffsetMax: getInt64Param(params, "pp_max"),
		PpYOffsetMin: getInt64Param(params, "pp_miy"),
		PpYOffsetMax: getInt64Param(params, "pp_may"),
	}
	sde := evnt.toSelfDescribing()
	e.SelfDescribingEvent = &sde
}

func setStructFields(e *SnowplowEvent, params map[string]interface{}) {
	evnt := StructEvent{
		SeCategory: getStringParam(params, "se_ca"),
		SeAction:   getStringParam(params, "se_ac"),
		SeLabel:    getStringParam(params, "se_la"),
		SeProperty: getStringParam(params, "se_pr"),
		SeValue:    getFloat64Param(params, "se_va"),
	}
	sde := evnt.toSelfDescribing()
	e.SelfDescribingEvent = &sde
}

func setTransactionFields(e *SnowplowEvent, params map[string]interface{}) {
	evnt := TransactionEvent{
		TrOrderId:     getStringParam(params, "tr_id"),
		TrAffiliation: getStringParam(params, "tr_af"),
		TrTotal:       getFloat64Param(params, "tr_tt"),
		TrTax:         getFloat64Param(params, "tr_tx"),
		TrShipping:    getFloat64Param(params, "tr_sh"),
		TrCity:        getStringParam(params, "tr_ci"),
		TrState:       getStringParam(params, "tr_st"),
		TrCountry:     getStringParam(params, "tr_co"),
		TrCurrency:    getStringParam(params, "tr_cu"),
	}
	sde := evnt.toSelfDescribing()
	e.SelfDescribingEvent = &sde
}

func setTransactionItemFields(e *SnowplowEvent, params map[string]interface{}) {
	evnt := TransactionItemEvent{
		TiOrderId:  getStringParam(params, "ti_id"),
		TiSku:      getStringParam(params, "ti_sk"),
		TiName:     getStringParam(params, "ti_nm"),
		TiCategory: getStringParam(params, "ti_ca"),
		TiPrice:    getFloat64Param(params, "ti_pr"),
		TiQuantity: getInt64Param(params, "ti_qu"),
		TiCurrency: getStringParam(params, "ti_cu"),
	}
	sde := evnt.toSelfDescribing()
	e.SelfDescribingEvent = &sde
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
