package gofab_test

import (
	"strings"
	"testing"

	"github.com/sivchari/gofab"
)

const (
	roleUser  = "user"
	roleAdmin = "admin"
)

type IntegratedTestUser struct {
	ID     int    `gofab:"sequence"`
	Name   string `gofab:"name"`
	Email  string `gofab:"email"`
	Role   string // No tag - will be set by factory defaults
	Active bool   `gofab:"-"` // Skip auto-generation
	NoTag  string // No tag, no factory default
}

func TestDefineWithAutoGeneration(t *testing.T) {
	// Create factory with defaults that will be applied after auto-generation
	userFactory := gofab.Define[IntegratedTestUser](func(u *IntegratedTestUser) {
		u.Role = roleUser // This will override any auto-generated value
		u.Active = true   // This will set the skipped field
	})

	user := userFactory.Build()

	// Check auto-generated fields
	if user.ID == 0 {
		t.Error("ID should be auto-generated")
	}

	if user.Name == "" {
		t.Error("Name should be auto-generated")
	}

	if user.Email == "" {
		t.Error("Email should be auto-generated")
	}

	// Check factory defaults
	if user.Role != roleUser {
		t.Errorf("Role should be set by factory default, got %s", user.Role)
	}

	if !user.Active {
		t.Error("Active should be set by factory default")
	}

	// Check untouched field
	if user.NoTag != "" {
		t.Error("NoTag should remain empty")
	}
}

func TestDefineWithOverrides(t *testing.T) {
	adminFactory := gofab.Define[IntegratedTestUser](func(u *IntegratedTestUser) {
		u.Role = roleAdmin
		u.Active = true
	})

	// Override some fields including auto-generated ones
	superAdmin := adminFactory.Build(func(u *IntegratedTestUser) {
		u.Name = "Super Admin" // Override auto-generated name
		u.NoTag = "Special"    // Set the untouched field
	})

	// Check that overrides work
	if superAdmin.Name != "Super Admin" {
		t.Errorf("Name should be overridden, got %s", superAdmin.Name)
	}

	if superAdmin.NoTag != "Special" {
		t.Errorf("NoTag should be set by override, got %s", superAdmin.NoTag)
	}

	// Check that factory defaults are still applied
	if superAdmin.Role != roleAdmin {
		t.Errorf("Role should be set by factory, got %s", superAdmin.Role)
	}

	// Check that auto-generation still works for non-overridden fields
	if superAdmin.Email == "" {
		t.Error("Email should still be auto-generated")
	}
}

func TestFactoryBuildList(t *testing.T) {
	adminFactory := gofab.Define[IntegratedTestUser](func(u *IntegratedTestUser) {
		u.Role = roleAdmin
		u.Active = true
	})

	admins := adminFactory.BuildList(3)

	if len(admins) != 3 {
		t.Errorf("Expected 3 admins, got %d", len(admins))
	}

	// Check that all admins have unique IDs (sequence)
	ids := make(map[int]bool)
	for i, admin := range admins {
		if ids[admin.ID] {
			t.Errorf("Duplicate ID found: %d", admin.ID)
		}

		ids[admin.ID] = true

		// Check that each admin has the factory defaults
		if admin.Role != roleAdmin {
			t.Errorf("Admin %d should have role 'admin', got %s", i, admin.Role)
		}

		if !admin.Active {
			t.Errorf("Admin %d should be active", i)
		}

		// Check that auto-generation works
		if admin.Name == "" || admin.Email == "" {
			t.Errorf("Admin %d should have auto-generated name and email", i)
		}
	}
}

func TestFactoryBuildListWithCustomizer(t *testing.T) {
	userFactory := gofab.Define[IntegratedTestUser](func(u *IntegratedTestUser) {
		u.Role = roleUser
	})

	users := userFactory.BuildList(2, func(u *IntegratedTestUser) {
		u.Active = true
		u.NoTag = "Batch User"
	})

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	for i, user := range users {
		// Check factory defaults
		if user.Role != roleUser {
			t.Errorf("User %d should have role 'user', got %s", i, user.Role)
		}

		// Check customizer
		if !user.Active {
			t.Errorf("User %d should be active from customizer", i)
		}

		if user.NoTag != "Batch User" {
			t.Errorf("User %d should have NoTag 'Batch User', got %s", i, user.NoTag)
		}

		// Check auto-generation
		if user.Name == "" || user.Email == "" {
			t.Errorf("User %d should have auto-generated name and email", i)
		}
	}
}

type BlogPost struct {
	ID        int    `gofab:"sequence"`
	Title     string `gofab:"sentence:3"`
	Body      string `gofab:"sentence:10"`
	Slug      string // Generated in AfterBuild
	WordCount int    // Calculated in AfterBuild
}

func TestAfterBuild(t *testing.T) {
	postFactory := gofab.Define[BlogPost]().
		AfterBuild(func(p *BlogPost) {
			// Generate slug from title
			p.Slug = strings.ToLower(strings.ReplaceAll(p.Title, " ", "-"))
		}).
		AfterBuild(func(p *BlogPost) {
			// Calculate word count
			p.WordCount = len(strings.Fields(p.Body))
		})

	post := postFactory.Build()

	// Check that AfterBuild callbacks were executed
	if post.Slug == "" {
		t.Error("Slug should be generated in AfterBuild")
	}

	if post.WordCount == 0 {
		t.Error("WordCount should be calculated in AfterBuild")
	}

	// Check that slug is derived from title
	expectedSlug := strings.ToLower(strings.ReplaceAll(post.Title, " ", "-"))
	if post.Slug != expectedSlug {
		t.Errorf("Slug should be derived from title, expected %s, got %s", expectedSlug, post.Slug)
	}

	// Check that word count is correct
	expectedWordCount := len(strings.Fields(post.Body))
	if post.WordCount != expectedWordCount {
		t.Errorf("WordCount should be %d, got %d", expectedWordCount, post.WordCount)
	}
}

func TestAfterBuildWithOverrides(t *testing.T) {
	postFactory := gofab.Define[BlogPost]().
		AfterBuild(func(p *BlogPost) {
			p.Slug = strings.ToLower(strings.ReplaceAll(p.Title, " ", "-"))
		}).
		AfterBuild(func(p *BlogPost) {
			p.WordCount = len(strings.Fields(p.Body))
		})

	// Build with custom title and body
	post := postFactory.Build(func(p *BlogPost) {
		p.Title = "Custom Title Here"
		p.Body = "This is a custom body"
	})

	// Check that AfterBuild works with custom values
	if post.Slug != "custom-title-here" {
		t.Errorf("Slug should be 'custom-title-here', got %s", post.Slug)
	}

	if post.WordCount != 5 {
		t.Errorf("WordCount should be 5, got %d", post.WordCount)
	}
}

func TestAfterBuildInBuildList(t *testing.T) {
	postFactory := gofab.Define[BlogPost]().
		AfterBuild(func(p *BlogPost) {
			p.Slug = strings.ToLower(strings.ReplaceAll(p.Title, " ", "-"))
		}).
		AfterBuild(func(p *BlogPost) {
			p.WordCount = len(strings.Fields(p.Body))
		})

	posts := postFactory.BuildList(3)

	if len(posts) != 3 {
		t.Errorf("Expected 3 posts, got %d", len(posts))
	}

	for i, post := range posts {
		// Check that AfterBuild was called for each post
		if post.Slug == "" {
			t.Errorf("Post %d should have a slug", i)
		}

		if post.WordCount == 0 {
			t.Errorf("Post %d should have a word count", i)
		}

		// Verify slug matches title
		expectedSlug := strings.ToLower(strings.ReplaceAll(post.Title, " ", "-"))
		if post.Slug != expectedSlug {
			t.Errorf("Post %d slug mismatch, expected %s, got %s", i, expectedSlug, post.Slug)
		}
	}
}

func TestMultipleAfterBuildCallbacks(t *testing.T) {
	callOrder := []string{}

	factory := gofab.Define[IntegratedTestUser]().
		AfterBuild(func(u *IntegratedTestUser) {
			callOrder = append(callOrder, "first")
			u.NoTag = "first-"
		}).
		AfterBuild(func(u *IntegratedTestUser) {
			callOrder = append(callOrder, "second")
			u.NoTag += "second-"
		}).
		AfterBuild(func(u *IntegratedTestUser) {
			callOrder = append(callOrder, "third")
			u.NoTag += "third"
		})

	user := factory.Build()

	// Check callbacks were called in order
	if len(callOrder) != 3 {
		t.Errorf("Expected 3 callbacks, got %d", len(callOrder))
	}

	if callOrder[0] != "first" || callOrder[1] != "second" || callOrder[2] != "third" {
		t.Errorf("Callbacks should be called in order, got %v", callOrder)
	}

	// Check that all callbacks affected the result
	if user.NoTag != "first-second-third" {
		t.Errorf("NoTag should be 'first-second-third', got %s", user.NoTag)
	}
}

func TestTraits(t *testing.T) {
	factory := gofab.Define[IntegratedTestUser]().
		Trait("admin", func(u *IntegratedTestUser) {
			u.Role = roleAdmin
			u.Active = true
		}).
		Trait("inactive", func(u *IntegratedTestUser) {
			u.Active = false
			u.NoTag = "deactivated"
		})

	// Test single trait
	admin := factory.Build(factory.WithTrait("admin")...)
	if admin.Role != roleAdmin {
		t.Errorf("Admin should have role 'admin', got %s", admin.Role)
	}

	if !admin.Active {
		t.Error("Admin should be active")
	}

	// Test another trait
	inactiveUser := factory.Build(factory.WithTrait("inactive")...)
	if inactiveUser.Active {
		t.Error("Inactive user should not be active")
	}

	if inactiveUser.NoTag != "deactivated" {
		t.Errorf("Inactive user should have NoTag 'deactivated', got %s", inactiveUser.NoTag)
	}
}

func TestMultipleTraits(t *testing.T) {
	factory := gofab.Define[IntegratedTestUser]().
		Trait("admin", func(u *IntegratedTestUser) {
			u.Role = roleAdmin
		}).
		Trait("active", func(u *IntegratedTestUser) {
			u.Active = true
		}).
		Trait("tagged", func(u *IntegratedTestUser) {
			u.NoTag = "tagged-user"
		})

	// Test multiple traits using WithTraits
	user := factory.Build(factory.WithTraits("admin", "active", "tagged")...)

	if user.Role != roleAdmin {
		t.Errorf("User should have role 'admin', got %s", user.Role)
	}

	if !user.Active {
		t.Error("User should be active")
	}

	if user.NoTag != "tagged-user" {
		t.Errorf("User should have NoTag 'tagged-user', got %s", user.NoTag)
	}
}

func TestTraitsWithOverrides(t *testing.T) {
	factory := gofab.Define[IntegratedTestUser]().
		Trait("admin", func(u *IntegratedTestUser) {
			u.Role = roleAdmin
			u.Active = true
		}).
		Trait("inactive", func(u *IntegratedTestUser) {
			u.Active = false
		})

	// Apply admin trait then inactive trait (inactive should override active)
	var builders []gofab.Builder[IntegratedTestUser]
	builders = append(builders, factory.WithTrait("admin")...)
	builders = append(builders, factory.WithTrait("inactive")...)

	user := factory.Build(builders...)

	if user.Role != roleAdmin {
		t.Errorf("User should have role 'admin', got %s", user.Role)
	}

	if user.Active {
		t.Error("User should be inactive (inactive trait should override admin's active setting)")
	}
}

func TestTraitsWithCustomBuilder(t *testing.T) {
	factory := gofab.Define[IntegratedTestUser]().
		Trait("admin", func(u *IntegratedTestUser) {
			u.Role = roleAdmin
			u.NoTag = "admin-tag"
		})

	// Apply trait and then custom builder
	var builders []gofab.Builder[IntegratedTestUser]
	builders = append(builders, factory.WithTrait("admin")...)
	builders = append(builders, func(u *IntegratedTestUser) {
		u.NoTag = "custom-override"
	})

	user := factory.Build(builders...)

	if user.Role != roleAdmin {
		t.Errorf("User should have role 'admin', got %s", user.Role)
	}

	if user.NoTag != "custom-override" {
		t.Errorf("NoTag should be overridden by custom builder, got %s", user.NoTag)
	}
}

func TestTraitsInBuildList(t *testing.T) {
	factory := gofab.Define[IntegratedTestUser]().
		Trait("admin", func(u *IntegratedTestUser) {
			u.Role = roleAdmin
			u.Active = true
		})

	admins := factory.BuildList(3, factory.WithTrait("admin")...)

	if len(admins) != 3 {
		t.Errorf("Expected 3 admins, got %d", len(admins))
	}

	for i, admin := range admins {
		if admin.Role != roleAdmin {
			t.Errorf("Admin %d should have role 'admin', got %s", i, admin.Role)
		}

		if !admin.Active {
			t.Errorf("Admin %d should be active", i)
		}
	}
}

func TestNonExistentTrait(t *testing.T) {
	factory := gofab.Define[IntegratedTestUser]()

	// Using a non-existent trait should return nil/empty builders
	builders := factory.WithTrait("nonexistent")
	if len(builders) != 0 {
		t.Errorf("Non-existent trait should return empty builders, got %d", len(builders))
	}

	// Should still be able to build without errors
	user := factory.Build(factory.WithTrait("nonexistent")...)
	if user.ID == 0 {
		t.Error("User should still be built with auto-generated fields")
	}
}
