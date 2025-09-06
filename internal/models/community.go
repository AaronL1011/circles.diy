package models

type Circle struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Thumbnail    string `json:"thumbnail"`
	Banner       string `json:"banner"`
	MemberCount  string `json:"member_count"`
	OnlineCount  string `json:"online_count"`
	UserRole     string `json:"user_role"`     // owner, admin, member
	JoinedDate   string `json:"joined_date"`
	LastActivity string `json:"last_activity"`
	Active       bool   `json:"active"`
}

type Discussion struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Preview string `json:"preview"`
	Circle  string `json:"circle"`
	TimeAgo string `json:"time_ago"`
}

type Event struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Time  string `json:"time"`
	Day   string `json:"day"`
	Month string `json:"month"`
}

type Ripple struct {
	ID          string      `json:"id"`
	User        User        `json:"user"`
	Content     string      `json:"content"`
	ContentType string      `json:"content_type"`
	Image       *MediaItem  `json:"image,omitempty"`
	Video       *MediaItem  `json:"video,omitempty"`
	Gallery     []MediaItem `json:"gallery,omitempty"`
	Link        *LinkPreview `json:"link,omitempty"`
	ExpiresIn   string      `json:"expires_in"`
	Circle      string      `json:"circle"`
	ViewCount   int         `json:"view_count"`
}

type LinkPreview struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image,omitempty"`
	Domain      string `json:"domain"`
}

type CircleActivity struct {
	ID       string `json:"id"`
	CircleID string `json:"circle_id"`
	Type     string `json:"type"` // post, member_joined, event, announcement
	Title    string `json:"title"`
	Content  string `json:"content"`
	User     string `json:"user"`
	TimeAgo  string `json:"time_ago"`
}

type CircleStats struct {
	TotalPosts      int    `json:"total_posts"`
	ActiveMembers   int    `json:"active_members"`
	RecentActivity  string `json:"recent_activity"`
	WeeklyGrowth    string `json:"weekly_growth"`
	EngagementRate  string `json:"engagement_rate"`
}