package snowplow

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/tidwall/gjson"
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
	Name_tracker             string                              `json:"name_tracker"`
	App_id                   string                              `json:"app_id"`
	Platform                 string                              `json:"platform"`
	Etl_tstamp               time.Time                           `json:"etl_tstamp"`
	Dvce_created_tstamp      MillisecondTimestampField           `json:"dvce_created_tstamp"`
	Dvce_sent_tstamp         MillisecondTimestampField           `json:"dvce_sent_tstamp"`
	True_tstamp              *MillisecondTimestampField          `json:"true_tstamp"`
	Collector_tstamp         time.Time                           `json:"collector_tstamp"`
	Derived_tstamp           time.Time                           `json:"derived_tstamp"`
	Os_timezone              *string                             `json:"os_timezone"`
	Event                    EventTypeField                      `json:"event"`
	Txn_id                   *string                             `json:"txn_id"` // deprecated
	Event_id                 *string                             `json:"event_id"`
	Event_fingerprint        *string                             `json:"event_fingerprint"`
	Tracker_version          *string                             `json:"v_tracker"`
	Collector_version        *string                             `json:"v_collector"`
	Etl_version              *string                             `json:"v_etl"`
	Domain_userid            *string                             `json:"domain_userid"`
	Network_userid           *string                             `json:"network_userid"`
	Userid                   *string                             `json:"user_id"`
	Domain_sessionidx        *int64                              `json:"domain_sessionidx"`
	Domain_sessionid         *string                             `json:"domain_sessionid"`
	User_ipaddress           *string                             `json:"user_ipaddress"`
	Page_url                 *string                             `json:"page_url"`
	Page_urlscheme           *string                             `json:"page_urlscheme"`
	Page_urlhost             *string                             `json:"page_urlhost"`
	Page_urlport             *string                             `json:"page_urlport"`
	Page_urlpath             *string                             `json:"page_urlpath"`
	Page_urlquery            *string                             `json:"page_urlquery"`
	Page_urlfragment         *string                             `json:"page_urlfragment"`
	Mkt_medium               *string                             `json:"mkt_medium"`
	Mkt_source               *string                             `json:"mkt_source"`
	Mkt_term                 *string                             `json:"mkt_term"`
	Mkt_content              *string                             `json:"mkt_content"`
	Mkt_campaign             *string                             `json:"mkt_campaign"`
	Useragent                *string                             `json:"useragent"`
	Page_title               *string                             `json:"page_title"`
	Page_referrer            *string                             `json:"page_referrer"`
	Refr_urlscheme           *string                             `json:"refr_urlscheme"`
	Refr_urlhost             *string                             `json:"refr_urlhost"`
	Refr_urlport             *string                             `json:"refr_urlport"`
	Refr_urlpath             *string                             `json:"refr_urlpath"`
	Refr_urlquery            *string                             `json:"refr_urlquery"`
	Refr_urlfragment         *string                             `json:"refr_urlfragment"`
	Refr_medium              *string                             `json:"refr_medium"`
	Refr_source              *string                             `json:"refr_source"`
	Refr_term                *string                             `json:"refr_term"`
	Refr_content             *string                             `json:"refr_content"`
	Refr_campaign            *string                             `json:"refr_campaign"`
	User_fingerprint         *string                             `json:"user_fingerprint"`
	Br_cookies               FlexibleBoolField                   `json:"br_cookies"`
	Br_lang                  *string                             `json:"br_lang"`
	Br_features_pdf          FlexibleBoolField                   `json:"br_features_pdf"`          // to deprecate
	Br_features_quicktime    FlexibleBoolField                   `json:"br_features_quicktime"`    // to deprecate
	Br_features_realplayer   FlexibleBoolField                   `json:"br_features_realplayer"`   // to deprecate
	Br_features_windowsmedia FlexibleBoolField                   `json:"br_features_windowsmedia"` // to deprecate
	Br_features_director     FlexibleBoolField                   `json:"br_features_director"`     // to deprecate
	Br_features_flash        FlexibleBoolField                   `json:"br_features_flash"`        // to deprecate
	Br_features_java         FlexibleBoolField                   `json:"br_features_java"`         // to deprecate
	Br_features_gears        FlexibleBoolField                   `json:"br_features_gears"`        // to deprecate
	Br_features_silverlight  FlexibleBoolField                   `json:"br_features_silverlight"`  // to deprecate
	Br_colordepth            *int                                `json:"br_colordepth"`
	Doc_charset              *string                             `json:"doc_charset"`
	Doc_size                 *string                             `json:"doc_size"`
	Doc_width                *int                                `json:"doc_width"`
	Doc_height               *int                                `json:"doc_height"`
	Viewport_size            *string                             `json:"viewport_size"`
	Br_viewwidth             *int                                `json:"br_viewwidth"`
	Br_viewheight            *int                                `json:"br_viewheight"`
	Monitor_resolution       *string                             `json:"monitor_resolution"`
	Dvce_screenwidth         *int                                `json:"dvce_screenwidth"`
	Dvce_screenheight        *int                                `json:"dvce_screenheight"`
	Mac_address              *string                             `json:"mac_address"`
	Contexts                 *Base64EncodedContexts              `json:"contexts"`
	Self_describing_event    *Base64EncodedSelfDescribingPayload `json:"self_describing_event"` // Self Describing Event
	Pp_xoffset_min           *int                                `json:"pp_xoffset_min"`        // Page Ping Event
	Pp_xoffset_max           *int                                `json:"pp_xoffset_max"`        // Page Ping Event
	Pp_yoffset_min           *int                                `json:"pp_yoffset_min"`        // Page Ping Event
	Pp_yoffset_max           *int                                `json:"pp_yoffset_max"`        // Page Ping Event
	Se_category              *string                             `json:"se_category"`           // Struct Event
	Se_action                *string                             `json:"se_action"`             // Struct Event
	Se_label                 *string                             `json:"se_label"`              // Struct Event
	Se_property              *string                             `json:"se_property"`           // Struct Event
	Se_value                 *float64                            `json:"se_value"`              // Struct Event
	Tr_orderid               *string                             `json:"tr_orderid"`            // Transaction Event
	Tr_affiliation           *string                             `json:"tr_affiliation"`        // Transaction Event
	Tr_total                 *float64                            `json:"tr_total"`              // Transaction Event
	Tr_tax                   *float64                            `json:"tr_tax"`                // Transaction Event
	Tr_shipping              *float64                            `json:"tr_shipping"`           // Transaction Event
	Tr_city                  *string                             `json:"tr_city"`               // Transaction Event
	Tr_state                 *string                             `json:"tr_state"`              // Transaction Event
	Tr_country               *string                             `json:"tr_country"`            // Transaction Event
	Tr_currency              *string                             `json:"tr_currency"`           // Transaction Event
	Ti_orderid               *string                             `json:"ti_orderid"`            // Transaction Item Event
	Ti_sku                   *string                             `json:"ti_sku"`                // Transaction Item Event
	Ti_name                  *string                             `json:"ti_name"`               // Transaction Item Event
	Ti_category              *string                             `json:"ti_category"`           // Transaction Item Event
	Ti_price                 *float64                            `json:"ti_price,string"`       // Transaction Item Event
	Ti_quantity              *int                                `json:"ti_quantity"`           // Transaction Item Event
	Ti_currency              *string                             `json:"ti_currency"`           // Transaction Item Event
	Refr_domain_userid       *string                             `json:"refr_domain_userid"`    // FIXME! Domain Linker
	Refr_domain_tstamp       *time.Time                          `json:"refr_domain_tstamp"`    // FIXME! Domain Linker
	Event_vendor             *string                             `json:"event_vendor"`
	Event_name               *string                             `json:"event_name"`
	Event_format             *string                             `json:"event_format"`
	Event_version            *string                             `json:"event_version"`
}

func (e SnowplowEvent) Schema() *string {
	switch e.Event {
	case SELF_DESCRIBING_EVENT:
		schemaName := e.Self_describing_event.Schema
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
	payloadBytes, err := json.Marshal(e.Self_describing_event.Data)
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

type ShortenedSnowplowEvent struct { //A struct used to quickly parse incoming json props or query params. Leverages Go type conversion to long-form props.
	Name_tracker             string                              `json:"tna"`
	App_id                   string                              `json:"aid"`
	Platform                 string                              `json:"p"`
	Etl_tstamp               time.Time                           `json:"etl_tstamp"` // not in the tracker protocol
	Dvce_created_tstamp      MillisecondTimestampField           `json:"dtm"`
	Dvce_sent_tstamp         MillisecondTimestampField           `json:"stm"`
	True_tstamp              *MillisecondTimestampField          `json:"ttm"`
	Collector_tstamp         time.Time                           `json:"collector_tstamp"` // not in the tracker protocol
	Derived_tstamp           time.Time                           `json:"derived_tstamp"`   // not in the tracker protocol
	Os_timezone              *string                             `json:"tz"`
	Event                    EventTypeField                      `json:"e"`
	Txn_id                   *string                             `json:"tid"` // deprecated
	Event_id                 *string                             `json:"eid"`
	Event_fingerprint        *string                             `json:"event_fingerprint"`
	Tracker_version          *string                             `json:"tv"`
	Collector_version        *string                             `json:"v_collector"`
	Etl_version              *string                             `json:"v_etl"`
	Domain_userid            *string                             `json:"duid"`
	Network_userid           *string                             `json:"nuid"`
	Userid                   *string                             `json:"uid"`
	Domain_sessionidx        *int64                              `json:"vid,string"`
	Domain_sessionid         *string                             `json:"sid"`
	User_ipaddress           *string                             `json:"ip"`
	Page_url                 *string                             `json:"url"`
	Page_urlscheme           *string                             `json:"page_urlscheme"`
	Page_urlhost             *string                             `json:"page_urlhost"`
	Page_urlport             *string                             `json:"page_urlport"`
	Page_urlpath             *string                             `json:"page_urlpath"`
	Page_urlquery            *string                             `json:"page_urlquery"`
	Page_urlfragment         *string                             `json:"page_urlfragment"`
	Mkt_medium               *string                             `json:"mkt_medium"`
	Mkt_source               *string                             `json:"mkt_source"`
	Mkt_term                 *string                             `json:"mkt_term"`
	Mkt_content              *string                             `json:"mkt_content"`
	Mkt_campaign             *string                             `json:"mkt_campaign"`
	Useragent                *string                             `json:"ua"`
	Page_title               *string                             `json:"page"`
	Page_referrer            *string                             `json:"refr"`
	Refr_urlscheme           *string                             `json:"refr_urlscheme"`
	Refr_urlhost             *string                             `json:"refr_urlhost"`
	Refr_urlport             *string                             `json:"refr_urlport"`
	Refr_urlpath             *string                             `json:"refr_urlpath"`
	Refr_urlquery            *string                             `json:"refr_urlquery"`
	Refr_urlfragment         *string                             `json:"refr_urlfragment"`
	Refr_medium              *string                             `json:"refr_medium"`
	Refr_source              *string                             `json:"refr_source"`
	Refr_term                *string                             `json:"refr_term"`
	Refr_content             *string                             `json:"refr_content"`
	Refr_campaign            *string                             `json:"refr_campaign"`
	User_fingerprint         *string                             `json:"fp"` // deprecated
	Br_cookies               FlexibleBoolField                   `json:"cookie"`
	Br_lang                  *string                             `json:"lang"`
	Br_features_pdf          FlexibleBoolField                   `json:"f_pdf"`   // to deprecate
	Br_features_quicktime    FlexibleBoolField                   `json:"f_qt"`    // to deprecate
	Br_features_realplayer   FlexibleBoolField                   `json:"f_realp"` // to deprecate
	Br_features_windowsmedia FlexibleBoolField                   `json:"f_wma"`   // to deprecate
	Br_features_director     FlexibleBoolField                   `json:"f_dir"`   // to deprecate
	Br_features_flash        FlexibleBoolField                   `json:"f_fla"`   // to deprecate
	Br_features_java         FlexibleBoolField                   `json:"f_java"`  // to deprecate
	Br_features_gears        FlexibleBoolField                   `json:"f_gears"` // to deprecate
	Br_features_silverlight  FlexibleBoolField                   `json:"f_ag"`    // to deprecate
	Br_colordepth            *int                                `json:"cd,string"`
	Doc_charset              *string                             `json:"cs"`
	Doc_size                 *string                             `json:"ds"`
	Doc_width                *int                                `json:"doc_width"`
	Doc_height               *int                                `json:"doc_height"`
	Viewport_size            *string                             `json:"vp"`
	Br_viewwidth             *int                                `json:"br_viewwidth"`
	Br_viewheight            *int                                `json:"br_viewheight"`
	Monitor_resolution       *string                             `json:"res"`
	Dvce_screenwidth         *int                                `json:"dvce_screenwidth"`
	Dvce_screenheight        *int                                `json:"dvce_screenheight"`
	Mac_address              *string                             `json:"mac"`
	Contexts                 *Base64EncodedContexts              `json:"cx"`
	Self_describing_event    *Base64EncodedSelfDescribingPayload `json:"ue_px"`              // Self Describing Event
	Pp_xoffset_min           *int                                `json:"pp_mix,string"`      // Page Ping Event
	Pp_xoffset_max           *int                                `json:"pp_max,string"`      // Page Ping Event
	Pp_yoffset_min           *int                                `json:"pp_miy,string"`      // Page Ping Event
	Pp_yoffset_max           *int                                `json:"pp_may,string"`      // Page Ping Event
	Se_category              *string                             `json:"se_ca"`              // Struct Event
	Se_action                *string                             `json:"se_ac"`              // Struct Event
	Se_label                 *string                             `json:"se_la"`              // Struct Event
	Se_property              *string                             `json:"se_pr"`              // Struct Event
	Se_value                 *float64                            `json:"se_va,string"`       // Struct Event
	Tr_orderid               *string                             `json:"tr_id"`              // Transaction Event
	Tr_affiliation           *string                             `json:"tr_af"`              // Transaction Event
	Tr_total                 *float64                            `json:"tr_tt,string"`       // Transaction Event
	Tr_tax                   *float64                            `json:"tr_tx,string"`       // Transaction Event
	Tr_shipping              *float64                            `json:"tr_sh,string"`       // Transaction Event
	Tr_city                  *string                             `json:"tr_ci"`              // Transaction Event
	Tr_state                 *string                             `json:"tr_st"`              // Transaction Event
	Tr_country               *string                             `json:"tr_co"`              // Transaction Event
	Tr_currency              *string                             `json:"tr_cu"`              // Transaction Event
	Ti_orderid               *string                             `json:"ti_id"`              // Transaction Item Event
	Ti_sku                   *string                             `json:"ti_sk"`              // Transaction Item Event
	Ti_name                  *string                             `json:"ti_nm"`              // Transaction Item Event
	Ti_category              *string                             `json:"ti_ca"`              // Transaction Item Event
	Ti_price                 *float64                            `json:"ti_pr,string"`       // Transaction Item Event
	Ti_quantity              *int                                `json:"ti_qu,string"`       // Transaction Item Event
	Ti_currency              *string                             `json:"ti_cu"`              // Transaction Item Event
	Refr_domain_userid       *string                             `json:"refr_domain_userid"` // FIXME! Domain Linker
	Refr_domain_tstamp       *time.Time                          `json:"refr_domain_tstamp"` // FIXME! Domain Linker
	Event_vendor             *string                             `json:"event_vendor"`       // Self-describing event metadata
	Event_name               *string                             `json:"event_name"`         // Self-describing event metadata
	Event_format             *string                             `json:"event_format"`       // Self-describing event metadata
	Event_version            *string                             `json:"event_version"`      // Self-describing event metadata
}

type SelfDescribingMetadata struct {
	Event_vendor  *string `json:"event_vendor"`
	Event_name    *string `json:"event_name"`
	Event_format  *string `json:"event_format"`
	Event_version *string `json:"event_version"`
}

type Base64EncodedContexts []event.SelfDescribingContext

func (c *Base64EncodedContexts) UnmarshalJSON(bytes []byte) error {
	var encodedPayload string
	var contexts []event.SelfDescribingContext
	err := json.Unmarshal(bytes, &encodedPayload)
	decodedPayload, err := b64.RawStdEncoding.DecodeString(encodedPayload)
	if err != nil {
		fmt.Printf("error decoding b64 encoded contexts %s\n", err)
	}
	contextPayload := gjson.Parse(string(decodedPayload))
	for _, pl := range contextPayload.Get("data").Array() {
		context := event.SelfDescribingContext{
			Schema: pl.Get("schema").String(),
			Data:   pl.Get("data").Value().(map[string]interface{}),
		}
		contexts = append(contexts, context)
	}
	*c = contexts
	return nil
}

type Base64EncodedSelfDescribingPayload event.SelfDescribingPayload

func (f *Base64EncodedSelfDescribingPayload) UnmarshalJSON(bytes []byte) error {
	var encodedPayload string
	err := json.Unmarshal(bytes, &encodedPayload)
	decodedPayload, err := b64.RawStdEncoding.DecodeString(encodedPayload)
	if err != nil {
		fmt.Printf("error decoding b64 encoded self describing payload %s\n", err)
	}
	schema := gjson.GetBytes(decodedPayload, "data.schema").String()
	data := gjson.GetBytes(decodedPayload, "data.data").Value().(map[string]interface{})
	f.Schema = schema
	f.Data = data
	return nil
}

type FlexibleBoolField bool

func (b *FlexibleBoolField) UnmarshalJSON(bytes []byte) error {
	var payload string
	err := json.Unmarshal(bytes, &payload)
	if err != nil {
		fmt.Printf("error decoding FlexibleBoolField %s\n", err)
	}
	val, err := strconv.ParseBool(payload)
	*b = FlexibleBoolField(val)
	return nil
}

type MillisecondTimestampField struct {
	time.Time
}

func (t *MillisecondTimestampField) UnmarshalJSON(bytes []byte) error {
	var msString string
	err := json.Unmarshal(bytes, &msString)
	msInt, err := strconv.ParseInt(msString, 10, 64)
	if err != nil {
		fmt.Printf("error decoding timestamp: %s\n", err)
		return err
	}
	t.Time = time.Unix(0, msInt*int64(time.Millisecond))
	return nil
}

type EventTypeField string

func (f *EventTypeField) UnmarshalJSON(bytes []byte) error {
	var payload string
	err := json.Unmarshal(bytes, &payload)
	if err != nil {
		fmt.Printf("error unmarshaling EventTypeField %s\n", err)
	}
	eventType := getEventType(payload)
	*f = EventTypeField(eventType)
	return nil
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
