package examples_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/sivchari/gofab"
)

// ExampleFactory_combined demonstrates combining Traits and AfterBuild
// for complex test data setup.
func ExampleFactory_combined() {
	type Article struct {
		ID          int
		Title       string
		Body        string
		Author      string
		Status      string
		Tags        []string
		Slug        string    // Generated in AfterBuild
		WordCount   int       // Calculated in AfterBuild
		PublishedAt time.Time // Set by trait or AfterBuild
	}

	// Create a factory with traits and AfterBuild callbacks
	articleFactory := gofab.Define[Article]().
		Trait("published", func(a *Article) {
			a.Status = "published"
			a.PublishedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		}).
		Trait("draft", func(a *Article) {
			a.Status = "draft"
			a.PublishedAt = time.Time{} // Zero value
		}).
		Trait("featured", func(a *Article) {
			a.Tags = append(a.Tags, "featured")
		}).
		AfterBuild(func(a *Article) {
			// Generate slug from title
			a.Slug = strings.ToLower(strings.ReplaceAll(a.Title, " ", "-"))
		}).
		AfterBuild(func(a *Article) {
			// Calculate word count
			a.WordCount = len(strings.Fields(a.Body))
		}).
		AfterBuild(func(a *Article) {
			// Add default tag based on word count
			if a.WordCount < 100 {
				a.Tags = append(a.Tags, "short-read")
			} else {
				a.Tags = append(a.Tags, "long-read")
			}
		})

	// Create a featured published article
	article := articleFactory.Build(
		append(articleFactory.WithTraits("published", "featured"),
			func(a *Article) {
				a.ID = 1
				a.Title = "Getting Started with Go"
				a.Body = "Go is a statically typed, compiled language..."
				a.Author = "Jane Doe"
			})...)

	fmt.Println("Title:", article.Title)
	fmt.Println("Slug:", article.Slug)
	fmt.Println("Status:", article.Status)
	fmt.Println("Tags:", strings.Join(article.Tags, ", "))
	fmt.Println("Word count:", article.WordCount)

	// Output:
	// Title: Getting Started with Go
	// Slug: getting-started-with-go
	// Status: published
	// Tags: featured, short-read
	// Word count: 7
}

// ExampleFactory_testScenarios demonstrates using factory patterns
// for different test scenarios.
func ExampleFactory_testScenarios() {
	userFactory := createUserTestFactory()

	// Test scenario 1: New user registration
	newUser := userFactory.Build(
		append(userFactory.WithTrait("newUser"),
			func(u *User) {
				u.ID = 1
				u.Email = "newbie@example.com"
			})...)

	fmt.Println("Scenario 1 - New User:")
	fmt.Printf("  Username: %s, Verified: %v, Login Count: %d\n",
		newUser.Username, newUser.EmailVerified, newUser.LoginCount)

	// Test scenario 2: Admin user
	admin := userFactory.Build(
		append(userFactory.WithTrait("admin"),
			func(u *User) {
				u.ID = 2
				u.Email = "admin@example.com"
			})...)

	fmt.Println("Scenario 2 - Admin:")
	fmt.Printf("  Role: %s, Active: %v, Verified: %v\n",
		admin.Role, admin.IsActive, admin.EmailVerified)

	// Test scenario 3: Banned active user
	bannedUser := userFactory.Build(
		append(userFactory.WithTraits("activeUser", "banned"),
			func(u *User) {
				u.ID = 3
				u.Email = "banned@example.com"
			})...)

	fmt.Println("Scenario 3 - Banned User:")
	fmt.Printf("  Active: %v, Verified: %v (reset by AfterBuild)\n",
		bannedUser.IsActive, bannedUser.EmailVerified)

	// Output:
	// Scenario 1 - New User:
	//   Username: newbie, Verified: false, Login Count: 0
	// Scenario 2 - Admin:
	//   Role: admin, Active: true, Verified: true
	// Scenario 3 - Banned User:
	//   Active: false, Verified: false (reset by AfterBuild)
}

type User struct {
	ID            int
	Email         string
	Username      string
	Role          string
	IsActive      bool
	EmailVerified bool
	LoginCount    int
	LastLogin     time.Time
}

func createUserTestFactory() *gofab.Factory[User] {
	// Create a comprehensive user factory for testing
	return gofab.Define[User]().
		// Define common user types as traits
		Trait("newUser", func(u *User) {
			u.IsActive = true
			u.EmailVerified = false
			u.LoginCount = 0
			u.Role = "user"
		}).
		Trait("activeUser", func(u *User) {
			u.IsActive = true
			u.EmailVerified = true
			u.LoginCount = 10
			u.Role = "user"
			u.LastLogin = time.Now().Add(-1 * time.Hour)
		}).
		Trait("admin", func(u *User) {
			u.Role = "admin"
			u.IsActive = true
			u.EmailVerified = true
		}).
		Trait("banned", func(u *User) {
			u.IsActive = false
		}).
		// AfterBuild to ensure data consistency
		AfterBuild(func(u *User) {
			// Generate username from email if not set
			if u.Username == "" && u.Email != "" {
				parts := strings.Split(u.Email, "@")
				u.Username = parts[0]
			}
		}).
		AfterBuild(func(u *User) {
			// Banned users should have verified email reset
			if !u.IsActive && u.Role != "admin" {
				u.EmailVerified = false
			}
		})
}
