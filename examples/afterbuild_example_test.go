package examples_test

import (
	"fmt"
	"strings"

	"github.com/sivchari/gofab"
)

// ExampleFactory_AfterBuild demonstrates using AfterBuild callbacks
// to automatically process instances after building.
func ExampleFactory_AfterBuild() {
	type BlogPost struct {
		Title     string
		Body      string
		Slug      string // Generated in AfterBuild
		WordCount int    // Calculated in AfterBuild
	}

	// Define a factory with AfterBuild callbacks
	postFactory := gofab.Define[BlogPost]().
		AfterBuild(func(p *BlogPost) {
			// Generate slug from title
			p.Slug = strings.ToLower(strings.ReplaceAll(p.Title, " ", "-"))
		}).
		AfterBuild(func(p *BlogPost) {
			// Calculate word count
			p.WordCount = len(strings.Fields(p.Body))
		})

	// Build a post with custom title and body
	post := postFactory.Build(func(p *BlogPost) {
		p.Title = "Hello World"
		p.Body = "This is my first blog post"
	})

	fmt.Println("Title:", post.Title)
	fmt.Println("Slug:", post.Slug)
	fmt.Println("Word count:", post.WordCount)

	// Output:
	// Title: Hello World
	// Slug: hello-world
	// Word count: 6
}

// ExampleFactory_AfterBuild_chain demonstrates chaining multiple AfterBuild callbacks
// where each callback can use the results of previous callbacks.
func ExampleFactory_AfterBuild_chain() {
	type User struct {
		Name     string
		Username string // Generated from Name
		Email    string // Generated from Username
	}

	userFactory := gofab.Define[User]().
		AfterBuild(func(u *User) {
			// First callback: generate username from name
			u.Username = strings.ToLower(strings.ReplaceAll(u.Name, " ", "."))
		}).
		AfterBuild(func(u *User) {
			// Second callback: generate email from username
			u.Email = fmt.Sprintf("%s@example.com", u.Username)
		})

	user := userFactory.Build(func(u *User) {
		u.Name = "John Doe"
	})

	fmt.Println("Name:", user.Name)
	fmt.Println("Username:", user.Username)
	fmt.Println("Email:", user.Email)

	// Output:
	// Name: John Doe
	// Username: john.doe
	// Email: john.doe@example.com
}

// ExampleFactory_AfterBuild_withList demonstrates that AfterBuild callbacks
// are applied to each instance when using BuildList.
func ExampleFactory_AfterBuild_withList() {
	type Product struct {
		Name  string
		Price float64
		Tax   float64 // Calculated in AfterBuild
		Total float64 // Calculated in AfterBuild
	}

	productFactory := gofab.Define[Product]().
		AfterBuild(func(p *Product) {
			p.Tax = p.Price * 0.1      // 10% tax
			p.Total = p.Price + p.Tax
		})

	products := productFactory.BuildList(2, func(p *Product) {
		p.Name = "Widget"
		p.Price = 100.00
	})

	for i, product := range products {
		fmt.Printf("Product %d - Price: %.2f, Tax: %.2f, Total: %.2f\n",
			i+1, product.Price, product.Tax, product.Total)
	}

	// Output:
	// Product 1 - Price: 100.00, Tax: 10.00, Total: 110.00
	// Product 2 - Price: 100.00, Tax: 10.00, Total: 110.00
}