package models

type MarketplaceItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Price    string `json:"price"`
	Image    string `json:"image"`
	Location string `json:"location"`
	TimeAgo  string `json:"time_ago"`
}