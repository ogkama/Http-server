package domain

type User struct {
	Id       string   `json:"id"`
	Login    string   `json:"login"`
	Password [32]byte `json:"password"`
}