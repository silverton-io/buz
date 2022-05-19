package envelope

type Page struct {
	Page     PageAttrs `json:"page"`
	Referrer PageAttrs `json:"referrer"`
}

type PageAttrs struct {
	Url      string                  `json:"url"`
	Title    *string                 `json:"title"`
	Scheme   string                  `json:"scheme"`
	Host     string                  `json:"host"`
	Port     string                  `json:"port"`
	Path     string                  `json:"path"`
	Query    *map[string]interface{} `json:"query"`
	Fragment *string                 `json:"fragment"`
	Medium   *string                 `json:"medium"`
	Source   *string                 `json:"source"`
	Term     *string                 `json:"term"`
	Content  *string                 `json:"content"`
	Campaign *string                 `json:"campaign"`
}
