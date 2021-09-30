package models

type User struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Cover   string `json:"cover"`
	Address string `json:"address"`
}
