package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func msStringToTime(ms string) time.Time {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		fmt.Println("FIXME!!")
	}
	return time.Unix(0, msInt*int64(time.Millisecond))
}

func stringToInt(val string) int {
	// FIXME! Handle this entire thing better so we don't swallow params that don't exist
	i, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("FIXME! stringToInt")
	}
	return i
}

func stringToBool(val string) bool {
	if val == "1" {
		return true
	} else {
		return false
	}
}

func stringToWidth(val string) int {
	split := strings.Split(val, "x")
	return stringToInt(split[0])
}

func stringToHeight(val string) int {
	split := strings.Split(val, "x")
	return stringToInt(split[1])
}

func stringToFloat64(val string) float64 {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Fatal(err) // FIXME!
	}
	return f
}

func b64ToMap(encodedJson string) map[string]interface{} {
	var decodedJSON map[string]interface{}
	bytes, err := b64.RawStdEncoding.DecodeString(encodedJson)
	if err != nil {
		log.Fatal(err)
		fmt.Println("FIXME!! b64 string could not be decoded")
	}
	err = json.Unmarshal(bytes, &decodedJSON)
	return decodedJSON
}

// func getJsonEncodedParam(param string) {

// }

// func getTimeParam(param string) time.Time {

// }

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
