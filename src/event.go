package main

import "time"

type Envelope struct {
	isValid bool
	errors  *[]string
	data    Event
}

type Event struct {
	// Application parameters - https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/snowplow-tracker-protocol/#common-parameters-platform-and-event-independent
	Name_tracker             string                     `json:"name_tracker"`
	Event_vendor             *string                    `json:"event_vendor"` // deprecated
	App_id                   string                     `json:"app_id"`
	Platform                 string                     `json:"platform"`
	Dvce_created_tstamp      MillisecondTimestampField  `json:"dvce_created_tstamp"`
	Dvce_sent_tstamp         MillisecondTimestampField  `json:"dvce_sent_tstamp"`
	True_tstamp              *MillisecondTimestampField `json:"true_tstamp"`
	Collector_tstamp         time.Time                  `json:"collector_tstamp"`
	Derived_tstamp           *time.Time                 `json:"derived_tstamp"`
	Os_timezone              string                     `json:"os_timezone"`
	Event                    string                     `json:"event"`
	Txn_id                   *string                    `json:"txn_id"` // deprecated
	Event_id                 string                     `json:"event_id"`
	Tracker_version          string                     `json:"v_tracker"`
	Domain_userid            string                     `json:"domain_userid"`
	Network_userid           string                     `json:"network_userid"`
	Userid                   *string                    `json:"user_id"`
	Domain_sessionidx        int64                      `json:"domain_sessionidx"`
	Domain_sessionid         string                     `json:"domain_sessionid"`
	User_ipaddress           string                     `json:"user_ipaddress"`
	Page_url                 string                     `json:"page_url"`
	Useragent                string                     `json:"useragent"`
	Page_title               string                     `json:"page_title"`
	Page_referrer            string                     `json:"page_referrer"`
	User_fingerprint         string                     `json:"user_fingerprint"`
	Br_cookies               bool                       `json:"br_cookies"`
	Br_lang                  string                     `json:"br_lang"`
	Br_features_pdf          bool                       `json:"br_features_pdf"`          // deprecate
	Br_features_quicktime    bool                       `json:"br_features_quicktime"`    // deprecate
	Br_features_realplayer   bool                       `json:"br_features_realplayer"`   // deprecate
	Br_features_windowsmedia bool                       `json:"br_features_windowsmedia"` // deprecate
	Br_features_director     bool                       `json:"br_features_director"`     // deprecate
	Br_features_flash        bool                       `json:"br_features_flash"`        // deprecate
	Br_features_java         bool                       `json:"br_features_java"`         // deprecate
	Br_features_gears        bool                       `json:"br_features_gears"`        // deprecate
	Br_features_silverlight  bool                       `json:"br_features_silverlight"`  // deprecate
	Br_colordepth            int                        `json:"br_colordepth"`
	Doc_charset              string                     `json:"doc_charset"`
	Mac_address              string                     `json:"mac_address"`
	// More here
	Contexts              []SelfDescribingPayload            `json:"contexts"`
	Self_describing_event Base64EncodedSelfDescribingPayload `json:"self_describing_event"`
}

type ShortenedEvent struct {
	Name_tracker             string                             `json:"tna"`
	Event_vendor             *string                            `json:"evn"` // deprecated
	App_id                   string                             `json:"aid"`
	Platform                 string                             `json:"p"`
	Dvce_created_tstamp      MillisecondTimestampField          `json:"dtm"`
	Dvce_sent_tstamp         MillisecondTimestampField          `json:"stm"`
	True_tstamp              *MillisecondTimestampField         `json:"ttm"`
	Collector_tstamp         time.Time                          `json:"collector_tstamp"`
	Derived_tstamp           *time.Time                         `json:"drtm"` // not actually in the tracker protocol
	Os_timezone              string                             `json:"tz"`
	Event                    string                             `json:"e"`
	Txn_id                   *string                            `json:"tid"` // deprecated
	Event_id                 string                             `json:"eid"`
	Tracker_version          string                             `json:"tv"`
	Domain_userid            string                             `json:"duid"`
	Network_userid           string                             `json:"nuid"`
	Userid                   *string                            `json:"uid"`
	Domain_sessionidx        int64                              `json:"vid,string"`
	Domain_sessionid         string                             `json:"sid"`
	User_ipaddress           string                             `json:"ip"`
	Page_url                 string                             `json:"url"`
	Useragent                string                             `json:"ua"`
	Page_title               string                             `json:"page"`
	Page_referrer            string                             `json:"refr"`
	User_fingerprint         string                             `json:"fp"`
	Br_cookies               bool                               `json:"cookie,string"`
	Br_lang                  string                             `json:"lang"`
	Br_features_pdf          bool                               `json:"f_pdf,string"`   // deprecate
	Br_features_quicktime    bool                               `json:"f_qt,string"`    // deprecate
	Br_features_realplayer   bool                               `json:"f_realp,string"` // deprecate
	Br_features_windowsmedia bool                               `json:"f_wma,string"`   // deprecate
	Br_features_director     bool                               `json:"f_dir,string"`   // deprecate
	Br_features_flash        bool                               `json:"f_fla,string"`   // deprecate
	Br_features_java         bool                               `json:"f_java,string"`  // deprecate
	Br_features_gears        bool                               `json:"f_gears,string"` // deprecate
	Br_features_silverlight  bool                               `json:"f_ag,string"`    // deprecate
	Br_colordepth            int                                `json:"cd,string"`
	Doc_charset              string                             `json:"cs"`
	Mac_address              string                             `json:"mac"`
	Contexts                 []SelfDescribingPayload            `json:"cx"`
	Self_describing_event    Base64EncodedSelfDescribingPayload `json:"ue_px"`
}
