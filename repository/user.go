package repository

type User interface {
	GetID(string) (*string, error)
	GetPassword(string) (*string, error)
	AddUser(string, string, string) (string, error)
	DeleteUser(string) error
}