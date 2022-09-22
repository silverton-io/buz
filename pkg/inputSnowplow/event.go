// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputsnowplow

import (
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/event"
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
	AD_IMPRESSION           = "ad_impression" // NOTE - already a self-describing event.
	UNKNOWN_EVENT           = "unknown_event"
	UNKNOWN_SCHEMA          = "unknown_schema"
	SELF_DESCRIBING_EVENT   = "self_describing"
	IGLU                    = "iglu"
)

type SnowplowEvent struct {
	// Application parameters - https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/snowplow-tracker-protocol/#common-parameters-platform-and-event-independent
	NameTracker            *string                      `json:"name_tracker"`
	AppId                  *string                      `json:"app_id"`
	Platform               string                       `json:"platform"`
	EtlTstamp              *time.Time                   `json:"etl_tstamp"`
	DvceCreatedTstamp      *time.Time                   `json:"dvce_created_tstamp"`
	DvceSentTstamp         *time.Time                   `json:"dvce_sent_tstamp"`
	TrueTstamp             *time.Time                   `json:"true_tstamp"`
	CollectorTstamp        time.Time                    `json:"collector_tstamp"`
	DerivedTstamp          time.Time                    `json:"derived_tstamp"`
	OsTimezone             *string                      `json:"os_timezone"`
	Event                  string                       `json:"event"`
	TxnId                  *string                      `json:"txn_id"`
	EventId                *string                      `json:"event_id"`
	EventFingerprint       uuid.UUID                    `json:"event_fingerprint"`
	TrackerVersion         *string                      `json:"v_tracker"`
	CollectorVersion       *string                      `json:"v_collector"`
	EtlVersion             *string                      `json:"v_etl"`
	DomainUserid           *string                      `json:"domain_userid"`
	NetworkUserid          *string                      `json:"network_userid"`
	Userid                 *string                      `json:"user_id"`
	DomainSessionIdx       *int64                       `json:"domain_sessionidx"`
	DomainSessionId        *string                      `json:"domain_sessionid"`
	UserIpAddress          *string                      `json:"user_ipaddress"`
	Useragent              *string                      `json:"useragent"`
	UserFingerprint        *string                      `json:"user_fingerprint"`
	MacAddress             *string                      `json:"mac_address"`
	PageUrl                *string                      `json:"page_url"`
	PageTitle              *string                      `json:"page_title"`
	PageUrlScheme          *string                      `json:"page_urlscheme"`
	PageUrlHost            *string                      `json:"page_urlhost"`
	PageUrlPort            *string                      `json:"page_urlport"`
	PageUrlPath            *string                      `json:"page_urlpath"`
	PageUrlQuery           *map[string]interface{}      `json:"page_urlquery"`
	PageUrlFragment        *string                      `json:"page_urlfragment"`
	MktMedium              *string                      `json:"mkt_medium"`
	MktSource              *string                      `json:"mkt_source"`
	MktTerm                *string                      `json:"mkt_term"`
	MktContent             *string                      `json:"mkt_content"`
	MktCampaign            *string                      `json:"mkt_campaign"`
	PageReferrer           *string                      `json:"page_referrer"`
	RefrUrlScheme          *string                      `json:"refr_urlscheme"`
	RefrUrlHost            *string                      `json:"refr_urlhost"`
	RefrUrlPort            *string                      `json:"refr_urlport"`
	RefrUrlPath            *string                      `json:"refr_urlpath"`
	RefrUrlQuery           *map[string]interface{}      `json:"refr_urlquery"`
	RefrUrlFragment        *string                      `json:"refr_urlfragment"`
	RefrMedium             *string                      `json:"refr_medium"`
	RefrSource             *string                      `json:"refr_source"`
	RefrTerm               *string                      `json:"refr_term"`
	RefrContent            *string                      `json:"refr_content"`
	RefrCampaign           *string                      `json:"refr_campaign"`
	RefrDomainUserId       *string                      `json:"refr_domain_userid"` // FIXME!
	RefrDomainTstamp       *time.Time                   `json:"refr_domain_tstamp"` // FIXME!
	BrCookies              *bool                        `json:"br_cookies"`
	BrLang                 *string                      `json:"br_lang"`
	BrFeaturesPdf          *bool                        `json:"br_features_pdf"`
	BrFeaturesQuicktime    *bool                        `json:"br_features_quicktime"`
	BrFeaturesRealplayer   *bool                        `json:"br_features_realplayer"`
	BrFeaturesWindowsmedia *bool                        `json:"br_features_windowsmedia"`
	BrFeaturesDirector     *bool                        `json:"br_features_director"`
	BrFeaturesFlash        *bool                        `json:"br_features_flash"`
	BrFeaturesJava         *bool                        `json:"br_features_java"`
	BrFeaturesGears        *bool                        `json:"br_features_gears"`
	BrFeaturesSilverlight  *bool                        `json:"br_features_silverlight"`
	BrColordepth           *int64                       `json:"br_colordepth"`
	ViewportSize           *string                      `json:"viewport_size"`
	BrViewWidth            *int                         `json:"br_viewwidth"`
	BrViewHeight           *int                         `json:"br_viewheight"`
	DocCharset             *string                      `json:"doc_charset"`
	DocSize                *string                      `json:"doc_size"`
	DocWidth               *int                         `json:"doc_width"`
	DocHeight              *int                         `json:"doc_height"`
	DvceScreenResolution   *string                      `json:"dvce_screenresolution"`
	DvceScreenWidth        *int                         `json:"dvce_screenwidth"`
	DvceScreenHeight       *int                         `json:"dvce_screenheight"`
	Contexts               *map[string]interface{}      `json:"contexts"`
	SelfDescribingEvent    *event.SelfDescribingPayload `json:"self_describing_event"`
	PpXOffsetMin           *int64                       `json:"pp_xoffset_min"`
	PpXOffsetMax           *int64                       `json:"pp_xoffset_max"`
	PpYOffsetMin           *int64                       `json:"pp_yoffset_min"`
	PpYOffsetMax           *int64                       `json:"pp_yoffset_max"`
	SeCategory             *string                      `json:"se_category"`
	SeAction               *string                      `json:"se_action"`
	SeLabel                *string                      `json:"se_label"`
	SeProperty             *string                      `json:"se_property"`
	SeValue                *float64                     `json:"se_value"`
	TrOrderId              *string                      `json:"tr_orderid"`
	TrAffiliation          *string                      `json:"tr_affiliation"`
	TrTotal                *float64                     `json:"tr_total"`
	TrTax                  *float64                     `json:"tr_tax"`
	TrShipping             *float64                     `json:"tr_shipping"`
	TrCity                 *string                      `json:"tr_city"`
	TrState                *string                      `json:"tr_state"`
	TrCountry              *string                      `json:"tr_country"`
	TrCurrency             *string                      `json:"tr_currency"`
	TiOrderId              *string                      `json:"ti_orderid"`
	TiSku                  *string                      `json:"ti_sku"`
	TiName                 *string                      `json:"ti_name"`
	TiCategory             *string                      `json:"ti_category"`
	TiPrice                *float64                     `json:"ti_price,string"`
	TiQuantity             *int64                       `json:"ti_quantity"`
	TiCurrency             *string                      `json:"ti_currency"`
	EventVendor            *string                      `json:"event_vendor"`
	EventName              *string                      `json:"event_name"`
	EventFormat            *string                      `json:"event_format"`
	EventVersion           *string                      `json:"event_version"`
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

func (e *PageViewEvent) toSelfDescribing() event.SelfDescribingPayload {
	return event.SelfDescribingPayload{
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

func (e *PagePingEvent) toSelfDescribing() event.SelfDescribingPayload {
	return event.SelfDescribingPayload{
		Schema: PAGE_PING_SCHEMA,
		Data:   util.StructToMap(e),
	}
}

type StructEvent struct {
	SeCategory *string  `json:"se_category"`
	SeAction   *string  `json:"se_action"`
	SeLabel    *string  `json:"se_label"`
	SeProperty *string  `json:"se_property"`
	SeValue    *float64 `json:"se_value"`
}

func (e *StructEvent) toSelfDescribing() event.SelfDescribingPayload {
	return event.SelfDescribingPayload{
		Schema: STRUCT_EVENT_SCHEMA,
		Data:   util.StructToMap(e),
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

func (e *TransactionEvent) toSelfDescribing() event.SelfDescribingPayload {
	return event.SelfDescribingPayload{
		Schema: TRANSACTION_SCHEMA,
		Data:   util.StructToMap(e),
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

func (e *TransactionItemEvent) toSelfDescribing() event.SelfDescribingPayload {
	return event.SelfDescribingPayload{
		Schema: TRANSACTION_ITEM_SCHEMA,
		Data:   util.StructToMap(e),
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
