package models

type MediaItem struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
}

type PostStats struct {
	Replies int `json:"replies"`
	Shares  int `json:"shares"`
	Views   int `json:"views"`
}

type FeedItem struct {
	ID      string      `json:"id"`
	User    User        `json:"user"`
	Content string      `json:"content"`
	TimeAgo string      `json:"time_ago"`
	Circle  string      `json:"circle"`
	Image   *MediaItem  `json:"image,omitempty"`
	Video   *MediaItem  `json:"video,omitempty"`
	Gallery []MediaItem `json:"gallery,omitempty"`
	CanBuy  bool        `json:"can_buy"`
}

type Post struct {
	ID      string      `json:"id"`
	User    User        `json:"user"`
	Content string      `json:"content"`
	TimeAgo string      `json:"time_ago"`
	Circle  string      `json:"circle"`
	Image   *MediaItem  `json:"image,omitempty"`
	Video   *MediaItem  `json:"video,omitempty"`
	Gallery []MediaItem `json:"gallery,omitempty"`
	CanBuy  bool        `json:"can_buy"`
	Stats   *PostStats  `json:"stats,omitempty"`
}

type DraftPost struct {
	ID      string      `json:"id"`
	Content string      `json:"content"`
	Image   *MediaItem  `json:"image,omitempty"`
	Video   *MediaItem  `json:"video,omitempty"`
	Gallery []MediaItem `json:"gallery,omitempty"`
}