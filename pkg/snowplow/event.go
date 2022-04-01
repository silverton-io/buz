package snowplow

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/util"
)

const (
	PAGE_PING               = "page_ping"
	PAGE_PING_SCHEMA        = "com.silverton.io/snowplow/page_ping/v1.0.json"
	PAGE_VIEW               = "page_view"
	PAGE_VIEW_SCHEMA        = "com.silverton.io/snowplow/page_view/v1.0.json"
	STRUCT_EVENT            = "struct_event"
	STRUCT_EVENT_SCHEMA     = "com.silverton.io/snowplow/struct/v1.0.json"
	TRANSACTION             = "transaction"
	TRANSACTION_SCHEMA      = "com.silverton.io/snowplow/transaction/v1.0.json"
	TRANSACTION_ITEM        = "transaction_item"
	TRANSACTION_ITEM_SCHEMA = "com.silverton.io/snowplow/transaction_item/v1.0.json"
	AD_IMPRESSION           = "ad_impression" // NOTE - already a self-describing event.
	UNKNOWN_EVENT           = "unknown_event"
	UNKNOWN_SCHEMA          = "unknown_schema"
	SELF_DESCRIBING_EVENT   = "self_describing"
	IGLU                    = "iglu"
)

type SnowplowEvent struct {
	Tstamp                 `json:"tstamp"`
	PlatformMetadata       `json:"platform_metadata"`
	Event                  `json:"event"`
	User                   `json:"user"`
	Session                `json:"session"`
	Page                   `json:"page"`
	Referrer               Page `json:"referrer"`
	DomainLinker           `json:"domain_linker"`
	Device                 `json:"device"`
	Browser                `json:"browser"`
	Screen                 `json:"screen"`
	Contexts               *[]event.SelfDescribingContext `json:"contexts"`
	SelfDescribingEvent    *event.SelfDescribingPayload   `json:"self_describing_event"` // Self Describing Event
	SelfDescribingMetadata `json:"self_describing_metadata"`
}

func (e SnowplowEvent) Schema() *string {
	evnt := e.Event.Event
	switch evnt {
	case SELF_DESCRIBING_EVENT:
		schemaName := e.SelfDescribingEvent.Schema
		if schemaName[:4] == IGLU {
			schemaName = schemaName[5:]
		}
		return &schemaName
	default:
		schemaName := string(evnt)
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

type PlatformMetadata struct {
	NameTracker      string  `json:"name_tracker"`
	TrackerVersion   *string `json:"v_tracker"`
	CollectorVersion *string `json:"v_collector"`
	EtlVersion       *string `json:"v_etl"`
}

type Event struct {
	AppId            string    `json:"app_id"`
	Platform         string    `json:"platform"`
	Event            string    `json:"event"`
	TxnId            *string   `json:"txn_id"` // deprecated
	EventId          *string   `json:"event_id"`
	EventFingerprint uuid.UUID `json:"event_fingerprint"`
}

type Tstamp struct {
	DvceCreatedTstamp time.Time  `json:"dvce_created_tstamp"`
	DvceSentTstamp    time.Time  `json:"dvce_sent_tstamp"`
	TrueTstamp        *time.Time `json:"true_tstamp"`
	CollectorTstamp   time.Time  `json:"collector_tstamp"`
	EtlTstamp         time.Time  `json:"etl_tstamp"`
	DerivedTstamp     time.Time  `json:"derived_tstamp"`
}

type User struct {
	DomainUserid    *string `json:"domain_userid"`
	NetworkUserid   *string `json:"network_userid"`
	Userid          *string `json:"user_id"`
	UserIpAddress   *string `json:"user_ipaddress"`
	UserFingerprint *string `json:"user_fingerprint"`
}

type Device struct {
	Useragent  *string `json:"useragent"`
	MacAddress *string `json:"mac_address"`
	OsTimezone *string `json:"os_timezone"`
}

type Browser struct {
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
}

type Screen struct {
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
}

type Session struct {
	DomainSessionIdx *int64  `json:"domain_sessionidx"`
	DomainSessionId  *string `json:"domain_sessionid"`
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

type SelfDescribingMetadata struct {
	EventVendor  *string `json:"event_vendor"`
	EventName    *string `json:"event_name"`
	EventFormat  *string `json:"event_format"`
	EventVersion *string `json:"event_version"`
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

type Page struct {
	Url      string  `json:"url"`
	Title    *string `json:"title"`
	Scheme   string  `json:"scheme"`
	Host     string  `json:"host"`
	Port     string  `json:"port"`
	Path     string  `json:"path"`
	Query    *string `json:"query"`
	Fragment *string `json:"fragment"`
	Medium   *string `json:"medium"`
	Source   *string `json:"source"`
	Term     *string `json:"term"`
	Content  *string `json:"content"`
	Campaign *string `json:"campaign"`
}

type DomainLinker struct {
	RefrDomainUserId *string    `json:"refr_domain_userid"` // FIXME! Domain Linker
	RefrDomainTstamp *time.Time `json:"refr_domain_tstamp"` // FIXME! Domain Linker
}

type Dimension struct {
	height int
	width  int
}
