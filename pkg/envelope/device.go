package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type Device struct {
	Ip                string  `json:"ip"`
	Useragent         string  `json:"useragent"`
	Id                *string `json:"id"`
	Nid               *string `json:"nid"`
	Idfa              *string `json:"idfa"`      // [iOS]
	Idfv              *string `json:"idfv"`      // [iOS]
	AdId              *string `json:"adId"`      // [Android] Google play services advertising id
	AndroidId         *string `json:"androidId"` // [Android] Android id
	AdTrackingEnabled *bool   `json:"adTrackingEnabled"`
	Manufacturer      *string `json:"manufacturer"`
	Model             *string `json:"model"`
	Name              *string `json:"name"`
	Type              *string `json:"type"`
	Token             *string `json:"token"`
	Os                `json:"os"`
	Browser           `json:"browser"`
	Screen            `json:"screen"`
	Network           `json:"network"`
	App               `json:"app"`
	Location          `json:"location"`
	Traits            *map[string]interface{} `json:"traits"`
}

func (d Device) Value() (driver.Value, error) {
	b, err := json.Marshal(d)
	return string(b), err
}

func (d Device) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &d)
}

type Browser struct {
	Lang           *string `json:"language"`
	Cookies        *bool   `json:"cookies"`
	ColorDepth     *int64  `json:"colorDepth"`
	Charset        *string `json:"charset"`
	ViewportSize   *string `json:"viewportSize"`
	ViewportWidth  *int    `json:"viewportWidth"`
	ViewportHeight *int    `json:"viewportHeight"`
	DocumentSize   *string `json:"documentSize"`
	DocumentWidth  *int    `json:"documentWidth"`
	DocumentHeight *int    `json:"documentHeight"`
}

type Screen struct {
	Resolution *string `json:"screenResolution"`
	Width      *int    `json:"screenWidth"`
	Height     *int    `json:"screenHeight"`
}

type Os struct {
	Name     *string `json:"name"`
	Version  *string `json:"version"`
	Timezone *string `json:"timezone"`
}

type Network struct {
	Bluetooth *bool   `json:"bluetooth"`
	Cellular  *bool   `json:"cellular"`
	Wifi      *bool   `json:"wifi"`
	Carrier   *string `json:"carrier"`
}

type App struct {
	Name    *string `json:"name"`
	Version *string `json:"version"`
	Build   *string `json:"build"`
}

type Location struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Country   *string  `json:"country"`
	Region    *string  `json:"region"`
	City      *string  `json:"city"`
	Dma       *string  `json:"dma"`
}
