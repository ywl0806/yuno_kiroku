package store

type UserStore interface {
	FindUsers()
	FindUserById()
	CreateUser()
}
