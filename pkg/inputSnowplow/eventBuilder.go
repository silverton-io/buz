// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputsnowplow

import (
	b64 "encoding/base64"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/event"
	"github.com/silverton-io/buz/pkg/util"
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
			log.Error().Err(err).Msg("🔴 could not parse bool param")
		}
		return &bVal
	}
}

func getDimensions(dimensionString string) (Dimension, error) {
	dim := strings.Split(dimensionString, "x")
	wS, hS := dim[0], dim[1]
	width, err := strconv.Atoi(wS)
	if err != nil {
		return Dimension{}, err
	}
	height, err := strconv.Atoi(hS)
	if err != nil {
		return Dimension{}, err
	}
	return Dimension{
		width:  width,
		height: height,
	}, nil
}

func getContexts(b64encodedContexts *string) *map[string]interface{} {
	var contexts = make(map[string]interface{})
	payload, err := b64.RawStdEncoding.DecodeString(*b64encodedContexts)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not decode b64 encoded contexts")
		return nil
	}
	contextPayload := gjson.ParseBytes(payload)
	for _, pl := range contextPayload.Get("data").Array() {
		schema := pl.Get("schema").String()
		data := pl.Get("data").Value().(map[string]interface{})
		contexts[schema] = data
	}
	return &contexts
}

func getSdPayload(b64EncodedPayload *string) *event.SelfDescribingPayload {
	if b64EncodedPayload != nil {
		payload, err := b64.RawStdEncoding.DecodeString(*b64EncodedPayload)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not decode b64 encoded self describing payload")
			return nil
		}
		schema := gjson.GetBytes(payload, "data.schema").String()
		p := event.SelfDescribingPayload{
			Schema: schema,
			Data:   gjson.GetBytes(payload, "data.data").Value().(map[string]interface{}),
		}
		return &p
	} else {
		return nil
	}
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

func getPageFromParam(params map[string]interface{}, k string) (Page, error) {
	p := getStringParam(params, k)
	if p != nil {
		parsedUrl, err := url.Parse(*p)
		if err != nil {
			log.Error().Err(err).Interface("url", *p).Msg("🔴 could not parse url")
			return Page{}, err
		}
		qParams := util.QueryToMap(parsedUrl.Query())
		frag := parsedUrl.Fragment
		page := Page{
			Url:      *p,
			Scheme:   parsedUrl.Scheme,
			Host:     parsedUrl.Host,
			Port:     parsedUrl.Port(),
			Path:     parsedUrl.Path,
			Query:    qParams,
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
	now := time.Now().UTC()
	e.DvceCreatedTstamp = getTimeParam(params, "dtm")
	e.DvceSentTstamp = getTimeParam(params, "stm")
	e.TrueTstamp = getTimeParam(params, "ttm")
	e.CollectorTstamp = time.Now().UTC()
	e.EtlTstamp = &now
	if e.DvceCreatedTstamp != nil {
		timeOnDevice := e.DvceSentTstamp.Sub(*e.DvceCreatedTstamp)
		e.DerivedTstamp = e.CollectorTstamp.Add(-timeOnDevice)
	}
}

func setPlatformMetadata(e *SnowplowEvent, params map[string]interface{}, conf config.Config) {
	e.NameTracker = getStringParam(params, "tna")
	e.TrackerVersion = getStringParam(params, "tv")
	e.CollectorVersion = &conf.App.Version
	e.EtlVersion = &conf.App.Version
}

func setEvent(e *SnowplowEvent, params map[string]interface{}) {
	eName := getStringParam(params, "e")
	fingerprint := uuid.New()
	e.AppId = getStringParam(params, "aid")
	e.Platform = *getStringParam(params, "p")
	e.Event = getEventType(*eName)
	e.TxnId = getStringParam(params, "tid")
	e.EventId = getStringParam(params, "eid")
	e.EventFingerprint = fingerprint
}

func setUser(c *gin.Context, e *SnowplowEvent, params map[string]interface{}) {
	e.DomainUserid = getStringParam(params, "duid")
	e.Userid = getStringParam(params, "uid")
}

func setSession(e *SnowplowEvent, params map[string]interface{}) {
	vid := getInt64Param(params, "vid")
	sid := getStringParam(params, "sid")
	e.DomainSessionIdx = vid
	e.DomainSessionId = sid
}

func setPage(e *SnowplowEvent, params map[string]interface{}) {
	page, _ := getPageFromParam(params, "url")
	title := getStringParam(params, "page")
	page.Title = title
	e.PageUrl = &page.Url
	e.PageTitle = page.Title
	e.PageUrlScheme = &page.Scheme
	e.PageUrlHost = &page.Host
	e.PageUrlPort = &page.Port
	e.PageUrlPath = &page.Path
	e.PageUrlQuery = &page.Query
	e.PageUrlFragment = page.Fragment
	e.MktCampaign = page.Campaign
	e.MktContent = page.Content
	e.MktMedium = page.Medium
	e.MktSource = page.Source
	e.MktTerm = page.Term
}

func setReferrer(e *SnowplowEvent, params map[string]interface{}) {
	refr, _ := getPageFromParam(params, "refr")
	e.PageReferrer = &refr.Url
	e.RefrUrlScheme = &refr.Scheme
	e.RefrUrlHost = &refr.Host
	e.RefrUrlPort = &refr.Port
	e.RefrUrlPath = &refr.Path
	e.RefrUrlQuery = &refr.Query
	e.RefrUrlFragment = refr.Fragment
	e.RefrCampaign = refr.Campaign
	e.RefrContent = refr.Content
	e.RefrMedium = refr.Medium
	e.RefrSource = refr.Source
	e.RefrTerm = refr.Term
}

func setDevice(c *gin.Context, e *SnowplowEvent, params map[string]interface{}) {
	useragent := c.Request.UserAgent()
	e.Useragent = &useragent
	e.MacAddress = getStringParam(params, "mac")
	e.OsTimezone = getStringParam(params, "tz")
}

func setBrowser(e *SnowplowEvent, params map[string]interface{}) {
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
}

func setScreen(e *SnowplowEvent, params map[string]interface{}) {
	e.DocCharset = getStringParam(params, "cs")
	e.ViewportSize = getStringParam(params, "vp")
	e.DvceScreenResolution = getStringParam(params, "res")
	e.DocSize = getStringParam(params, "ds")
	if e.DocSize != nil {
		docDimension, _ := getDimensions(*e.DocSize)
		e.DocWidth, e.DocHeight = &docDimension.width, &docDimension.height
	}
	if e.ViewportSize != nil {
		vpDimension, _ := getDimensions(*e.ViewportSize)
		e.BrViewWidth, e.BrViewHeight = &vpDimension.width, &vpDimension.height
	}
	if e.DvceScreenResolution != nil {
		monDimension, _ := getDimensions(*e.DvceScreenResolution)
		e.DvceScreenWidth, e.DvceScreenHeight = &monDimension.width, &monDimension.height
	}
}

func setPageView(e *SnowplowEvent, params map[string]interface{}) {
	evnt := PageViewEvent{}
	sde := evnt.toSelfDescribing()
	e.SelfDescribingEvent = &sde
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
	if b64encodedContexts != nil {
		e.Contexts = getContexts(b64encodedContexts)
	}
}

func setSelfDescribing(e *SnowplowEvent, params map[string]interface{}) {
	b64EncodedPayload := getStringParam(params, "ue_px")
	e.SelfDescribingEvent = getSdPayload(b64EncodedPayload)
}

func buildEventFromMappedParams(c *gin.Context, params map[string]interface{}, conf config.Config) SnowplowEvent {
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
	switch event.Event {
	case PAGE_VIEW:
		setPageView(&event, params)
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
	return event
}
