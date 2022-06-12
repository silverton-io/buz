package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type Device struct {
	Ip                string  `json:"ip"`
	Useragent         string  `json:"useragent"`
	Id                string  `json:"id"`
	Nid               string  `json:"nid"`
	Idfa              *string `json:"idfa,omitempty"`      // [iOS]
	Idfv              *string `json:"idfv,omitempty"`      // [iOS]
	AdId              *string `json:"adId,omitempty"`      // [Android] Google play services advertising id
	AndroidId         *string `json:"androidId,omitempty"` // [Android] Android id
	AdTrackingEnabled *bool   `json:"adTrackingEnabled,omitempty"`
	Manufacturer      *string `json:"manufacturer,omitempty"`
	Model             *string `json:"model,omitempty"`
	Name              *string `json:"name,omitempty"`
	Type              *string `json:"type,omitempty"`
	Token             *string `json:"token,omitempty"`
	Os                `json:"os,omitempty"`
	Browser           `json:"browser,omitempty"`
	Screen            `json:"screen,omitempty"`
	Network           `json:"network,omitempty"`
	App               `json:"app,omitempty"`
	Location          `json:"location,omitempty"`
	Traits            *map[string]interface{} `json:"traits,omitempty"`
}

func (d Device) Value() (driver.Value, error) {
	b, err := json.Marshal(d)
	return string(b), err
}

func (d Device) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &d)
}

type Browser struct {
	Lang           *string `json:"language,omitempty"`
	Cookies        *bool   `json:"cookies,omitempty"`
	ColorDepth     *int64  `json:"colorDepth,omitempty"`
	Charset        *string `json:"charset,omitempty"`
	ViewportSize   *string `json:"viewportSize,omitempty"`
	ViewportWidth  *int    `json:"viewportWidth,omitempty"`
	ViewportHeight *int    `json:"viewportHeight,omitempty"`
	DocumentSize   *string `json:"documentSize,omitempty"`
	DocumentWidth  *int    `json:"documentWidth,omitempty"`
	DocumentHeight *int    `json:"documentHeight,omitempty"`
}

type Screen struct {
	Resolution *string `json:"screenResolution,omitempty"`
	Width      *int    `json:"screenWidth,omitempty"`
	Height     *int    `json:"screenHeight,omitempty"`
}

type Os struct {
	Name     *string `json:"name,omitempty"`
	Version  *string `json:"version,omitempty"`
	Timezone *string `json:"timezone,omitempty"`
}

type Network struct {
	Bluetooth *bool   `json:"bluetooth,omitempty"`
	Cellular  *bool   `json:"cellular,omitempty"`
	Wifi      *bool   `json:"wifi,omitempty"`
	Carrier   *string `json:"carrier,omitempty"`
}

type App struct {
	Name    *string `json:"name,omitempty"`
	Version *string `json:"version,omitempty"`
	Build   *string `json:"build,omitempty"`
}

type Location struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Country   *string  `json:"country,omitempty"`
	Region    *string  `json:"region,omitempty"`
	City      *string  `json:"city,omitempty"`
	Dma       *string  `json:"dma,omitempty"`
}
