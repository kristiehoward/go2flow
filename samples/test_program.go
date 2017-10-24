package schema

import (
	"time"
)

// UpdateProductRequest is used to create or update an existing product, it can be
// used for any type of product, and does not handle children objects
type UpdateProductRequest struct {
	Name             string             `json:"name"`
	PublisherID      string             `json:"publisher_id"`
	CreatedAt        time.Time          `json:"created_at"`
	ShortDescription string             `json:"short_description"`
	FullDescription  string             `json:"full_description,omitempty"`
	Source           string             `json:"source"`
	PopularityScore  int64              `json:"popularity"`
	Categories       []Category         `json:"categories"`
	Links            []URLInfo          `json:"links"`
	Screenshots      []URLInfo          `json:"screenshots"`
	LogoURLs         map[string]string  `json:"logo_url"`
	CrazyMap         map[AliasType]bool `json:"crazy_map"`
	IsOffline        bool               `json:"is_offline,omitempty"`
}

type AliasType string

/*URLInfo is the representation of a link and its label. It can be for external links that we expose
on the product details page, or for screenshots and other images */
type URLInfo struct {
	URL   string `json:"url"`
	Label string `json:"label"`
}

/*
Category contains the name and user-friendly label of a product category
*/
type Category struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

/*
OperatingSystem contains the name and user-friendly label representing an Operating system
*/
type OperatingSystem struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}
