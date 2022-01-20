package main

// func getParamValue(param string, c *gin.Context) string {
// 	if c.Request.Method == "GET" {
// 		return c.Query(param)
// 	} else if c.Request.Method == "POST" {

// 	}
// }

// func getB64Param(param string) {

// }

// func getJsonEncodedParam(param string) {

// }

// func getTimeParam(param string) time.Time {

// }

// func getFloat64Param(param string) float64 {
// 	f, err := strconv.ParseFloat(param, 64)
// 	if err != nil {
// 		fmt.Println("Cannot convert val ", param)
// 	}
// 	return f
// }

// func getIntParam(param string) int {
// 	// FIXME! Handle this entire thing better so we don't swallow params that don't exist
// 	i, err := strconv.Atoi(param)
// 	if err != nil {
// 		fmt.Println("Cannot convert ", param, " to int") // FIXME! Log or deadletter or something
// 	}
// 	return i
// }

// func getBoolParam(param string, c *gin.Context) bool {

// }

// func getEventType(param string) string {
// 	switch param {
// 	case "pp":
// 		return "page_ping"
// 	case "pv":
// 		return "page_view"
// 	case "se":
// 		return "struct_event"
// 	case "ue":
// 		return "self_describing"
// 	case "tr":
// 		return "transaction"
// 	case "ti":
// 		return "transaction_item"
// 	case "ad":
// 		return "ad_impression"
// 	}
// 	return "unknown"
// }
