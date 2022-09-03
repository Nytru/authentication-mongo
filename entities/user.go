package entities

type User struct {
	Name string `json:"name"`
	Id int `json:"id"` // `mongo:"id"`
}