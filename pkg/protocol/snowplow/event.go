// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package snowplow

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/util"
)

const (
	PAGE_PING               = "page_ping"
	PAGE_PING_SCHEMA        = "io.silverton/snowplow/page_ping/v1.0.json"
	PAGE_VIEW               = "page_view"
	PAGE_VIEW_SCHEMA        = "io.silverton/snowplow/page_view/v1.0.json"
	STRUCT_EVENT            = "struct_event"
	STRUCT_EVENT_SCHEMA     = "io.silverton/snowplow/struct/v1.0.json"
	TRANSACTION             = "transaction"
	TRANSACTION_SCHEMA      = "io.silverton/snowplow/transaction/v1.0.json"
	TRANSACTION_ITEM        = "transaction_item"
	TRANSACTION_ITEM_SCHEMA = "io.silverton/snowplow/transaction_item/v1.0.json"
	AD_IMPRESSION           = "ad_impression" // NOTE - already a self-describing event
	UNKNOWN_EVENT           = "unknown_event"
	UNKNOWN_SCHEMA          = "unknown_schema"
	SELF_DESCRIBING_EVENT   = "self_describing"
	IGLU                    = "iglu"
)

type SnowplowEvent struct {
	// Application parameters - https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/snowplow-tracker-protocol/#common-parameters-platform-and-event-independent
	NameTracker            *string                         `json:"name_tracker"`
	AppId                  *string                         `json:"app_id"`
	Platform               string                          `json:"platform"`
	EtlTstamp              *time.Time                      `json:"etl_tstamp"`
	DvceCreatedTstamp      *time.Time                      `json:"dvce_created_tstamp"`
	DvceSentTstamp         *time.Time                      `json:"dvce_sent_tstamp"`
	TrueTstamp             *time.Time                      `json:"true_tstamp"`
	CollectorTstamp        time.Time                       `json:"collector_tstamp"`
	DerivedTstamp          time.Time                       `json:"derived_tstamp"`
	OsTimezone             *string                         `json:"os_timezone"`
	Event                  string                          `json:"event"`
	TxnId                  *string                         `json:"txn_id"`
	EventId                *string                         `json:"event_id"`
	EventFingerprint       uuid.UUID                       `json:"event_fingerprint"`
	TrackerVersion         *string                         `json:"v_tracker"`
	CollectorVersion       *string                         `json:"v_collector"`
	EtlVersion             *string                         `json:"v_etl"`
	DomainUserid           *string                         `json:"domain_userid"`
	NetworkUserid          *string                         `json:"network_userid"`
	Userid                 *string                         `json:"user_id"`
	DomainSessionIdx       *int64                          `json:"domain_sessionidx"`
	DomainSessionId        *string                         `json:"domain_sessionid"`
	UserIpAddress          *string                         `json:"user_ipaddress"`
	Useragent              *string                         `json:"useragent"`
	UserFingerprint        *string                         `json:"user_fingerprint"`
	MacAddress             *string                         `json:"mac_address"`
	PageUrl                *string                         `json:"page_url"`
	PageTitle              *string                         `json:"page_title"`
	PageUrlScheme          *string                         `json:"page_urlscheme"`
	PageUrlHost            *string                         `json:"page_urlhost"`
	PageUrlPort            *string                         `json:"page_urlport"`
	PageUrlPath            *string                         `json:"page_urlpath"`
	PageUrlQuery           *map[string]interface{}         `json:"page_urlquery"`
	PageUrlFragment        *string                         `json:"page_urlfragment"`
	MktMedium              *string                         `json:"mkt_medium,omitempty"`
	MktSource              *string                         `json:"mkt_source,omitempty"`
	MktTerm                *string                         `json:"mkt_term,omitempty"`
	MktContent             *string                         `json:"mkt_content,omitempty"`
	MktCampaign            *string                         `json:"mkt_campaign,omitempty"`
	PageReferrer           *string                         `json:"page_referrer,omitempty"`
	RefrUrlScheme          *string                         `json:"refr_urlscheme,omitempty"`
	RefrUrlHost            *string                         `json:"refr_urlhost,omitempty"`
	RefrUrlPort            *string                         `json:"refr_urlport,omitempty"`
	RefrUrlPath            *string                         `json:"refr_urlpath,omitempty"`
	RefrUrlQuery           *map[string]interface{}         `json:"refr_urlquery,omitempty"`
	RefrUrlFragment        *string                         `json:"refr_urlfragment,omitempty"`
	RefrMedium             *string                         `json:"refr_medium,omitempty"`
	RefrSource             *string                         `json:"refr_source,omitempty"`
	RefrTerm               *string                         `json:"refr_term,omitempty"`
	RefrContent            *string                         `json:"refr_content,omitempty"`
	RefrCampaign           *string                         `json:"refr_campaign,omitempty"`
	RefrDomainUserId       *string                         `json:"refr_domain_userid,omitempty"`
	RefrDomainTstamp       *time.Time                      `json:"refr_domain_tstamp,omitempty"`
	BrCookies              *bool                           `json:"br_cookies,omitempty"`
	BrLang                 *string                         `json:"br_lang,omitempty"`
	BrFeaturesPdf          *bool                           `json:"br_features_pdf,omitempty"`
	BrFeaturesQuicktime    *bool                           `json:"br_features_quicktime,omitempty"`
	BrFeaturesRealplayer   *bool                           `json:"br_features_realplayer,omitempty"`
	BrFeaturesWindowsmedia *bool                           `json:"br_features_windowsmedia,omitempty"`
	BrFeaturesDirector     *bool                           `json:"br_features_director,omitempty"`
	BrFeaturesFlash        *bool                           `json:"br_features_flash,omitempty"`
	BrFeaturesJava         *bool                           `json:"br_features_java,omitempty"`
	BrFeaturesGears        *bool                           `json:"br_features_gears,omitempty"`
	BrFeaturesSilverlight  *bool                           `json:"br_features_silverlight,omitempty"`
	BrColordepth           *int64                          `json:"br_colordepth,omitempty"`
	ViewportSize           *string                         `json:"viewport_size,omitempty"`
	BrViewWidth            *int                            `json:"br_viewwidth,omitempty"`
	BrViewHeight           *int                            `json:"br_viewheight,omitempty"`
	DocCharset             *string                         `json:"doc_charset,omitempty"`
	DocSize                *string                         `json:"doc_size,omitempty"`
	DocWidth               *int                            `json:"doc_width,omitempty"`
	DocHeight              *int                            `json:"doc_height,omitempty"`
	DvceScreenResolution   *string                         `json:"dvce_screenresolution,omitempty"`
	DvceScreenWidth        *int                            `json:"dvce_screenwidth,omitempty"`
	DvceScreenHeight       *int                            `json:"dvce_screenheight,omitempty"`
	Contexts               *envelope.Contexts              `json:"contexts,omitempty"`
	SelfDescribingEvent    *envelope.SelfDescribingPayload `json:"self_describing_event"`
	PpXOffsetMin           *int64                          `json:"pp_xoffset_min,omitempty"`
	PpXOffsetMax           *int64                          `json:"pp_xoffset_max,omitempty"`
	PpYOffsetMin           *int64                          `json:"pp_yoffset_min,omitempty"`
	PpYOffsetMax           *int64                          `json:"pp_yoffset_max,omitempty"`
	SeCategory             *string                         `json:"se_category,omitempty"`
	SeAction               *string                         `json:"se_action,omitempty"`
	SeLabel                *string                         `json:"se_label,omitempty"`
	SeProperty             *string                         `json:"se_property,omitempty"`
	SeValue                *float64                        `json:"se_value,omitempty"`
	TrOrderId              *string                         `json:"tr_orderid,omitempty"`
	TrAffiliation          *string                         `json:"tr_affiliation,omitempty"`
	TrTotal                *float64                        `json:"tr_total,omitempty"`
	TrTax                  *float64                        `json:"tr_tax,omitempty"`
	TrShipping             *float64                        `json:"tr_shipping,omitempty"`
	TrCity                 *string                         `json:"tr_city,omitempty"`
	TrState                *string                         `json:"tr_state,omitempty"`
	TrCountry              *string                         `json:"tr_country,omitempty"`
	TrCurrency             *string                         `json:"tr_currency,omitempty"`
	TiOrderId              *string                         `json:"ti_orderid,omitempty"`
	TiSku                  *string                         `json:"ti_sku,omitempty"`
	TiName                 *string                         `json:"ti_name,omitempty"`
	TiCategory             *string                         `json:"ti_category,omitempty"`
	TiPrice                *float64                        `json:"ti_price,string,omitempty"`
	TiQuantity             *int64                          `json:"ti_quantity,omitempty"`
	TiCurrency             *string                         `json:"ti_currency,omitempty"`
	EventVendor            *string                         `json:"event_vendor,omitempty"`
	EventName              *string                         `json:"event_name,omitempty"`
	EventFormat            *string                         `json:"event_format,omitempty"`
	EventVersion           *string                         `json:"event_version,omitempty"`
}

func (e *SnowplowEvent) Map() map[string]interface{} {
	var i map[string]interface{}
	m, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("could not marshal snowplow event")
	} else {
		err := json.Unmarshal(m, &i)
		if err != nil {
			log.Error().Err(err).Msg("could not unmarshal snowplow event to map[string]interface{}")
		}
	}
	return i
}

type Page struct {
	Url      string                 `json:"url"`
	Title    *string                `json:"title"`
	Scheme   string                 `json:"scheme"`
	Host     string                 `json:"host"`
	Port     string                 `json:"port"`
	Path     string                 `json:"path"`
	Query    map[string]interface{} `json:"query"`
	Fragment *string                `json:"fragment"`
	Medium   *string                `json:"medium"`
	Source   *string                `json:"source"`
	Term     *string                `json:"term"`
	Content  *string                `json:"content"`
	Campaign *string                `json:"campaign"`
}

type PageViewEvent struct{}

func (e *PageViewEvent) toSelfDescribing() envelope.SelfDescribingPayload {
	return envelope.SelfDescribingPayload{
		Schema: PAGE_VIEW_SCHEMA,
		Data:   map[string]interface{}{},
	}
}

type PagePingEvent struct {
	PpXOffsetMin *int64 `json:"pp_xoffset_min"`
	PpXOffsetMax *int64 `json:"pp_xoffset_max"`
	PpYOffsetMin *int64 `json:"pp_yoffset_min"`
	PpYOffsetMax *int64 `json:"pp_yoffset_max"`
}

func (e *PagePingEvent) toSelfDescribing() envelope.SelfDescribingPayload {
	data, err := util.StructToMap(e)
	if err != nil {
		log.Error().Err(err).Msg("could not coerce page ping event to map")
	}
	return envelope.SelfDescribingPayload{
		Schema: PAGE_PING_SCHEMA,
		Data:   data,
	}
}

type StructEvent struct {
	SeCategory *string  `json:"se_category"`
	SeAction   *string  `json:"se_action"`
	SeLabel    *string  `json:"se_label"`
	SeProperty *string  `json:"se_property"`
	SeValue    *float64 `json:"se_value"`
}

func (e *StructEvent) toSelfDescribing() envelope.SelfDescribingPayload {
	data, err := util.StructToMap(e)
	if err != nil {
		log.Error().Err(err).Msg("could not coerce struct event to map")
	}
	return envelope.SelfDescribingPayload{
		Schema: STRUCT_EVENT_SCHEMA,
		Data:   data,
	}
}

type TransactionEvent struct {
	TrOrderId     *string  `json:"tr_orderid"`
	TrAffiliation *string  `json:"tr_affiliation"`
	TrTotal       *float64 `json:"tr_total"`
	TrTax         *float64 `json:"tr_tax"`
	TrShipping    *float64 `json:"tr_shipping"`
	TrCity        *string  `json:"tr_city"`
	TrState       *string  `json:"tr_state"`
	TrCountry     *string  `json:"tr_country"`
	TrCurrency    *string  `json:"tr_currency"`
}

func (e *TransactionEvent) toSelfDescribing() envelope.SelfDescribingPayload {
	data, err := util.StructToMap(e)
	if err != nil {
		log.Error().Err(err).Msg("could not coerce transaction event to map")
	}
	return envelope.SelfDescribingPayload{
		Schema: TRANSACTION_SCHEMA,
		Data:   data,
	}
}

type TransactionItemEvent struct {
	TiOrderId  *string  `json:"ti_orderid"`
	TiSku      *string  `json:"ti_sku"`
	TiName     *string  `json:"ti_name"`
	TiCategory *string  `json:"ti_category"`
	TiPrice    *float64 `json:"ti_price"`
	TiQuantity *int64   `json:"ti_quantity"`
	TiCurrency *string  `json:"ti_currency"`
}

func (e *TransactionItemEvent) toSelfDescribing() envelope.SelfDescribingPayload {
	data, err := util.StructToMap(e)
	if err != nil {
		log.Error().Err(err).Msg("could not coerce transaction item event to map")
	}
	return envelope.SelfDescribingPayload{
		Schema: TRANSACTION_ITEM_SCHEMA,
		Data:   data,
	}
}

func getEventType(param string) string {
	switch param {
	case "pp":
		return PAGE_PING
	case "pv":
		return PAGE_VIEW
	case "se":
		return STRUCT_EVENT
	case "ue":
		return SELF_DESCRIBING_EVENT
	case "tr":
		return TRANSACTION
	case "ti":
		return TRANSACTION_ITEM
	case "ad":
		return AD_IMPRESSION
	}
	return UNKNOWN_EVENT
}

type Dimension struct {
	height int
	width  int
}
