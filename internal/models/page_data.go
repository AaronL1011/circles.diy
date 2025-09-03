package models

type PageData struct {
	Success   bool
	CSRFToken string
}

type ThemeSettings struct {
	Mode   string `json:"mode"`   // light, dark, system
	Radius string `json:"radius"` // 0, 6, 12, 32
}

type BaseData struct {
	Title     string
	ActiveNav string
	Theme     ThemeSettings
	CSRFToken string
}

type DashboardData struct {
	BaseData
	Feed             []FeedItem        `json:"feed"`
	FeedOffset       int               `json:"feed_offset"`
	Circles          []Circle          `json:"circles"`
	Discussions      []Discussion      `json:"discussions"`
	Events           []Event           `json:"events"`
	Ripples          []Ripple          `json:"ripples"`
	MarketplaceItems []MarketplaceItem `json:"marketplace_items"`
	Impact           []ImpactItem      `json:"impact"`
}

type ProfileData struct {
	BaseData
	Profile      Profile     `json:"profile"`
	Posts        []Post      `json:"posts"`
	PostOffset   int         `json:"post_offset"`
	HasMorePosts bool        `json:"has_more_posts"`
	IsOwner      bool        `json:"is_owner"`
	Extensions   []Extension `json:"extensions"`
	Analytics    Analytics   `json:"analytics"`
	Drafts       []DraftPost `json:"drafts"`
	DraftCount   int         `json:"draft_count"`
}

type CirclesPageData struct {
	BaseData
	Circles         []Circle         `json:"circles"`
	RecentActivity  []CircleActivity `json:"recent_activity"`
	Stats           CircleStats      `json:"stats"`
	FeaturedCircles []Circle         `json:"featured_circles"`
}