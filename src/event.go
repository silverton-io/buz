package main

import "time"

type Envelope struct {
	isValid     bool
	parseErrors *[]string
	data        Event
}

type Event struct {
	// Application parameters - https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/snowplow-tracker-protocol/#common-parameters-platform-and-event-independent
	Name_tracker             string                     `json:"name_tracker"`
	Event_vendor             *string                    `json:"event_vendor,omitempty"` // deprecated
	App_id                   string                     `json:"app_id"`
	Platform                 string                     `json:"platform"`
	Etl_tstamp               *time.Time                 `json:"etl_tstamp,omitempty"`
	Dvce_created_tstamp      MillisecondTimestampField  `json:"dvce_created_tstamp"`
	Dvce_sent_tstamp         MillisecondTimestampField  `json:"dvce_sent_tstamp"`
	True_tstamp              *MillisecondTimestampField `json:"true_tstamp,omitempty"`
	Collector_tstamp         time.Time                  `json:"collector_tstamp"`
	Derived_tstamp           *time.Time                 `json:"derived_tstamp,omitempty"`
	Os_timezone              *string                    `json:"os_timezone,omitempty"`
	Event                    EventTypeField             `json:"event"`
	Txn_id                   *string                    `json:"txn_id,omitempty"` // deprecated
	Event_id                 string                     `json:"event_id"`
	Event_fingerprint        string                     `json:"event_fingerprint"`
	Tracker_version          string                     `json:"v_tracker"`
	Collector_version        *string                    `json:"v_collector"`
	Etl_version              *string                    `json:"v_etl,omitempty"`
	Domain_userid            string                     `json:"domain_userid"`
	Network_userid           *string                    `json:"network_userid,omitempty"`
	Userid                   *string                    `json:"user_id,omitempty"`
	Domain_sessionidx        *int64                     `json:"domain_sessionidx,omitempty"`
	Domain_sessionid         *string                    `json:"domain_sessionid,omitempty"`
	User_ipaddress           string                     `json:"user_ipaddress"`
	Page_url                 *string                    `json:"page_url,omitempty"`
	Useragent                *string                    `json:"useragent,omitempty"`
	Page_title               *string                    `json:"page_title,omitempty"`
	Page_referrer            *string                    `json:"page_referrer,omitempty"`
	User_fingerprint         *string                    `json:"user_fingerprint,omitempty"`
	Br_cookies               FlexibleBoolField          `json:"br_cookies"`
	Br_lang                  *string                    `json:"br_lang,omitempty"`
	Br_features_pdf          FlexibleBoolField          `json:"br_features_pdf"`          // deprecate
	Br_features_quicktime    FlexibleBoolField          `json:"br_features_quicktime"`    // deprecate
	Br_features_realplayer   FlexibleBoolField          `json:"br_features_realplayer"`   // deprecate
	Br_features_windowsmedia FlexibleBoolField          `json:"br_features_windowsmedia"` // deprecate
	Br_features_director     FlexibleBoolField          `json:"br_features_director"`     // deprecate
	Br_features_flash        FlexibleBoolField          `json:"br_features_flash"`        // deprecate
	Br_features_java         FlexibleBoolField          `json:"br_features_java"`         // deprecate
	Br_features_gears        FlexibleBoolField          `json:"br_features_gears"`        // deprecate
	Br_features_silverlight  FlexibleBoolField          `json:"br_features_silverlight"`  // deprecate
	Br_colordepth            *int                       `json:"br_colordepth,omitempty"`
	Doc_charset              *string                    `json:"doc_charset,omitempty"`
	Doc_size                 *string                    `json:"doc_size,omitempty"`
	Doc_width                *int                       `json:"doc_width,omitempty"`
	Doc_height               *int                       `json:"doc_height,omitempty"`
	Viewport_size            *string                    `json:"viewport_size,omitempty"`
	Br_viewwidth             *int                       `json:"br_viewwidth,omitempty"`
	Br_viewheight            *int                       `json:"br_viewheight,omitempty"`
	Monitor_resolution       *string                    `json:"monitor_resolution,omitempty"`
	Dvce_screenwidth         *int                       `json:"dvce_screenwidth,omitempty"`
	Dvce_screenheight        *int                       `json:"dvce_screenheight,omitempty"`
	Mac_address              *string                    `json:"mac_address,omitempty"`
	// Contexts                 *[]SelfDescribingPayload            `json:"contexts"`
	Self_describing_event *Base64EncodedSelfDescribingPayload `json:"self_describing_event,omitempty"` // Self Describing Event
	Pp_xoffset_min        *int                                `json:"pp_xoffset_min,omitempty"`        // Page Ping Event
	Pp_xoffset_max        *int                                `json:"pp_xoffset_max,omitempty"`        // Page Ping Event
	Pp_yoffset_min        *int                                `json:"pp_yoffset_min,omitempty"`        // Page Ping Event
	Pp_yoffset_max        *int                                `json:"pp_yoffset_max,omitempty"`        // Page Ping Event
	Se_category           *string                             `json:"se_category,omitempty"`           // Struct Event
	Se_action             *string                             `json:"se_action,omitempty"`             // Struct Event
	Se_label              *string                             `json:"se_label,omitempty"`              // Struct Event
	Se_property           *string                             `json:"se_property,omitempty"`           // Struct Event
	Se_value              *float64                            `json:"se_value,omitempty"`              // Struct Event
	Tr_orderid            *string                             `json:"tr_orderid,omitempty"`            // Transaction Event
	Tr_affiliation        *string                             `json:"tr_affiliation,omitempty"`        // Transaction Event
	Tr_total              *float64                            `json:"tr_total,omitempty"`              // Transaction Event
	Tr_tax                *float64                            `json:"tr_tax,omitempty"`                // Transaction Event
	Tr_shipping           *float64                            `json:"tr_shipping,omitempty"`           // Transaction Event
	Tr_city               *string                             `json:"tr_city,omitempty"`               // Transaction Event
	Tr_state              *string                             `json:"tr_state,omitempty"`              // Transaction Event
	Tr_country            *string                             `json:"tr_country,omitempty"`            // Transaction Event
	Tr_currency           *string                             `json:"tr_currency,omitempty"`           // Transaction Event
	Ti_orderid            *string                             `json:"ti_orderid,omitempty"`            // Transaction Item Event
	Ti_sku                *string                             `json:"ti_sku,omitempty"`                // Transaction Item Event
	Ti_name               *string                             `json:"ti_name,omitempty"`               // Transaction Item Event
	Ti_category           *string                             `json:"ti_category,omitempty"`           // Transaction Item Event
	Ti_price              *float64                            `json:"ti_price,string,omitempty"`       // Transaction Item Event
	Ti_quantity           *int                                `json:"ti_quantity,omitempty"`           // Transaction Item Event
	Ti_currency           *string                             `json:"ti_currency,omitempty"`           // Transaction Item Event
	Refr_domain_userid    *string                             `json:"refr_domain_userid,omitempty"`    // Domain Linking FIXME!!
	Refr_domain_tstamp    time.Time                           `json:"refr_domain_tstamp,omitempty"`    // Domain Linking FIXME!!
}

type ShortenedEvent struct {
	Name_tracker             string                     `json:"tna"`
	Event_vendor             *string                    `json:"evn,omitempty"` // deprecated
	App_id                   string                     `json:"aid"`
	Platform                 string                     `json:"p"`
	Etl_tstamp               *time.Time                 `json:"etl_tstamp,omitempty"` // not in the tracker protocol
	Dvce_created_tstamp      MillisecondTimestampField  `json:"dtm"`
	Dvce_sent_tstamp         MillisecondTimestampField  `json:"stm"`
	True_tstamp              *MillisecondTimestampField `json:"ttm,omitempty"`
	Collector_tstamp         time.Time                  `json:"collector_tstamp"`         // not in the tracker protocol
	Derived_tstamp           *time.Time                 `json:"derived_tstamp,omitempty"` // not in the tracker protocol
	Os_timezone              *string                    `json:"tz,omitempty"`
	Event                    EventTypeField             `json:"e"`
	Txn_id                   *string                    `json:"tid,omitempty"` // deprecated
	Event_id                 string                     `json:"eid"`
	Event_fingerprint        string                     `json:"event_fingerprint"`
	Tracker_version          string                     `json:"tv"`
	Collector_version        *string                    `json:"v_collector"`
	Etl_version              *string                    `json:"v_etl,omitempty"`
	Domain_userid            string                     `json:"duid"`
	Network_userid           *string                    `json:"nuid,omitempty"`
	Userid                   *string                    `json:"uid,omitempty"`
	Domain_sessionidx        *int64                     `json:"vid,string,omitempty"`
	Domain_sessionid         *string                    `json:"sid,omitempty"`
	User_ipaddress           string                     `json:"ip,omitempty"`
	Page_url                 *string                    `json:"url,omitempty"`
	Useragent                *string                    `json:"ua,omitempty"`
	Page_title               *string                    `json:"page,omitempty"`
	Page_referrer            *string                    `json:"refr,omitempty"`
	User_fingerprint         *string                    `json:"fp,omitempty"` // deprecated
	Br_cookies               FlexibleBoolField          `json:"cookie"`
	Br_lang                  *string                    `json:"lang,omitempty"`
	Br_features_pdf          FlexibleBoolField          `json:"f_pdf"`   // deprecate
	Br_features_quicktime    FlexibleBoolField          `json:"f_qt"`    // deprecate
	Br_features_realplayer   FlexibleBoolField          `json:"f_realp"` // deprecate
	Br_features_windowsmedia FlexibleBoolField          `json:"f_wma"`   // deprecate
	Br_features_director     FlexibleBoolField          `json:"f_dir"`   // deprecate
	Br_features_flash        FlexibleBoolField          `json:"f_fla"`   // deprecate
	Br_features_java         FlexibleBoolField          `json:"f_java"`  // deprecate
	Br_features_gears        FlexibleBoolField          `json:"f_gears"` // deprecate
	Br_features_silverlight  FlexibleBoolField          `json:"f_ag"`    // deprecate
	Br_colordepth            *int                       `json:"cd,string,omitempty"`
	Doc_charset              *string                    `json:"cs,omitempty"`
	Doc_size                 *string                    `json:"ds,omitempty"`
	Doc_width                *int                       `json:"doc_width,omitempty"`
	Doc_height               *int                       `json:"doc_height,omitempty"`
	Viewport_size            *string                    `json:"vp,omitempty"`
	Br_viewwidth             *int                       `json:"br_viewwidth,omitempty"`
	Br_viewheight            *int                       `json:"br_viewheight,omitempty"`
	Monitor_resolution       *string                    `json:"res,omitempty"`
	Dvce_screenwidth         *int                       `json:"dvce_screenwidth,omitempty"`
	Dvce_screenheight        *int                       `json:"dvce_screenheight,omitempty"`
	Mac_address              *string                    `json:"mac,omitempty"`
	// Contexts                 *[]SelfDescribingPayload            `json:"cx"`
	Self_describing_event *Base64EncodedSelfDescribingPayload `json:"ue_px,omitempty"`              // Self Describing Event
	Pp_xoffset_min        *int                                `json:"pp_mix,string,omitempty"`      // Page Ping Event
	Pp_xoffset_max        *int                                `json:"pp_max,string,omitempty"`      // Page Ping Event
	Pp_yoffset_min        *int                                `json:"pp_miy,string,omitempty"`      // Page Ping Event
	Pp_yoffset_max        *int                                `json:"pp_may,string,omitempty"`      // Page Ping Event
	Se_category           *string                             `json:"se_ca,omitempty"`              // Struct Event
	Se_action             *string                             `json:"se_ac,omitempty"`              // Struct Event
	Se_label              *string                             `json:"se_la,omitempty"`              // Struct Event
	Se_property           *string                             `json:"se_pr,omitempty"`              // Struct Event
	Se_value              *float64                            `json:"se_va,string,omitempty"`       // Struct Event
	Tr_orderid            *string                             `json:"tr_id,omitempty"`              // Transaction Event
	Tr_affiliation        *string                             `json:"tr_af,omitempty"`              // Transaction Event
	Tr_total              *float64                            `json:"tr_tt,string,omitempty"`       // Transaction Event
	Tr_tax                *float64                            `json:"tr_tx,string,omitempty"`       // Transaction Event
	Tr_shipping           *float64                            `json:"tr_sh,string,omitempty"`       // Transaction Event
	Tr_city               *string                             `json:"tr_ci,omitempty"`              // Transaction Event
	Tr_state              *string                             `json:"tr_st,omitempty"`              // Transaction Event
	Tr_country            *string                             `json:"tr_co,omitempty"`              // Transaction Event
	Tr_currency           *string                             `json:"tr_cu,omitempty"`              // Transaction Event
	Ti_orderid            *string                             `json:"ti_id,omitempty"`              // Transaction Item Event
	Ti_sku                *string                             `json:"ti_sk,omitempty"`              // Transaction Item Event
	Ti_name               *string                             `json:"ti_nm,omitempty"`              // Transaction Item Event
	Ti_category           *string                             `json:"ti_ca,omitempty"`              // Transaction Item Event
	Ti_price              *float64                            `json:"ti_pr,string,omitempty"`       // Transaction Item Event
	Ti_quantity           *int                                `json:"ti_qu,string,omitempty"`       // Transaction Item Event
	Ti_currency           *string                             `json:"ti_cu,omitempty"`              // Transaction Item Event
	Refr_domain_userid    *string                             `json:"refr_domain_userid,omitempty"` // Domain Linking FIXME!!
	Refr_domain_tstamp    time.Time                           `json:"refr_domain_tstamp,omitempty"` // Domain Linking FIXME!!
}
