package snowplow

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type Event struct {
	// Application parameters - https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/snowplow-tracker-protocol/#common-parameters-platform-and-event-independent
	Name_tracker             string                              `json:"name_tracker"`
	Event_vendor             *string                             `json:"event_vendor,omitempty"` // deprecated
	App_id                   string                              `json:"app_id"`
	Platform                 string                              `json:"platform"`
	Etl_tstamp               *time.Time                          `json:"etl_tstamp,omitempty"`
	Dvce_created_tstamp      MillisecondTimestampField           `json:"dvce_created_tstamp"`
	Dvce_sent_tstamp         MillisecondTimestampField           `json:"dvce_sent_tstamp"`
	True_tstamp              *MillisecondTimestampField          `json:"true_tstamp,omitempty"`
	Collector_tstamp         time.Time                           `json:"collector_tstamp"`
	Derived_tstamp           *time.Time                          `json:"derived_tstamp,omitempty"`
	Os_timezone              *string                             `json:"os_timezone,omitempty"`
	Event                    EventTypeField                      `json:"event"`
	Txn_id                   *string                             `json:"txn_id,omitempty"` // deprecated
	Event_id                 string                              `json:"event_id"`
	Event_fingerprint        string                              `json:"event_fingerprint"`
	Tracker_version          string                              `json:"v_tracker"`
	Collector_version        *string                             `json:"v_collector"`
	Etl_version              *string                             `json:"v_etl,omitempty"`
	Domain_userid            string                              `json:"domain_userid"`
	Network_userid           *string                             `json:"network_userid,omitempty"`
	Userid                   *string                             `json:"user_id,omitempty"`
	Domain_sessionidx        *int64                              `json:"domain_sessionidx,omitempty"`
	Domain_sessionid         *string                             `json:"domain_sessionid,omitempty"`
	User_ipaddress           string                              `json:"user_ipaddress"`
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
	User_fingerprint         *string                             `json:"user_fingerprint,omitempty"`
	Br_cookies               FlexibleBoolField                   `json:"br_cookies"`
	Br_lang                  *string                             `json:"br_lang,omitempty"`
	Br_features_pdf          FlexibleBoolField                   `json:"br_features_pdf"`          // to deprecate
	Br_features_quicktime    FlexibleBoolField                   `json:"br_features_quicktime"`    // to deprecate
	Br_features_realplayer   FlexibleBoolField                   `json:"br_features_realplayer"`   // to deprecate
	Br_features_windowsmedia FlexibleBoolField                   `json:"br_features_windowsmedia"` // to deprecate
	Br_features_director     FlexibleBoolField                   `json:"br_features_director"`     // to deprecate
	Br_features_flash        FlexibleBoolField                   `json:"br_features_flash"`        // to deprecate
	Br_features_java         FlexibleBoolField                   `json:"br_features_java"`         // to deprecate
	Br_features_gears        FlexibleBoolField                   `json:"br_features_gears"`        // to deprecate
	Br_features_silverlight  FlexibleBoolField                   `json:"br_features_silverlight"`  // to deprecate
	Br_colordepth            *int                                `json:"br_colordepth,omitempty"`
	Doc_charset              *string                             `json:"doc_charset,omitempty"`
	Doc_size                 *string                             `json:"doc_size,omitempty"`
	Doc_width                *int                                `json:"doc_width,omitempty"`
	Doc_height               *int                                `json:"doc_height,omitempty"`
	Viewport_size            *string                             `json:"viewport_size,omitempty"`
	Br_viewwidth             *int                                `json:"br_viewwidth,omitempty"`
	Br_viewheight            *int                                `json:"br_viewheight,omitempty"`
	Monitor_resolution       *string                             `json:"monitor_resolution,omitempty"`
	Dvce_screenwidth         *int                                `json:"dvce_screenwidth,omitempty"`
	Dvce_screenheight        *int                                `json:"dvce_screenheight,omitempty"`
	Mac_address              *string                             `json:"mac_address,omitempty"`
	Contexts                 *Base64EncodedContexts              `json:"contexts"`
	Self_describing_event    *Base64EncodedSelfDescribingPayload `json:"self_describing_event,omitempty"` // Self Describing Event
	Pp_xoffset_min           *int                                `json:"pp_xoffset_min,omitempty"`        // Page Ping Event
	Pp_xoffset_max           *int                                `json:"pp_xoffset_max,omitempty"`        // Page Ping Event
	Pp_yoffset_min           *int                                `json:"pp_yoffset_min,omitempty"`        // Page Ping Event
	Pp_yoffset_max           *int                                `json:"pp_yoffset_max,omitempty"`        // Page Ping Event
	Se_category              *string                             `json:"se_category,omitempty"`           // Struct Event
	Se_action                *string                             `json:"se_action,omitempty"`             // Struct Event
	Se_label                 *string                             `json:"se_label,omitempty"`              // Struct Event
	Se_property              *string                             `json:"se_property,omitempty"`           // Struct Event
	Se_value                 *float64                            `json:"se_value,omitempty"`              // Struct Event
	Tr_orderid               *string                             `json:"tr_orderid,omitempty"`            // Transaction Event
	Tr_affiliation           *string                             `json:"tr_affiliation,omitempty"`        // Transaction Event
	Tr_total                 *float64                            `json:"tr_total,omitempty"`              // Transaction Event
	Tr_tax                   *float64                            `json:"tr_tax,omitempty"`                // Transaction Event
	Tr_shipping              *float64                            `json:"tr_shipping,omitempty"`           // Transaction Event
	Tr_city                  *string                             `json:"tr_city,omitempty"`               // Transaction Event
	Tr_state                 *string                             `json:"tr_state,omitempty"`              // Transaction Event
	Tr_country               *string                             `json:"tr_country,omitempty"`            // Transaction Event
	Tr_currency              *string                             `json:"tr_currency,omitempty"`           // Transaction Event
	Ti_orderid               *string                             `json:"ti_orderid,omitempty"`            // Transaction Item Event
	Ti_sku                   *string                             `json:"ti_sku,omitempty"`                // Transaction Item Event
	Ti_name                  *string                             `json:"ti_name,omitempty"`               // Transaction Item Event
	Ti_category              *string                             `json:"ti_category,omitempty"`           // Transaction Item Event
	Ti_price                 *float64                            `json:"ti_price,string,omitempty"`       // Transaction Item Event
	Ti_quantity              *int                                `json:"ti_quantity,omitempty"`           // Transaction Item Event
	Ti_currency              *string                             `json:"ti_currency,omitempty"`           // Transaction Item Event
	Refr_domain_userid       *string                             `json:"refr_domain_userid,omitempty"`    // FIXME! Domain Linker
	Refr_domain_tstamp       time.Time                           `json:"refr_domain_tstamp,omitempty"`    // FIXME! Domain Linker
}

type ShortenedEvent struct { //A struct used to quickly parse incoming json props or query params. Leverages Go type conversion to long-form props.
	Name_tracker             string                              `json:"tna"`
	Event_vendor             *string                             `json:"evn,omitempty"` // deprecated
	App_id                   string                              `json:"aid"`
	Platform                 string                              `json:"p"`
	Etl_tstamp               *time.Time                          `json:"etl_tstamp,omitempty"` // not in the tracker protocol
	Dvce_created_tstamp      MillisecondTimestampField           `json:"dtm"`
	Dvce_sent_tstamp         MillisecondTimestampField           `json:"stm"`
	True_tstamp              *MillisecondTimestampField          `json:"ttm,omitempty"`
	Collector_tstamp         time.Time                           `json:"collector_tstamp"`         // not in the tracker protocol
	Derived_tstamp           *time.Time                          `json:"derived_tstamp,omitempty"` // not in the tracker protocol
	Os_timezone              *string                             `json:"tz,omitempty"`
	Event                    EventTypeField                      `json:"e"`
	Txn_id                   *string                             `json:"tid,omitempty"` // deprecated
	Event_id                 string                              `json:"eid"`
	Event_fingerprint        string                              `json:"event_fingerprint"`
	Tracker_version          string                              `json:"tv"`
	Collector_version        *string                             `json:"v_collector"`
	Etl_version              *string                             `json:"v_etl,omitempty"`
	Domain_userid            string                              `json:"duid"`
	Network_userid           *string                             `json:"nuid,omitempty"`
	Userid                   *string                             `json:"uid,omitempty"`
	Domain_sessionidx        *int64                              `json:"vid,string,omitempty"`
	Domain_sessionid         *string                             `json:"sid,omitempty"`
	User_ipaddress           string                              `json:"ip,omitempty"`
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
	Refr_campaign            *string                             `json:"refr_campaign,omitempty"`
	User_fingerprint         *string                             `json:"fp,omitempty"` // deprecated
	Br_cookies               FlexibleBoolField                   `json:"cookie"`
	Br_lang                  *string                             `json:"lang,omitempty"`
	Br_features_pdf          FlexibleBoolField                   `json:"f_pdf"`   // to deprecate
	Br_features_quicktime    FlexibleBoolField                   `json:"f_qt"`    // to deprecate
	Br_features_realplayer   FlexibleBoolField                   `json:"f_realp"` // to deprecate
	Br_features_windowsmedia FlexibleBoolField                   `json:"f_wma"`   // to deprecate
	Br_features_director     FlexibleBoolField                   `json:"f_dir"`   // to deprecate
	Br_features_flash        FlexibleBoolField                   `json:"f_fla"`   // to deprecate
	Br_features_java         FlexibleBoolField                   `json:"f_java"`  // to deprecate
	Br_features_gears        FlexibleBoolField                   `json:"f_gears"` // to deprecate
	Br_features_silverlight  FlexibleBoolField                   `json:"f_ag"`    // to deprecate
	Br_colordepth            *int                                `json:"cd,string,omitempty"`
	Doc_charset              *string                             `json:"cs,omitempty"`
	Doc_size                 *string                             `json:"ds,omitempty"`
	Doc_width                *int                                `json:"doc_width,omitempty"`
	Doc_height               *int                                `json:"doc_height,omitempty"`
	Viewport_size            *string                             `json:"vp,omitempty"`
	Br_viewwidth             *int                                `json:"br_viewwidth,omitempty"`
	Br_viewheight            *int                                `json:"br_viewheight,omitempty"`
	Monitor_resolution       *string                             `json:"res,omitempty"`
	Dvce_screenwidth         *int                                `json:"dvce_screenwidth,omitempty"`
	Dvce_screenheight        *int                                `json:"dvce_screenheight,omitempty"`
	Mac_address              *string                             `json:"mac,omitempty"`
	Contexts                 *Base64EncodedContexts              `json:"cx"`
	Self_describing_event    *Base64EncodedSelfDescribingPayload `json:"ue_px,omitempty"`              // Self Describing Event
	Pp_xoffset_min           *int                                `json:"pp_mix,string,omitempty"`      // Page Ping Event
	Pp_xoffset_max           *int                                `json:"pp_max,string,omitempty"`      // Page Ping Event
	Pp_yoffset_min           *int                                `json:"pp_miy,string,omitempty"`      // Page Ping Event
	Pp_yoffset_max           *int                                `json:"pp_may,string,omitempty"`      // Page Ping Event
	Se_category              *string                             `json:"se_ca,omitempty"`              // Struct Event
	Se_action                *string                             `json:"se_ac,omitempty"`              // Struct Event
	Se_label                 *string                             `json:"se_la,omitempty"`              // Struct Event
	Se_property              *string                             `json:"se_pr,omitempty"`              // Struct Event
	Se_value                 *float64                            `json:"se_va,string,omitempty"`       // Struct Event
	Tr_orderid               *string                             `json:"tr_id,omitempty"`              // Transaction Event
	Tr_affiliation           *string                             `json:"tr_af,omitempty"`              // Transaction Event
	Tr_total                 *float64                            `json:"tr_tt,string,omitempty"`       // Transaction Event
	Tr_tax                   *float64                            `json:"tr_tx,string,omitempty"`       // Transaction Event
	Tr_shipping              *float64                            `json:"tr_sh,string,omitempty"`       // Transaction Event
	Tr_city                  *string                             `json:"tr_ci,omitempty"`              // Transaction Event
	Tr_state                 *string                             `json:"tr_st,omitempty"`              // Transaction Event
	Tr_country               *string                             `json:"tr_co,omitempty"`              // Transaction Event
	Tr_currency              *string                             `json:"tr_cu,omitempty"`              // Transaction Event
	Ti_orderid               *string                             `json:"ti_id,omitempty"`              // Transaction Item Event
	Ti_sku                   *string                             `json:"ti_sk,omitempty"`              // Transaction Item Event
	Ti_name                  *string                             `json:"ti_nm,omitempty"`              // Transaction Item Event
	Ti_category              *string                             `json:"ti_ca,omitempty"`              // Transaction Item Event
	Ti_price                 *float64                            `json:"ti_pr,string,omitempty"`       // Transaction Item Event
	Ti_quantity              *int                                `json:"ti_qu,string,omitempty"`       // Transaction Item Event
	Ti_currency              *string                             `json:"ti_cu,omitempty"`              // Transaction Item Event
	Refr_domain_userid       *string                             `json:"refr_domain_userid,omitempty"` // FIXME! Domain Linker
	Refr_domain_tstamp       time.Time                           `json:"refr_domain_tstamp,omitempty"` // FIXME! Domain Linker
}

type SelfDescribingPayload struct {
	Schema string                 `json:"schema"`
	Data   map[string]interface{} `json:"data"`
}

type Context SelfDescribingPayload

type Base64EncodedContexts []Context

func (c *Base64EncodedContexts) UnmarshalJSON(bytes []byte) error {
	var encodedPayload string
	var contexts []Context
	err := json.Unmarshal(bytes, &encodedPayload)
	decodedPayload, err := b64.RawStdEncoding.DecodeString(encodedPayload)
	if err != nil {
		fmt.Printf("error decoding b64 encoded contexts %s\n", err)
	}
	contextPayload := gjson.Parse(string(decodedPayload))
	for _, pl := range contextPayload.Get("data").Array() {
		context := Context{
			Schema: pl.Get("schema").String(),
			Data:   pl.Get("data").Value().(map[string]interface{}),
		}
		contexts = append(contexts, context)
	}
	*c = contexts
	return nil
}

type Base64EncodedSelfDescribingPayload SelfDescribingPayload

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
		return "page_ping"
	case "pv":
		return "page_view"
	case "se":
		return "struct_event"
	case "ue":
		return "self_describing"
	case "tr":
		return "transaction"
	case "ti":
		return "transaction_item"
	case "ad":
		return "ad_impression"
	}
	return "unknown"
}

type Dimension struct {
	height int
	width  int
}

type PageFields struct {
	scheme   string
	host     string
	port     int
	path     string
	query    string
	fragment string
	medium   string
	source   string
	term     string
	content  string
	campaign string
}
