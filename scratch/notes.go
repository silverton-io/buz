
Doc_width                int                    `json:"doc_width"`  // ds position 1
Doc_height               int                    `json:"doc_height"` // ds position 2
Br_viewwidth             int                    `json:"br_viewwidth"`      // vp position 1
Br_viewheight            int                    `json:"br_viewheight"`     // vp position 2
Dvce_screenwidth         int                    `json:"dvce_screenwidth"`  // res position 1
Dvce_screenheight        int                    `json:"dvce_screenheight"` // res position 2
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



Doc_width:                stringToWidth(e.Get("ds").String()),
Doc_height:               stringToHeight(e.Get("ds").String()),
Br_viewwidth:             stringToWidth(e.Get("vp").String()),
Br_viewheight:            stringToWidth(e.Get("vp").String()),
Dvce_screenwidth:         stringToWidth(e.Get("res").String()),
Dvce_screenheight:        stringToHeight(e.Get("res").String()),
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