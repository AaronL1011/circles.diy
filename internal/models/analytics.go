package models

type Analytics struct {
	ProfileViews         int `json:"profile_views"`
	ProfileViewsChange   int `json:"profile_views_change"`
	PostEngagement       int `json:"post_engagement"`
	PostEngagementChange int `json:"post_engagement_change"`
	NewConnections       int `json:"new_connections"`
	NewConnectionsChange int `json:"new_connections_change"`
}

type ImpactItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}