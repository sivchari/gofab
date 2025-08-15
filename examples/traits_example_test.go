package examples_test

import (
	"fmt"
	"strings"

	"github.com/sivchari/gofab"
)

// ExampleFactory_Trait demonstrates defining and using traits
// for reusable attribute sets.
func ExampleFactory_Trait() {
	type User struct {
		Name        string
		Role        string
		Permissions []string
		Active      bool
	}

	// Define a factory with named traits
	userFactory := gofab.Define[User]().
		Trait("admin", func(u *User) {
			u.Role = "admin"
			u.Permissions = []string{"read", "write", "delete"}
			u.Active = true
		}).
		Trait("moderator", func(u *User) {
			u.Role = "moderator"
			u.Permissions = []string{"read", "write"}
			u.Active = true
		}).
		Trait("inactive", func(u *User) {
			u.Active = false
		})

	// Build an admin user
	admin := userFactory.Build(
		append(userFactory.WithTrait("admin"),
			func(u *User) {
				u.Name = "Alice"
			})...)

	fmt.Println("Admin:", admin.Name)
	fmt.Println("Role:", admin.Role)
	fmt.Println("Permissions:", strings.Join(admin.Permissions, ", "))
	fmt.Println("Active:", admin.Active)

	// Output:
	// Admin: Alice
	// Role: admin
	// Permissions: read, write, delete
	// Active: true
}

// ExampleFactory_WithTraits demonstrates combining multiple traits.
func ExampleFactory_WithTraits() {
	type User struct {
		Name      string
		Role      string
		Active    bool
		Verified  bool
		Premium   bool
	}

	userFactory := gofab.Define[User]().
		Trait("admin", func(u *User) {
			u.Role = "admin"
			u.Active = true
		}).
		Trait("verified", func(u *User) {
			u.Verified = true
		}).
		Trait("premium", func(u *User) {
			u.Premium = true
		}).
		Trait("suspended", func(u *User) {
			u.Active = false
		})

	// Combine multiple traits
	user := userFactory.Build(
		append(userFactory.WithTraits("admin", "verified", "premium"),
			func(u *User) {
				u.Name = "Bob"
			})...)

	fmt.Println("Name:", user.Name)
	fmt.Println("Role:", user.Role)
	fmt.Println("Active:", user.Active)
	fmt.Println("Verified:", user.Verified)
	fmt.Println("Premium:", user.Premium)

	// Output:
	// Name: Bob
	// Role: admin
	// Active: true
	// Verified: true
	// Premium: true
}

// ExampleFactory_Trait_override demonstrates how traits can override each other
// when applied in sequence.
func ExampleFactory_Trait_override() {
	type Account struct {
		Type   string
		Active bool
		Limit  int
	}

	accountFactory := gofab.Define[Account]().
		Trait("basic", func(a *Account) {
			a.Type = "basic"
			a.Active = true
			a.Limit = 100
		}).
		Trait("premium", func(a *Account) {
			a.Type = "premium"
			a.Active = true
			a.Limit = 1000
		}).
		Trait("suspended", func(a *Account) {
			a.Active = false
			a.Limit = 0
		})

	// Apply premium then suspended - suspended overrides active status
	var builders []gofab.Builder[Account]
	builders = append(builders, accountFactory.WithTrait("premium")...)
	builders = append(builders, accountFactory.WithTrait("suspended")...)

	account := accountFactory.Build(builders...)

	fmt.Println("Type:", account.Type)       // Keeps premium type
	fmt.Println("Active:", account.Active)   // Overridden by suspended
	fmt.Println("Limit:", account.Limit)     // Overridden by suspended

	// Output:
	// Type: premium
	// Active: false
	// Limit: 0
}

// ExampleFactory_Trait_withList demonstrates using traits with BuildList.
func ExampleFactory_Trait_withList() {
	type Employee struct {
		ID         int
		Department string
		Level      string
	}

	employeeFactory := gofab.Define[Employee]().
		Trait("engineering", func(e *Employee) {
			e.Department = "Engineering"
			e.Level = "Senior"
		}).
		Trait("junior", func(e *Employee) {
			e.Level = "Junior"
		})

	// Create multiple engineers
	engineers := employeeFactory.BuildList(3,
		employeeFactory.WithTrait("engineering")...)

	for i, eng := range engineers {
		eng.ID = i + 1 // Manually set IDs for example
		fmt.Printf("Engineer %d: %s %s\n", eng.ID, eng.Level, eng.Department)
	}

	// Output:
	// Engineer 1: Senior Engineering
	// Engineer 2: Senior Engineering
	// Engineer 3: Senior Engineering
}