package config

type PixelPath struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Type     string `json:"type"` // Either b64 + configured param or params
	B64Param string `json:"b64Param"`
}

type Pixel struct {
	Enabled bool        `json:"enabled"`
	Paths   []PixelPath `json:"paths"`
}
