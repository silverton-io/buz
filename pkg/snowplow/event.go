package snowplow

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
)

const (
	PAGE_PING             = "page_ping"
	PAGE_VIEW             = "page_view"
	STRUCT_EVENT          = "struct_event"
	SELF_DESCRIBING_EVENT = "self_describing"
	TRANSACTION           = "transaction"
	TRANSACTION_ITEM      = "transaction_item"
	AD_IMPRESSION         = "ad_impression"
	UNKNOWN_EVENT         = "unknown_event"
	UNKNOWN_SCHEMA        = "unknown_schema"
	IGLU                  = "iglu"
)

type SnowplowEvent struct {
	// Application parameters - https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/snowplow-tracker-protocol/#common-parameters-platform-and-event-independent
	NameTracker       string     `json:"name_tracker"`
	AppId             string     `json:"app_id"`
	Platform          string     `json:"platform"`
	EtlTstamp         time.Time  `json:"etl_tstamp"`
	DvceCreatedTstamp time.Time  `json:"dvce_created_tstamp"`
	DvceSentTstamp    time.Time  `json:"dvce_sent_tstamp"`
	TrueTstamp        *time.Time `json:"true_tstamp"`
	CollectorTstamp   time.Time  `json:"collector_tstamp"`
	DerivedTstamp     time.Time  `json:"derived_tstamp"`
	OsTimezone        *string    `json:"os_timezone"`
	Event             string     `json:"event"`
	TxnId             *string    `json:"txn_id"` // deprecated
	EventId           *string    `json:"event_id"`
	EventFingerprint  uuid.UUID  `json:"event_fingerprint"`
	TrackerVersion    *string    `json:"v_tracker"`
	CollectorVersion  *string    `json:"v_collector"`
	EtlVersion        *string    `json:"v_etl"`
	// User fields
	DomainUserid     *string `json:"domain_userid"`
	NetworkUserid    *string `json:"network_userid"`
	Userid           *string `json:"user_id"`
	DomainSessionIdx *int64  `json:"domain_sessionidx"`
	DomainSessionId  *string `json:"domain_sessionid"`
	UserIpAddress    *string `json:"user_ipaddress"`
	Useragent        *string `json:"useragent"`
	UserFingerprint  *string `json:"user_fingerprint"`
	MacAddress       *string `json:"mac_address"`
	// Page fields
	PageUrl         *string `json:"page_url"`
	PageTitle       *string `json:"page_title"`
	PageUrlScheme   *string `json:"page_urlscheme"`
	PageUrlHost     *string `json:"page_urlhost"`
	PageUrlPort     *string `json:"page_urlport"`
	PageUrlPath     *string `json:"page_urlpath"`
	PageUrlQuery    *string `json:"page_urlquery"`
	PageUrlFragment *string `json:"page_urlfragment"`
	MktMedium       *string `json:"mkt_medium"`
	MktSource       *string `json:"mkt_source"`
	MktTerm         *string `json:"mkt_term"`
	MktContent      *string `json:"mkt_content"`
	MktCampaign     *string `json:"mkt_campaign"`
	// Referrer fields
	PageReferrer     *string    `json:"page_referrer"`
	RefrUrlScheme    *string    `json:"refr_urlscheme"`
	RefrUrlHost      *string    `json:"refr_urlhost"`
	RefrUrlPort      *string    `json:"refr_urlport"`
	RefrUrlPath      *string    `json:"refr_urlpath"`
	RefrUrlQuery     *string    `json:"refr_urlquery"`
	RefrUrlFragment  *string    `json:"refr_urlfragment"`
	RefrMedium       *string    `json:"refr_medium"`
	RefrSource       *string    `json:"refr_source"`
	RefrTerm         *string    `json:"refr_term"`
	RefrContent      *string    `json:"refr_content"`
	RefrCampaign     *string    `json:"refr_campaign"`
	RefrDomainUserId *string    `json:"refr_domain_userid"` // FIXME! Domain Linker
	RefrDomainTstamp *time.Time `json:"refr_domain_tstamp"` // FIXME! Domain Linker
	// Br features fields
	BrCookies              *bool   `json:"br_cookies"`
	BrLang                 *string `json:"br_lang"`
	BrFeaturesPdf          *bool   `json:"br_features_pdf"`          // to deprecate
	BrFeaturesQuicktime    *bool   `json:"br_features_quicktime"`    // to deprecate
	BrFeaturesRealplayer   *bool   `json:"br_features_realplayer"`   // to deprecate
	BrFeaturesWindowsmedia *bool   `json:"br_features_windowsmedia"` // to deprecate
	BrFeaturesDirector     *bool   `json:"br_features_director"`     // to deprecate
	BrFeaturesFlash        *bool   `json:"br_features_flash"`        // to deprecate
	BrFeaturesJava         *bool   `json:"br_features_java"`         // to deprecate
	BrFeaturesGears        *bool   `json:"br_features_gears"`        // to deprecate
	BrFeaturesSilverlight  *bool   `json:"br_features_silverlight"`  // to deprecate
	BrColordepth           *int64  `json:"br_colordepth"`
	// Dimension fields
	ViewportSize      *string `json:"viewport_size"`
	BrViewWidth       *int    `json:"br_viewwidth"`
	BrViewHeight      *int    `json:"br_viewheight"`
	DocCharset        *string `json:"doc_charset"`
	DocSize           *string `json:"doc_size"`
	DocWidth          *int    `json:"doc_width"`
	DocHeight         *int    `json:"doc_height"`
	MonitorResolution *string `json:"monitor_resolution"`
	DvceScreenWidth   *int    `json:"dvce_screenwidth"`
	DvceScreenHeight  *int    `json:"dvce_screenheight"`
	// Payload/context fields
	Contexts            *[]event.SelfDescribingContext `json:"contexts"`
	SelfDescribingEvent *event.SelfDescribingPayload   `json:"self_describing_event"` // Self Describing Event
	// Page ping fields
	PpXOffsetMin *int64 `json:"pp_xoffset_min"` // Page Ping Event
	PpXOffsetMax *int64 `json:"pp_xoffset_max"` // Page Ping Event
	PpYOffsetMin *int64 `json:"pp_yoffset_min"` // Page Ping Event
	PpYOffsetMax *int64 `json:"pp_yoffset_max"` // Page Ping Event
	// Struct fields
	SeCategory *string  `json:"se_category"` // Struct Event
	SeAction   *string  `json:"se_action"`   // Struct Event
	SeLabel    *string  `json:"se_label"`    // Struct Event
	SeProperty *string  `json:"se_property"` // Struct Event
	SeValue    *float64 `json:"se_value"`    // Struct Event
	// Transaction fields
	TrOrderId     *string  `json:"tr_orderid"`     // Transaction Event
	TrAffiliation *string  `json:"tr_affiliation"` // Transaction Event
	TrTotal       *float64 `json:"tr_total"`       // Transaction Event
	TrTax         *float64 `json:"tr_tax"`         // Transaction Event
	TrShipping    *float64 `json:"tr_shipping"`    // Transaction Event
	TrCity        *string  `json:"tr_city"`        // Transaction Event
	TrState       *string  `json:"tr_state"`       // Transaction Event
	TrCountry     *string  `json:"tr_country"`     // Transaction Event
	TrCurrency    *string  `json:"tr_currency"`    // Transaction Event
	// Transaction item fields
	TiOrderId  *string  `json:"ti_orderid"`      // Transaction Item Event
	TiSku      *string  `json:"ti_sku"`          // Transaction Item Event
	TiName     *string  `json:"ti_name"`         // Transaction Item Event
	TiCategory *string  `json:"ti_category"`     // Transaction Item Event
	TiPrice    *float64 `json:"ti_price,string"` // Transaction Item Event
	TiQuantity *int64   `json:"ti_quantity"`     // Transaction Item Event
	TiCurrency *string  `json:"ti_currency"`     // Transaction Item Event
	// Event fields
	EventVendor  *string `json:"event_vendor"`
	EventName    *string `json:"event_name"`
	EventFormat  *string `json:"event_format"`
	EventVersion *string `json:"event_version"`
}

func (e SnowplowEvent) Schema() *string {
	switch e.Event {
	case SELF_DESCRIBING_EVENT:
		schemaName := e.SelfDescribingEvent.Schema
		if schemaName[:4] == IGLU {
			schemaName = schemaName[5:]
		}
		return &schemaName
	default:
		schemaName := string(e.Event)
		return &schemaName
	}
}

func (e SnowplowEvent) Protocol() string {
	return protocol.SNOWPLOW
}

func (e SnowplowEvent) PayloadAsByte() ([]byte, error) {
	payloadBytes, err := json.Marshal(e.SelfDescribingEvent.Data)
	if err != nil {
		return nil, err
	}
	return payloadBytes, nil
}

func (e SnowplowEvent) AsByte() ([]byte, error) {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return eventBytes, nil
}

func (e SnowplowEvent) AsMap() (map[string]interface{}, error) {
	var m map[string]interface{}
	b, err := e.AsByte()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type SelfDescribingMetadata struct {
	Event_vendor  *string `json:"event_vendor"`
	Event_name    *string `json:"event_name"`
	Event_format  *string `json:"event_format"`
	Event_version *string `json:"event_version"`
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

type PageFields struct {
	scheme   string
	host     string
	path     string
	query    *string
	fragment *string
	medium   *string
	source   *string
	term     *string
	content  *string
	campaign *string
}
