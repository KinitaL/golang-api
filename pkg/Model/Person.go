package model

type Person struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	IsAdmin  string `json:"isAdmin"`
}
