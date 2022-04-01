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

func getPageParam(params map[string]interface{}, k string) (Page, error) {
	p := getStringParam(params, k)
	if p != nil {
		parsedUrl, err := url.Parse(*p)
		if err != nil {
			log.Error().Err(err).Interface("url", *p).Msg("could not parse url")
			return Page{}, err
		}
		unescapedQry, err := url.QueryUnescape(parsedUrl.RawQuery)
		if err != nil {
			log.Error().Err(err).Interface("query", parsedUrl.RawQuery).Msg("could not unescape query params")
		}
		frag := parsedUrl.Fragment
		page := Page{
			Url:      *p,
			Scheme:   parsedUrl.Scheme,
			Host:     parsedUrl.Host,
			Port:     parsedUrl.Port(),
			Path:     parsedUrl.Path,
			Query:    &unescapedQry,
			Fragment: &frag,
			Medium:   getQueryParam(*parsedUrl, "utm_medium"),
			Source:   getQueryParam(*parsedUrl, "utm_source"),
			Term:     getQueryParam(*parsedUrl, "utm_term"),
			Content:  getQueryParam(*parsedUrl, "utm_content"),
			Campaign: getQueryParam(*parsedUrl, "utm_campaign"),
		}
		return page, nil
	} else {
		return Page{}, nil
	}
}

func setTstamps(e *SnowplowEvent, params map[string]interface{}) {
	ts := Tstamp{
		DvceCreatedTstamp: *getTimeParam(params, "dtm"),
		DvceSentTstamp:    *getTimeParam(params, "stm"),
		TrueTstamp:        getTimeParam(params, "ttm"),
		CollectorTstamp:   time.Now().UTC(),
		EtlTstamp:         time.Now().UTC(),
	}
	timeOnDevice := ts.DvceSentTstamp.Sub(ts.DvceCreatedTstamp)
	ts.DerivedTstamp = ts.CollectorTstamp.Add(-timeOnDevice)
	e.Tstamp = ts
}

func setPlatformMetadata(e *SnowplowEvent, params map[string]interface{}, conf config.Config) {
	pm := PlatformMetadata{
		NameTracker:      *getStringParam(params, "tna"),
		TrackerVersion:   getStringParam(params, "tv"),
		CollectorVersion: &conf.App.Version,
		EtlVersion:       &conf.App.Version,
	}
	e.PlatformMetadata = pm
}

func setEvent(e *SnowplowEvent, params map[string]interface{}) {
	eName := getStringParam(params, "e")
	fingerprint := uuid.New()
	evnt := Event{
		AppId:            *getStringParam(params, "aid"),
		Platform:         *getStringParam(params, "p"),
		Event:            getEventType(*eName),
		TxnId:            getStringParam(params, "tid"),
		EventId:          getStringParam(params, "eid"),
		EventFingerprint: fingerprint,
	}
	e.Event = evnt
}

func setUser(c *gin.Context, e *SnowplowEvent, params map[string]interface{}) {
	nuid := c.GetString("identity")
	ip, _ := c.RemoteIP()
	sIp := ip.String() // FIXME! Should incorporate other IP sources
	user := User{
		DomainUserid:  getStringParam(params, "duid"),
		NetworkUserid: &nuid,
		Userid:        getStringParam(params, "uid"),
		UserIpAddress: &sIp,
	}
	e.User = user
}

func setSession(e *SnowplowEvent, params map[string]interface{}) {
	s := Session{
		DomainSessionIdx: getInt64Param(params, "vid"),
		DomainSessionId:  getStringParam(params, "sid"),
	}
	e.Session = s
}

func setPage(e *SnowplowEvent, params map[string]interface{}) {
	page, _ := getPageParam(params, "url")
	title := getStringParam(params, "page")
	page.Title = title
	e.Page = page
}

func setReferrer(e *SnowplowEvent, params map[string]interface{}) {
	referrer, _ := getPageParam(params, "refr")
	e.Referrer = referrer
}

func setDevice(c *gin.Context, e *SnowplowEvent, params map[string]interface{}) {
	useragent := c.Request.UserAgent()
	d := Device{
		Useragent:  &useragent,
		MacAddress: getStringParam(params, "mac"),
		OsTimezone: getStringParam(params, "tz"),
	}
	e.Device = d
}

func setBrowser(e *SnowplowEvent, params map[string]interface{}) {
	b := Browser{
		BrCookies:              getBoolParam(params, "cookie"),
		BrLang:                 getStringParam(params, "lang"),
		BrFeaturesPdf:          getBoolParam(params, "f_pdf"),
		BrFeaturesQuicktime:    getBoolParam(params, "f_qt"),
		BrFeaturesRealplayer:   getBoolParam(params, "f_realp"),
		BrFeaturesWindowsmedia: getBoolParam(params, "f_wma"),
		BrFeaturesDirector:     getBoolParam(params, "f_dir"),
		BrFeaturesFlash:        getBoolParam(params, "f_fla"),
		BrFeaturesJava:         getBoolParam(params, "f_java"),
		BrFeaturesGears:        getBoolParam(params, "f_gears"),
		BrFeaturesSilverlight:  getBoolParam(params, "f_ag"),
		BrColordepth:           getInt64Param(params, "cd"),
	}
	e.Browser = b
}

func setScreen(e *SnowplowEvent, params map[string]interface{}) {
	s := Screen{
		DocCharset:        getStringParam(params, "cs"),
		ViewportSize:      getStringParam(params, "vp"),
		DocSize:           getStringParam(params, "ds"),
		MonitorResolution: getStringParam(params, "res"),
	}
	if s.DocSize != nil {
		docDimension, _ := getDimensions(*s.DocSize)
		s.DocWidth, s.DocHeight = &docDimension.width, &docDimension.height
	}
	if s.ViewportSize != nil {
		vpDimension, _ := getDimensions(*s.ViewportSize)
		s.BrViewWidth, s.BrViewHeight = &vpDimension.width, &vpDimension.height
	}
	if s.MonitorResolution != nil {
		monDimension, _ := getDimensions(*s.MonitorResolution)
		s.DvceScreenWidth, s.DvceScreenHeight = &monDimension.width, &monDimension.height
	}
	e.Screen = s
}

func setPagePing(e *SnowplowEvent, params map[string]interface{}) {
	evnt := PagePingEvent{
		PpXOffsetMin: getInt64Param(params, "pp_mix"),
		PpXOffsetMax: getInt64Param(params, "pp_max"),
		PpYOffsetMin: getInt64Param(params, "pp_miy"),
		PpYOffsetMax: getInt64Param(params, "pp_may"),
	}
	sde := evnt.toSelfDescribing()
	e.SelfDescribingEvent = &sde
}

func setStruct(e *SnowplowEvent, params map[string]interface{}) {
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

func setTransaction(e *SnowplowEvent, params map[string]interface{}) {
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

func setTransactionItem(e *SnowplowEvent, params map[string]interface{}) {
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

func setSelfDescribing(e *SnowplowEvent, params map[string]interface{}) {
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

func BuildEventFromMappedParams(c *gin.Context, params map[string]interface{}, conf config.Config) SnowplowEvent {
	event := SnowplowEvent{}
	setTstamps(&event, params)
	setPlatformMetadata(&event, params, conf)
	setEvent(&event, params)
	setUser(c, &event, params)
	setSession(&event, params)
	setPage(&event, params)
	setReferrer(&event, params)
	setDevice(c, &event, params)
	setBrowser(&event, params)
	setScreen(&event, params)
	setContexts(&event, params)
	switch event.Event.Event {
	case PAGE_PING:
		setPagePing(&event, params)
	case STRUCT_EVENT:
		setStruct(&event, params)
	case TRANSACTION:
		setTransaction(&event, params)
	case TRANSACTION_ITEM:
		setTransactionItem(&event, params)
	case SELF_DESCRIBING_EVENT:
		setSelfDescribing(&event, params)
	}
	anonymizeFields(&event, conf.Snowplow)
	return event
}
