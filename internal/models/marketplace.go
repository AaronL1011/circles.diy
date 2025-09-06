package models

type MarketplaceItem struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Price       string     `json:"price"`
	PriceType   string     `json:"price_type"` // sale, trade, free, negotiable
	Image       *MediaItem `json:"image,omitempty"`
	Images      []MediaItem `json:"images,omitempty"`
	Location    string     `json:"location"`
	Distance    string     `json:"distance,omitempty"`
	TimeAgo     string     `json:"time_ago"`
	Seller      User       `json:"seller"`
	Circle      string     `json:"circle,omitempty"`
	Category    string     `json:"category"`
	Tags        []string   `json:"tags"`
	Condition   string     `json:"condition"` // new, like-new, good, fair, poor
	IsAvailable bool       `json:"is_available"`
	ViewCount   int        `json:"view_count"`
	IsFeatured  bool       `json:"is_featured"`
}

type MarketplaceCategory struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Count int    `json:"count"`
}

type MarketplaceFilter struct {
	PriceTypes []string   `json:"price_types"`
	Categories []string   `json:"categories"`
	Conditions []string   `json:"conditions"`
	Locations  []Location `json:"locations"`
	MaxDistance int       `json:"max_distance"`
}

type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Count     int     `json:"count"`
}