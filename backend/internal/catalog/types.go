package catalog

const (
	CatImage = "image"
	CatVideo = "video"
)

const (
	TypeString    = "string"
	TypeInt       = "int"
	TypeFloat     = "float"
	TypeBool      = "bool"
	TypeSelect    = "select"
	TypeMedia     = "media"
	TypeMediaList = "media_list"
	TypeJSON      = "json"
)

const (
	ScopeInput      = "input"
	ScopeParameters = "parameters"
)

type ModelEntry struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Category    string     `json:"category"`
	Endpoint    string     `json:"endpoint"`
	Async       bool       `json:"async"`
	Payload     string     `json:"payload"`
	Description string     `json:"description"`
	Parameters  []ParamDef `json:"parameters"`
}

type ParamDef struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Type        string   `json:"type"`
	Scope       string   `json:"scope"`
	Required    bool     `json:"required"`
	Default     any      `json:"default,omitempty"`
	Min         *float64 `json:"min,omitempty"`
	Max         *float64 `json:"max,omitempty"`
	Step        *float64 `json:"step,omitempty"`
	Options     []Option `json:"options,omitempty"`
	Description string   `json:"description,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
}

type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Catalog struct {
	Models []ModelEntry `json:"models"`
}

func (c *Catalog) ByCategory(cat string) []ModelEntry {
	var result []ModelEntry
	for _, m := range c.Models {
		if m.Category == cat {
			result = append(result, m)
		}
	}
	return result
}

func (c *Catalog) Find(id string) *ModelEntry {
	for i := range c.Models {
		if c.Models[i].ID == id {
			return &c.Models[i]
		}
	}
	return nil
}
