package models

type User struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`
	Banner string `json:"banner"`
}

type ProfileStats struct {
	Posts       int `json:"posts"`
	Connections int `json:"connections"`
	Circles     int `json:"circles"`
}

type Profile struct {
	ID          string       `json:"id"`
	Handle      string       `json:"handle"`
	Name        string       `json:"name"`
	Avatar      string       `json:"avatar"`
	Banner      string       `json:"banner"`
	Bio         string       `json:"bio"`
	Stats       ProfileStats `json:"stats"`
	IsConnected bool         `json:"is_connected"`
}