package gofab_test

import (
	"fmt"

	"github.com/sivchari/gofab"
)

type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

func Example() {
	// Define a factory for User
	userFactory := gofab.Define[User]().
		Default(func(u *User) {
			u.ID = 1
			u.Name = "John Doe"
			u.Email = "john@example.com"
			u.Age = 25
		})

	// Create a user using factory defaults
	user1 := userFactory.Build()
	fmt.Printf("User1: %+v\n", user1)

	// Create a user with factory defaults and custom overrides
	user2 := userFactory.Build(func(u *User) {
		u.Name = "Jane Smith"
		u.Age = 30
	})
	fmt.Printf("User2: %+v\n", user2)

	// Create a user without using factory (standalone)
	user3 := gofab.Create[User](func(u *User) {
		u.ID = 3
		u.Name = "Bob Wilson"
		u.Email = "bob@example.com"
		u.Age = 35
	})
	fmt.Printf("User3: %+v\n", user3)

	// Output:
	// User1: {ID:1 Name:John Doe Email:john@example.com Age:25}
	// User2: {ID:1 Name:Jane Smith Email:john@example.com Age:30}
	// User3: {ID:3 Name:Bob Wilson Email:bob@example.com Age:35}
}
