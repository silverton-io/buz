package config

type PixelPath struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Pixel struct {
	Enabled bool        `json:"enabled"`
	Paths   []PixelPath `json:"paths"`
}
