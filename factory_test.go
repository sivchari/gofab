package gofab_test

import (
	"testing"

	"github.com/sivchari/gofab"
)

type IntegratedTestUser struct {
	ID      int    `gofab:"sequence"`
	Name    string `gofab:"name"`
	Email   string `gofab:"email"`
	Role    string // No tag - will be set by factory defaults
	Active  bool   `gofab:"-"` // Skip auto-generation
	NoTag   string             // No tag, no factory default
}

func TestDefineWithAutoGeneration(t *testing.T) {
	// Create factory with defaults that will be applied after auto-generation
	userFactory := gofab.Define[IntegratedTestUser](func(u *IntegratedTestUser) {
		u.Role = "user"  // This will override any auto-generated value
		u.Active = true  // This will set the skipped field
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
	if user.Role != "user" {
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
		u.Role = "admin"
		u.Active = true
	})

	// Override some fields including auto-generated ones
	superAdmin := adminFactory.Build(func(u *IntegratedTestUser) {
		u.Name = "Super Admin"  // Override auto-generated name
		u.NoTag = "Special"     // Set the untouched field
	})

	// Check that overrides work
	if superAdmin.Name != "Super Admin" {
		t.Errorf("Name should be overridden, got %s", superAdmin.Name)
	}
	if superAdmin.NoTag != "Special" {
		t.Errorf("NoTag should be set by override, got %s", superAdmin.NoTag)
	}

	// Check that factory defaults are still applied
	if superAdmin.Role != "admin" {
		t.Errorf("Role should be set by factory, got %s", superAdmin.Role)
	}

	// Check that auto-generation still works for non-overridden fields
	if superAdmin.Email == "" {
		t.Error("Email should still be auto-generated")
	}
}

func TestFactoryBuildList(t *testing.T) {
	adminFactory := gofab.Define[IntegratedTestUser](func(u *IntegratedTestUser) {
		u.Role = "admin"
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
		if admin.Role != "admin" {
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
		u.Role = "user"
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
		if user.Role != "user" {
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

