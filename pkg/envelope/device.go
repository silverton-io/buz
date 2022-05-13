package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type Device struct {
	Ip                *string `json:"ip"`
	Useragent         *string `json:"useragent"`
	Duid              *string `json:"duid"`
	Nuid              *string `json:"nuid"`
	Timezone          *string `json:"timezone"`
	Id                *string `json:"id,omitempty"`
	AdId              *string `json:"adId,omitempty"`
	AdTrackingEnabled *bool   `json:"adTrackingEnabled,omitempty"`
	Manufacturer      *string `json:"manufacturer,omitempty"`
	Model             *string `json:"model,omitempty"`
	Name              *string `json:"name,omitempty"`
	Type              *string `json:"type,omitempty"`
	Token             *string `json:"token,omitempty"`
	Os                `json:"os"`
	Browser           `json:"browser"`
	Screen            `json:"screen"`
}

func (d Device) Value() (driver.Value, error) {
	b, err := json.Marshal(d)
	return string(b), err
}

func (d Device) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &d)
}

type Browser struct {
	Cookies    *bool   `json:"cookies"`
	Lang       *string `json:"lang"`
	ColorDepth *int64  `json:"colorDepth"`
}

type Screen struct {
	ViewportSize     *string `json:"viewportSize"`
	ViewportWidth    *int    `json:"viewportWidth"`
	ViewportHeight   *int    `json:"viewportHeight"`
	Charset          *string `json:"charset"`
	DocumentSize     *string `json:"documentSize"`
	DocumentWidth    *int    `json:"documentWidth"`
	DocumentHeight   *int    `json:"documentHeight"`
	ScreenResolution *string `json:"screenResolution"`
	ScreenWidth      *int    `json:"screenWidth"`
	ScreenHeight     *int    `json:"screenHeight"`
}

type Os struct {
	Name    *string `json:"name"`
	Version *string `json:"version"`
}

type Network struct {
	Bluetooth *bool   `json:"bluetooth"`
	Cellular  *bool   `json:"cellular"`
	Wifi      *bool   `json:"wifi"`
	Carrier   *string `json:"carrier"`
}
