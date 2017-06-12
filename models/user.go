package models

import (
)

func init() {
}

type User struct {
	Id       string
	Username string
	Password string
}

func AddUser(u User) string {
	return nil
}
