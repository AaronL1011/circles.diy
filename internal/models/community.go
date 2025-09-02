package models

type Circle struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Thumbnail   string `json:"thumbnail"`
	MemberCount string `json:"member_count"`
	Active      bool   `json:"active"`
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
	ID        string `json:"id"`
	User      string `json:"user"`
	Content   string `json:"content"`
	ExpiresIn string `json:"expires_in"`
}