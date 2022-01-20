package main

import "time"

type SelfDescribingEnvelope struct {
	schema string
	data   map[string]interface{}
}

type Event struct {
	// Application parameters - https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/snowplow-tracker-protocol/#common-parameters-platform-and-event-independent
	// FIXME'S! ------------------------------------------> dvce_screenheight, dvce_screenwidth, br_viewport stuff, res
	Name_tracker             string                 `json:"name_tracker"`
	Event_vendor             string                 `json:"event_vendor"` // deprecated
	App_id                   string                 `json:"app_id"`
	Platform                 string                 `json:"platform"`
	Dvce_created_tstamp      time.Time              `json:"dvce_created_tstamp"`
	Dvce_sent_tstamp         time.Time              `json:"dvce_sent_tstamp"`
	Collector_tstamp         time.Time              `json:"collector_tstamp"`
	Derived_tstamp           time.Time              `json:"derived_tstamp"`
	True_tstamp              time.Time              `json:"true_tstamp"`
	Os_timezone              string                 `json:"os_timezone"`
	Event                    string                 `json:"event"`
	Txn_id                   string                 `json:"txn_id"` // deprecated
	Event_id                 string                 `json:"event_id"`
	V_tracker                string                 `json:"v_tracker"`
	Domain_userid            string                 `json:"domain_userid"`
	Network_userid           string                 `json:"network_userid"`
	User_id                  string                 `json:"user_id"`
	Domain_sessionidx        int                    `json:"domain_sessionidx"`
	Domain_sessionid         string                 `json:"domain_sessionid"`
	User_ipaddress           string                 `json:"user_ipaddress"`
	Page_url                 string                 `json:"page_url"`
	Useragent                string                 `json:"useragent"`
	Page_title               string                 `json:"page_title"`
	Page_referrer            string                 `json:"page_referrer"`
	User_fingerprint         string                 `json:"user_fingerprint"`
	Br_cookies               bool                   `json:"br_cookies"`
	Br_lang                  string                 `json:"br_lang"`
	Br_features_pdf          bool                   `json:"br_features_pdf"`
	Br_features_quicktime    bool                   `json:"br_features_quicktime"`
	Br_features_realplayer   bool                   `json:"br_features_realplayer"`
	Br_features_windowsmedia bool                   `json:"br_features_windowsmedia"`
	Br_features_director     bool                   `json:"br_features_director"`
	Br_features_flash        bool                   `json:"br_features_flash"`
	Br_features_java         bool                   `json:"br_features_java"`
	Br_features_gears        bool                   `json:"br_features_gears"`
	Br_features_silverlight  bool                   `json:"br_features_silverlight"`
	Br_colordepth            int                    `json:"br_colordepth"`
	Doc_width                int                    `json:"doc_width"`  // ds position 1
	Doc_height               int                    `json:"doc_height"` // ds position 2
	Doc_charset              string                 `json:"doc_charset"`
	Br_viewwidth             int                    `json:"br_viewwidth"`      // vp position 1
	Br_viewheight            int                    `json:"br_viewheight"`     // vp position 2
	Dvce_screenwidth         int                    `json:"dvce_screenwidth"`  // res position 1
	Dvce_screenheight        int                    `json:"dvce_screenheight"` // res position 2
	Mac_address              string                 `json:"mac_address"`
	Contexts                 map[string]interface{} `json:"context"` // either co or cx
	// Page ping
	Pp_xoffset_min int `json:"pp_xoffset_min"`
	Pp_xoffset_max int `json:"pp_xoffset_max"`
	Pp_yoffset_min int `json:"pp_yoffset_min"`
	Pp_yoffset_max int `json:"pp_yoffset_max"`
	// Transaction
	Tr_orderid     string  `json:"tr_orderid"`
	Tr_affiliation string  `json:"tr_affiliation"`
	Tr_total       float64 `json:"tr_total"`
	Tr_tax         float64 `json:"tr_tax"`
	Tr_shipping    float64 `json:"tr_shipping"`
	Tr_city        string  `json:"tr_city"`
	Tr_state       string  `json:"tr_state"`
	Tr_country     string  `json:"tr_country"`
	Tr_currency    string  `json:"tr_currency"`
	// Transaction item
	Ti_orderid  string  `json:"ti_orderid"`
	Ti_sku      string  `json:"ti_sku"`
	Ti_name     string  `json:"ti_name"`
	Ti_category string  `json:"ti_category"`
	Ti_price    float64 `json:"ti_price"`
	Ti_quantity int     `json:"ti_quantity"`
	Ti_currency string  `json:"ti_currency"`
	// Struct
	Se_category string  `json:"se_category"`
	Se_action   string  `json:"se_action"`
	Se_label    string  `json:"se_label"`
	Se_property string  `json:"se_property"`
	Se_value    float64 `json:"se_value"`
	// Self Describing
	Self_describing_event map[string]interface{} `json:"self_describing_event"` // either ue_px or
}
