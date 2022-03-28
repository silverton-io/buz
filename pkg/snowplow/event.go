package snowplow

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
)

// Event Types

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
)

// Other

const (
	IGLU = "iglu"
)

type SnowplowEvent struct {
	// Application parameters - https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/snowplow-tracker-protocol/#common-parameters-platform-and-event-independent
	Name_tracker        string     `json:"name_tracker"`
	App_id              string     `json:"app_id"`
	Platform            string     `json:"platform"`
	Etl_tstamp          time.Time  `json:"etl_tstamp"`
	Dvce_created_tstamp time.Time  `json:"dvce_created_tstamp"`
	Dvce_sent_tstamp    time.Time  `json:"dvce_sent_tstamp"`
	True_tstamp         *time.Time `json:"true_tstamp"`
	Collector_tstamp    time.Time  `json:"collector_tstamp"`
	Derived_tstamp      time.Time  `json:"derived_tstamp"`
	Os_timezone         *string    `json:"os_timezone"`
	Event               string     `json:"event"`
	Txn_id              *string    `json:"txn_id"` // deprecated
	Event_id            *string    `json:"event_id"`
	Event_fingerprint   uuid.UUID  `json:"event_fingerprint"`
	Tracker_version     *string    `json:"v_tracker"`
	Collector_version   *string    `json:"v_collector"`
	Etl_version         *string    `json:"v_etl"`
	// User fields
	Domain_userid     *string `json:"domain_userid"`
	Network_userid    *string `json:"network_userid"`
	Userid            *string `json:"user_id"`
	Domain_sessionidx *int64  `json:"domain_sessionidx"`
	Domain_sessionid  *string `json:"domain_sessionid"`
	User_ipaddress    *string `json:"user_ipaddress"`
	Useragent         *string `json:"useragent"`
	User_fingerprint  *string `json:"user_fingerprint"`
	Mac_address       *string `json:"mac_address"`
	// Page fields
	Page_url         *string `json:"page_url"`
	Page_urlscheme   *string `json:"page_urlscheme"`
	Page_urlhost     *string `json:"page_urlhost"`
	Page_urlport     *string `json:"page_urlport"`
	Page_urlpath     *string `json:"page_urlpath"`
	Page_urlquery    *string `json:"page_urlquery"`
	Page_urlfragment *string `json:"page_urlfragment"`
	Mkt_medium       *string `json:"mkt_medium"`
	Mkt_source       *string `json:"mkt_source"`
	Mkt_term         *string `json:"mkt_term"`
	Mkt_content      *string `json:"mkt_content"`
	Mkt_campaign     *string `json:"mkt_campaign"`
	Page_title       *string `json:"page_title"`
	// Referrer fields
	Page_referrer      *string    `json:"page_referrer"`
	Refr_urlscheme     *string    `json:"refr_urlscheme"`
	Refr_urlhost       *string    `json:"refr_urlhost"`
	Refr_urlport       *string    `json:"refr_urlport"`
	Refr_urlpath       *string    `json:"refr_urlpath"`
	Refr_urlquery      *string    `json:"refr_urlquery"`
	Refr_urlfragment   *string    `json:"refr_urlfragment"`
	Refr_medium        *string    `json:"refr_medium"`
	Refr_source        *string    `json:"refr_source"`
	Refr_term          *string    `json:"refr_term"`
	Refr_content       *string    `json:"refr_content"`
	Refr_campaign      *string    `json:"refr_campaign"`
	Refr_domain_userid *string    `json:"refr_domain_userid"` // FIXME! Domain Linker
	Refr_domain_tstamp *time.Time `json:"refr_domain_tstamp"` // FIXME! Domain Linker
	// Br features fields
	Br_cookies               *bool   `json:"br_cookies"`
	Br_lang                  *string `json:"br_lang"`
	Br_features_pdf          *bool   `json:"br_features_pdf"`          // to deprecate
	Br_features_quicktime    *bool   `json:"br_features_quicktime"`    // to deprecate
	Br_features_realplayer   *bool   `json:"br_features_realplayer"`   // to deprecate
	Br_features_windowsmedia *bool   `json:"br_features_windowsmedia"` // to deprecate
	Br_features_director     *bool   `json:"br_features_director"`     // to deprecate
	Br_features_flash        *bool   `json:"br_features_flash"`        // to deprecate
	Br_features_java         *bool   `json:"br_features_java"`         // to deprecate
	Br_features_gears        *bool   `json:"br_features_gears"`        // to deprecate
	Br_features_silverlight  *bool   `json:"br_features_silverlight"`  // to deprecate
	Br_colordepth            *int64  `json:"br_colordepth"`
	// Dimension fields
	Viewport_size      *string `json:"viewport_size"`
	Br_viewwidth       *int    `json:"br_viewwidth"`
	Br_viewheight      *int    `json:"br_viewheight"`
	Doc_charset        *string `json:"doc_charset"`
	Doc_size           *string `json:"doc_size"`
	Doc_width          *int    `json:"doc_width"`
	Doc_height         *int    `json:"doc_height"`
	Monitor_resolution *string `json:"monitor_resolution"`
	Dvce_screenwidth   *int    `json:"dvce_screenwidth"`
	Dvce_screenheight  *int    `json:"dvce_screenheight"`
	// Payload/context fields
	Contexts              *[]event.SelfDescribingContext `json:"contexts"`
	Self_describing_event *event.SelfDescribingEvent     `json:"self_describing_event"` // Self Describing Event
	// Page ping fields
	Pp_xoffset_min *int64 `json:"pp_xoffset_min"` // Page Ping Event
	Pp_xoffset_max *int64 `json:"pp_xoffset_max"` // Page Ping Event
	Pp_yoffset_min *int64 `json:"pp_yoffset_min"` // Page Ping Event
	Pp_yoffset_max *int64 `json:"pp_yoffset_max"` // Page Ping Event
	// Struct fields
	Se_category *string  `json:"se_category"` // Struct Event
	Se_action   *string  `json:"se_action"`   // Struct Event
	Se_label    *string  `json:"se_label"`    // Struct Event
	Se_property *string  `json:"se_property"` // Struct Event
	Se_value    *float64 `json:"se_value"`    // Struct Event
	// Transaction fields
	Tr_orderid     *string  `json:"tr_orderid"`     // Transaction Event
	Tr_affiliation *string  `json:"tr_affiliation"` // Transaction Event
	Tr_total       *float64 `json:"tr_total"`       // Transaction Event
	Tr_tax         *float64 `json:"tr_tax"`         // Transaction Event
	Tr_shipping    *float64 `json:"tr_shipping"`    // Transaction Event
	Tr_city        *string  `json:"tr_city"`        // Transaction Event
	Tr_state       *string  `json:"tr_state"`       // Transaction Event
	Tr_country     *string  `json:"tr_country"`     // Transaction Event
	Tr_currency    *string  `json:"tr_currency"`    // Transaction Event
	// Transaction item fields
	Ti_orderid  *string  `json:"ti_orderid"`      // Transaction Item Event
	Ti_sku      *string  `json:"ti_sku"`          // Transaction Item Event
	Ti_name     *string  `json:"ti_name"`         // Transaction Item Event
	Ti_category *string  `json:"ti_category"`     // Transaction Item Event
	Ti_price    *float64 `json:"ti_price,string"` // Transaction Item Event
	Ti_quantity *int64   `json:"ti_quantity"`     // Transaction Item Event
	Ti_currency *string  `json:"ti_currency"`     // Transaction Item Event
	// Event fields
	Event_vendor  *string `json:"event_vendor"`
	Event_name    *string `json:"event_name"`
	Event_format  *string `json:"event_format"`
	Event_version *string `json:"event_version"`
}

func (e SnowplowEvent) Schema() *string {
	switch e.Event {
	case SELF_DESCRIBING_EVENT:
		schemaName := e.Self_describing_event.Payload.Schema
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
	payloadBytes, err := json.Marshal(e.Self_describing_event.Payload)
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
