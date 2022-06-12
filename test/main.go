package main

import (
	"fmt"

	"github.com/emehrkay/khadijah"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	mark := User{
		ID:    "someID",
		Name:  "emehrkay",
		Email: "spam@aol.com",
	}
	mag := khadijah.New()
	create := mag.CreateNode(mark, "User", true)

	fmt.Printf(`%+v`, create)
}
