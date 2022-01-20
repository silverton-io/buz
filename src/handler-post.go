package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func HandlePost(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	// payloadSchema := gjson.GetBytes(body, "schema")
	payloadData := gjson.GetBytes(body, "data")
	for _, e := range payloadData.Array() {
		event := Event{
			Name_tracker:             e.Get("tna").String(),
			Event_vendor:             e.Get("evn").String(),
			App_id:                   e.Get("aid").String(),
			Platform:                 e.Get("p").String(),
			Dvce_created_tstamp:      msStringToTime(e.Get("dtm").String()),
			Dvce_sent_tstamp:         msStringToTime(e.Get("stm").String()),
			True_tstamp:              msStringToTime(e.Get("ttm").String()),
			Collector_tstamp:         time.Now(),
			Os_timezone:              e.Get("tz").String(),
			Event:                    getEventType(e.Get("e").String()),
			Txn_id:                   e.Get("tid").String(),
			Event_id:                 e.Get("eid").String(),
			V_tracker:                e.Get("tv").String(),
			Domain_userid:            e.Get("duid").String(),
			Network_userid:           e.Get("nuid").String(),
			User_id:                  e.Get("uid").String(),
			Domain_sessionidx:        stringToInt(e.Get("vid").String()),
			Domain_sessionid:         e.Get("sid").String(),
			User_ipaddress:           e.Get("ip").String(),
			Page_url:                 e.Get("url").String(),
			Useragent:                e.Get("ua").String(), // "or from request context"
			Page_title:               e.Get("page").String(),
			Page_referrer:            e.Get("refr").String(),
			User_fingerprint:         e.Get("fp").String(),
			Br_cookies:               stringToBool(e.Get("cookie").String()),
			Br_lang:                  e.Get("lang").String(),
			Br_features_pdf:          stringToBool(e.Get("f_pdf").String()),
			Br_features_quicktime:    stringToBool(e.Get("f_qt").String()),
			Br_features_realplayer:   stringToBool(e.Get("f_realp").String()),
			Br_features_windowsmedia: stringToBool(e.Get("f_wma").String()),
			Br_features_director:     stringToBool(e.Get("f_dir").String()),
			Br_features_flash:        stringToBool(e.Get("f_fla").String()),
			Br_features_java:         stringToBool(e.Get("f_java").String()),
			Br_features_gears:        stringToBool(e.Get("f_gears").String()),
			Br_features_silverlight:  stringToBool(e.Get("f_ag").String()),
			Br_colordepth:            stringToInt(e.Get("cd").String()),
			Doc_width:                stringToWidth(e.Get("ds").String()),
			Doc_height:               stringToHeight(e.Get("ds").String()),
			Doc_charset:              e.Get("cs").String(),
			Br_viewwidth:             stringToWidth(e.Get("vp").String()),
			Br_viewheight:            stringToWidth(e.Get("vp").String()),
			Dvce_screenwidth:         stringToWidth(e.Get("res").String()),
			Dvce_screenheight:        stringToHeight(e.Get("res").String()),
			Mac_address:              e.Get("mac").String(),
			Contexts:                 b64ToMap(e.Get("cx").String()),
			Pp_xoffset_min:           stringToInt(e.Get("pp_mix").String()),
			Pp_xoffset_max:           stringToInt(e.Get("pp_max").String()),
			Pp_yoffset_min:           stringToInt(e.Get("pp_miy").String()),
			Pp_yoffset_max:           stringToInt(e.Get("pp_may").String()),

			Tr_orderid:     e.Get("tr_id").String(),
			Tr_affiliation: e.Get("tr_af").String(),
			Tr_total:       stringToFloat64(e.Get("tr_tt").String()),
			Tr_tax:         stringToFloat64(e.Get("tr_tx").String()),
			Tr_shipping:    stringToFloat64(e.Get("tr_sh").String()),
			Tr_city:        e.Get("tr_ci").String(),
			Tr_state:       e.Get("tr_st").String(),
			Tr_country:     e.Get("tr_co").String(),
			Tr_currency:    e.Get("tr_cu").String(),

			Ti_orderid:  e.Get("ti_id").String(),
			Ti_sku:      e.Get("ti_sku").String(),
			Ti_name:     e.Get("ti_nm").String(),
			Ti_category: e.Get("ti_ca").String(),
			Ti_price:    stringToFloat64(e.Get("ti_pr").String()),
			Ti_quantity: stringToInt(e.Get("ti_qu").String()),
			Ti_currency: e.Get("ti_cu").String(),

			Se_category: e.Get("se_ca").String(),
			Se_action:   e.Get("se_ac").String(),
			Se_label:    e.Get("se_la").String(),
			Se_property: e.Get("se_pr").String(),
			Se_value:    stringToFloat64(e.Get("se_va").String()),

			Self_describing_event: b64ToMap(e.Get("ue_px").String()),
		}
		// event.Network_userid = c.Get("identity")
		data, _ := json.Marshal(event)
		fmt.Printf("%s\n", data)
	}
	// if err != nil {
	// 	log.Fatal(err)
	// 	fmt.Println("Something bad happened")
	// }
	// for _, i := range shortenedEvents {
	// 	fmt.Println(i)
	// }
	// err := json.Unmarshal(data, &events)

	c.JSON(200, gin.H{
		"message": "received",
	})
}
