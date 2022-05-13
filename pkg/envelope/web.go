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
	Query    *map[string]interface{} `json:"query,omitempty"`
	Fragment *string                 `json:"fragment,omitempty"`
	Medium   *string                 `json:"medium,omitempty"`
	Source   *string                 `json:"source,omitempty"`
	Term     *string                 `json:"term,omitempty"`
	Content  *string                 `json:"content,omitempty"`
	Campaign *string                 `json:"campaign,omitempty"`
}
