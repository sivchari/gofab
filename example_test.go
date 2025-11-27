package gofab

import (
	"fmt"
)

// Example_build demonstrates the basic Build function.
func Example_build() {
	type User struct {
		ID    int    `gofab:"sequence"`
		Name  string `gofab:"name"`
		Email string `gofab:"email"`
	}

	user := Build[User]()
	fmt.Printf("Has ID: %v\n", user.ID > 0)
	fmt.Printf("Has Name: %v\n", user.Name != "")
	fmt.Printf("Has Email: %v\n", user.Email != "")

	// Output:
	// Has ID: true
	// Has Name: true
	// Has Email: true
}

// Example_buildWithCustomizer demonstrates Build with custom values.
func Example_buildWithCustomizer() {
	type Product struct {
		Name  string
		Price float64
	}

	product := Build[Product](func(p *Product) {
		p.Name = "Widget"
		p.Price = 99.99
	})

	fmt.Printf("Name: %s\n", product.Name)
	fmt.Printf("Price: %.2f\n", product.Price)

	// Output:
	// Name: Widget
	// Price: 99.99
}

// Example_buildList demonstrates creating multiple instances.
func Example_buildList() {
	type Task struct {
		Title string
		Done  bool
	}

	const taskTitle = "Task"

	tasks := BuildList[Task](3, func(t *Task) {
		t.Title = taskTitle
		t.Done = false
	})

	fmt.Printf("Created %d tasks\n", len(tasks))
	fmt.Printf("All have title: %v\n", tasks[0].Title == taskTitle && tasks[1].Title == taskTitle)

	// Output:
	// Created 3 tasks
	// All have title: true
}

// Example_define demonstrates factory definition with defaults.
func Example_define() {
	type User struct {
		Name   string
		Active bool
	}

	userFactory := Define[User](func(u *User) {
		u.Active = true // Default value
	})

	user := userFactory.Build(func(u *User) {
		u.Name = "Alice"
	})

	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("Active: %v\n", user.Active)

	// Output:
	// Name: Alice
	// Active: true
}

// Example_trait demonstrates using a single trait.
func Example_trait() {
	type User struct {
		Email  string
		Role   string
		Active bool
	}

	userFactory := Define[User]().
		Trait("admin", func(u *User) {
			u.Role = "admin"
			u.Active = true
		})

	admin := userFactory.Build(
		append(userFactory.WithTrait("admin"),
			func(u *User) {
				u.Email = "admin@example.com"
			})...)

	fmt.Printf("Role: %s\n", admin.Role)
	fmt.Printf("Active: %v\n", admin.Active)

	// Output:
	// Role: admin
	// Active: true
}

// Example_withTraits demonstrates combining multiple traits.
func Example_withTraits() {
	type Product struct {
		Name     string
		Featured bool
		OnSale   bool
	}

	productFactory := Define[Product]().
		Trait("featured", func(p *Product) {
			p.Featured = true
		}).
		Trait("sale", func(p *Product) {
			p.OnSale = true
		})

	product := productFactory.Build(
		append(productFactory.WithTraits("featured", "sale"),
			func(p *Product) {
				p.Name = "Special Item"
			})...)

	fmt.Printf("Name: %s\n", product.Name)
	fmt.Printf("Featured: %v\n", product.Featured)
	fmt.Printf("On Sale: %v\n", product.OnSale)

	// Output:
	// Name: Special Item
	// Featured: true
	// On Sale: true
}

// Example_sequence demonstrates the sequence tag for auto-incrementing IDs.
func Example_sequence() {
	type Order struct {
		ID int `gofab:"sequence"`
	}

	// Build multiple orders to show sequence
	orders := make([]Order, 3)
	for i := range orders {
		orders[i] = Build[Order]()
	}

	// Check that IDs are sequential
	sequential := orders[1].ID > orders[0].ID && orders[2].ID > orders[1].ID
	fmt.Printf("IDs are sequential: %v\n", sequential)

	// Output:
	// IDs are sequential: true
}

// Example_sequenceWithStartValue demonstrates sequence with custom start value.
func Example_sequenceWithStartValue() {
	type Invoice struct {
		Number int `gofab:"sequence:1000"`
	}

	ResetAllSequences()

	invoice := Build[Invoice]()
	fmt.Printf("Invoice number starts at 1000: %v\n", invoice.Number >= 1000)

	// Output:
	// Invoice number starts at 1000: true
}

// Example_resetSequence demonstrates resetting sequences between tests.
func Example_resetSequence() {
	type Item struct {
		ID int `gofab:"sequence"`
	}

	ResetAllSequences()

	item1 := Build[Item]()
	item2 := Build[Item]()

	fmt.Printf("First item ID: %d\n", item1.ID)
	fmt.Printf("Second item ID: %d\n", item2.ID)

	// Reset and build again
	ResetAllSequences()

	item3 := Build[Item]()
	fmt.Printf("After reset, item ID: %d\n", item3.ID)

	// Output:
	// First item ID: 1
	// Second item ID: 2
	// After reset, item ID: 1
}

// Example_newTags demonstrates the new populate tags.
func Example_newTags() {
	type Profile struct {
		UUID     string `gofab:"uuid"`
		Username string `gofab:"username"`
		Website  string `gofab:"url"`
		Active   bool   `gofab:"bool:true"`
		Status   string `gofab:"oneof:active,pending,inactive"`
	}

	profile := Build[Profile]()

	fmt.Printf("Has UUID: %v\n", profile.UUID != "")
	fmt.Printf("Has Username: %v\n", profile.Username != "")
	fmt.Printf("Has Website: %v\n", profile.Website != "")
	fmt.Printf("Is Active: %v\n", profile.Active)
	fmt.Printf("Has Valid Status: %v\n", profile.Status == "active" || profile.Status == "pending" || profile.Status == "inactive")

	// Output:
	// Has UUID: true
	// Has Username: true
	// Has Website: true
	// Is Active: true
	// Has Valid Status: true
}
