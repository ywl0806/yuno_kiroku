package models

type User struct {
	ID       int    `firestore:"id"`
	Username string `firestore:"username"`
	Name     string `firestore:"name"`
}

type UserCred struct {
	Password string `firestore:"password"`
}
