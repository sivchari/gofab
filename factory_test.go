package gofab_test

import (
	"testing"

	"github.com/sivchari/gofab"
)

const (
	roleUser  = "user"
	roleAdmin = "admin"
)

type IntegratedTestUser struct {
	ID     int    `gofab:"range:1,1000"`
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

	// Check that all admins have IDs
	for i, admin := range admins {
		if admin.ID == 0 {
			t.Errorf("Admin %d should have an ID", i)
		}

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
