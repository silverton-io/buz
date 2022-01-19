package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Event struct {
	// Application parameters - https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/snowplow-tracker-protocol/#common-parameters-platform-and-event-independent
	name_tracker string `json:"name_tracker"` // tna
	event_vendor string `json:"event_vendor"` // evn - deprecated
	app_id       string `json:"app_id"`       // aid
	platform     string `json:"platform"`     // p
	// Date/time parameters
	dvce_created_tstamp int    `json:"dvce_created_tstamp"` // dtm
	dvce_sent_tstamp    int    `json:"dvce_sent_tstamp"`    // stm
	true_tstamp         int    `json:"true_tstamp"`         // ttm
	os_timezone         string `json:"os_timezone"`         // tz
	// Event/transaction parameters
	event    string `json:"event"`    // e
	txn_id   string `json:"txn_id"`   // tid - deprecated
	event_id string `json:"event_id"` // eid
	// Tracker version
	v_tracker string `json:"v_tracker"` // tv
	// User params
	domain_userid     string `json:"domain_userid"`     // duid
	network_userid    string `json:"network_userid"`    // nuid or tnuid OR VALUE OF COOKIE
	user_id           string `json:"user_id"`           // uid
	domain_sessionidx int    `json:"domain_sessionidx"` // vid
	domain_sessionid  string `json:"domain_sessionid"`  // sid
	user_ipaddress    string `json:"user_ipaddress"`    // ip
	// dvce_screenheight int    `json:"dvce_screenheight"` // res - FIXME! Implement wxh
	// dvce_screenwidth  int    `json:"dvce_screenwidth"` // res - FIXME! Implement wxh
	page_url         string `json:"page_url"`         // url
	useragent        string `json:"useragent"`        // ua
	page_title       string `json:"page_title"`       // page
	user_fingerprint string `json:"user_fingerprint"` // fp
	// connection_type string // NOT IMPLEMENTED - UNUSED
	br_cookies             bool   `json:"br_cookies"`             // cookie
	br_lang                string `json:"br_lang"`                // lang
	br_features_pdf        bool   `json:"br_features_pdf"`        // f_pdf
	br_features_quicktime  bool   `json:"br_features_quicktime"`  // f_qt
	br_features_realplayer bool   `json:"br_features_realplayer"` // f_realp
	br_features_wma        bool   `json:"br_features_wma"`        // f_wma

}

func toInt(val string) int {
	// FIXME! Handle this entire thing better so we don't swallow params that don't exist
	i, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("Cannot convert ", val, " to int") // FIXME! Log or deadletter or something
	}
	return i
}

func toBool(val string) bool {
	// FIXME!
	if val == "0" {
		return false
	} else {
		return true
	}

}

func parseQueryParams(c *gin.Context) Event {
	event := Event{
		name_tracker:           c.Query("tna"),
		event_vendor:           c.Query("evn"),
		app_id:                 c.Query("aid"),
		platform:               c.Query("p"),
		dvce_created_tstamp:    toInt(c.Query("dtm")),
		dvce_sent_tstamp:       toInt(c.Query("stm")),
		true_tstamp:            toInt(c.Query("ttm")),
		os_timezone:            c.Query("tz"),
		event:                  c.Query("e"),
		txn_id:                 c.Query("tid"),
		event_id:               c.Query("eid"),
		v_tracker:              c.Query("tv"),
		domain_userid:          c.Query("duid"),
		network_userid:         c.Query("nuid"), // FIXME! Also support tnuid
		user_id:                c.Query("uid"),
		domain_sessionidx:      toInt(c.Query("vid")),
		domain_sessionid:       c.Query("sid"),
		user_ipaddress:         c.Query("ip"),
		page_url:               c.Query("url"),
		useragent:              c.Query("ua"),
		page_title:             c.Query("page"),
		user_fingerprint:       c.Query("fp"),
		br_cookies:             toBool(c.Query("cookie")),
		br_lang:                c.Query("lang"),
		br_features_pdf:        toBool(c.Query("f_pdf")),
		br_features_quicktime:  toBool(c.Query("f_qt")),
		br_features_realplayer: toBool(c.Query("f_realp")),
		br_features_wma:        toBool(c.Query("f_wma")),
	}
	return event
}

func HandleGet(c *gin.Context) {
	// FIXME! Parse Query Params
	event := parseQueryParams(c)
	fmt.Printf("%+v\n", event)
	c.JSON(200, gin.H{
		"message": "received",
	})
}

func HandleRedirect(c *gin.Context) {
	event := parseQueryParams(c)
	fmt.Printf("%+v\n", event)
	// FIXME! Parse Query Parameters Here
	redirectUrl, _ := c.GetQuery("u")
	c.Redirect(302, redirectUrl)
}
