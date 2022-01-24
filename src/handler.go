package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func mapParams(c *gin.Context) map[string]interface{} {
	mappedParams := make(map[string]interface{})
	params := c.Request.URL.Query()
	for k, v := range params {
		mappedParams[k] = v[0]
	}
	return mappedParams
}

func setEventCollectorMetadataFields(e *Event) {
	fmt.Println(Version)
	e.Collector_version = &Version
}

func parseWidthHeight(dimensionString string) (Dimension, error) {
	dim := strings.Split(dimensionString, "x")
	wS, hS := dim[0], dim[1]
	width, err := strconv.Atoi(wS)
	height, err := strconv.Atoi(hS)
	if err != nil {
		return Dimension{}, err
	}
	return Dimension{
		width:  width,
		height: height,
	}, nil
}

func setEventFieldsFromRequest(c *gin.Context, e *Event) {
	nuid := c.GetString("identity")
	ip, _ := c.RemoteIP()
	useragent := c.Request.UserAgent()
	e.Network_userid = &nuid
	// NOTE!! Intentionally ignore the ip and useragent set via query param.
	e.User_ipaddress = ip.String()
	e.Useragent = &useragent
	e.Collector_tstamp = time.Now()
}

func setEventWidthHeightFields(e *Event) {
	// Doc
	docDimension, _ := parseWidthHeight(*e.Doc_size)
	e.Doc_width, e.Doc_height = &docDimension.width, &docDimension.height
	// Viewport
	vpDimension, _ := parseWidthHeight(*e.Viewport_size)
	e.Br_viewwidth, e.Br_viewheight = &vpDimension.width, &vpDimension.height
	// Screen
	monDimension, _ := parseWidthHeight(*e.Monitor_resolution)
	e.Dvce_screenwidth, e.Dvce_screenheight = &monDimension.width, &monDimension.height
}

func getPageFieldsFromUrl(rawUrl string) (PageFields, error) {
	parsedUrl, err := url.Parse(rawUrl)
	unescapedQry, err := url.QueryUnescape(parsedUrl.RawQuery)
	if err != nil {
		return PageFields{}, err
	}
	// medium := parsedUrl.Query()["utm_medium"][0]
	// source := parsedUrl.Query()["utm_source"][0]
	// term := parsedUrl.Query()["utm_term"][0]
	content := parsedUrl.Query()["utm_content"][0]
	campaign := parsedUrl.Query()["utm_campaign"][0]
	return PageFields{
		scheme: parsedUrl.Scheme,
		host:   parsedUrl.Host,
		// port: FIXME!!!
		path:     parsedUrl.Path,
		query:    unescapedQry,
		fragment: parsedUrl.Fragment,
		// medium:   medium,
		// source:   source,
		// term:     term,
		content:  content,
		campaign: campaign,
	}, nil
}

func setPageFields(e *Event) {
	pageFields, err := getPageFieldsFromUrl(*e.Page_url)
	if err != nil {
		fmt.Printf("error setting page fields %s\n", err)
	}
	e.Page_urlscheme = &pageFields.scheme
	e.Page_urlhost = &pageFields.host
	e.Page_urlpath = &pageFields.path
	e.Page_urlquery = &pageFields.query
	e.Page_urlfragment = &pageFields.fragment
	// e.Mkt_medium = &pageFields.medium
	// e.Mkt_source = &pageFields.source
	// e.Mkt_term = &pageFields.term
	e.Mkt_content = &pageFields.content
	e.Mkt_campaign = &pageFields.campaign
}

func buildEventFromMappedParams(c *gin.Context, e map[string]interface{}) Event {
	body, err := json.Marshal(e)
	if err != nil {
		fmt.Println(err)
	}
	shortenedEvent := ShortenedEvent{}
	err = json.Unmarshal(body, &shortenedEvent)
	if err != nil {
		fmt.Printf("error unmarshalling to shortened event %s", err)
	}
	event := Event(shortenedEvent)
	setEventCollectorMetadataFields(&event)
	setEventFieldsFromRequest(c, &event)
	setEventWidthHeightFields(&event)
	setPageFields(&event)
	prettyPrint(event)
	return event
}

func HandleRedirect(c *gin.Context) {
	mappedParams := mapParams(c)
	buildEventFromMappedParams(c, mappedParams)
	// fmt.Printf("\n%+v\n", se)
	// e, _ := json.Marshal(se)
	redirectUrl, _ := c.GetQuery("u")
	c.Redirect(302, redirectUrl)
}

func HandleGet(c *gin.Context) {
	// Parse query parameters to map[string]interface{}
	mappedParams := mapParams(c)
	se := buildEventFromMappedParams(c, mappedParams)
	// fmt.Printf("\n%+v\n", se)
	// e, _ := json.Marshal(se)
	c.JSON(200, se)
}

// func HandlePost(c *gin.Context) {
// 	body, err := ioutil.ReadAll(c.Request.Body)
// 	if err != nil {
// 		fmt.Println("FIXME! HANDLE ERRORS WHEN READING REQUEST BODY")
// 	}

// 	payloadData := gjson.GetBytes(body, "data")
// 	for _, e := range payloadData.Array() {
// 		// event := buildEvent(c, e)
// 		// fmt.Printf("\n%+v\n", event)
// 		fmt.Println(e)
// 	}
// 	c.JSON(200, AllOk)
// }
