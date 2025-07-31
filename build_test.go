package gofab_test

import (
	"testing"

	"github.com/sivchari/gofab"
)

type TestUser struct {
	ID      int    `gofab:"sequence"`
	Name    string `gofab:"name"`
	Email   string `gofab:"email"`
	Phone   string `gofab:"phone"`
	Company string `gofab:"company"`
	Address string `gofab:"address"`
	Age     int    `gofab:"range:18,65"`
	Bio     string `gofab:"sentence:3"`
	Title   string `gofab:"word"`
	Skip    string `gofab:"-"`
	NoTag   string
}

func TestBuildSingle(t *testing.T) {
	user := gofab.Build[TestUser]()
	
	// Check that tagged fields are populated
	if user.ID == 0 {
		t.Error("ID should be populated by sequence")
	}
	if user.Name == "" {
		t.Error("Name should be populated")
	}
	if user.Email == "" {
		t.Error("Email should be populated")
	}
	if user.Phone == "" {
		t.Error("Phone should be populated")
	}
	if user.Company == "" {
		t.Error("Company should be populated")
	}
	if user.Address == "" {
		t.Error("Address should be populated")
	}
	if user.Age == 0 {
		t.Error("Age should be populated with range")
	}
	if user.Bio == "" {
		t.Error("Bio should be populated with sentence")
	}
	if user.Title == "" {
		t.Error("Title should be populated with word")
	}
	
	// Check that skipped and no-tag fields are not populated
	if user.Skip != "" {
		t.Error("Skip field should remain empty")
	}
	if user.NoTag != "" {
		t.Error("NoTag field should remain empty")
	}
	
	// Check age range
	if user.Age < 18 || user.Age > 65 {
		t.Errorf("Age should be between 18 and 65, got %d", user.Age)
	}
}

func TestBuildWithCustomizer(t *testing.T) {
	user := gofab.Build[TestUser](func(u *TestUser) {
		u.NoTag = "Custom Value"
		u.Name = "Override Name"
	})
	
	// Check that customizer overrides work
	if user.NoTag != "Custom Value" {
		t.Errorf("NoTag should be customized, got %s", user.NoTag)
	}
	if user.Name != "Override Name" {
		t.Errorf("Name should be overridden, got %s", user.Name)
	}
	
	// Check that other fields are still auto-generated
	if user.Email == "" {
		t.Error("Email should still be auto-generated")
	}
}

func TestBuildList(t *testing.T) {
	users := gofab.BuildList[TestUser](3)
	
	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}
	
	// Check that all users have different sequence IDs
	ids := make(map[int]bool)
	for _, user := range users {
		if ids[user.ID] {
			t.Errorf("Duplicate ID found: %d", user.ID)
		}
		ids[user.ID] = true
		
		// Check that each user has populated fields
		if user.Name == "" || user.Email == "" {
			t.Error("All users should have populated fields")
		}
	}
}

func TestBuildListWithCustomizer(t *testing.T) {
	users := gofab.BuildList[TestUser](2, func(u *TestUser) {
		u.NoTag = "All Same"
	})
	
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
	
	for i, user := range users {
		if user.NoTag != "All Same" {
			t.Errorf("User %d should have customized NoTag", i)
		}
		if user.Name == "" {
			t.Errorf("User %d should have auto-generated Name", i)
		}
	}
}

func TestSequenceIncrement(t *testing.T) {
	user1 := gofab.Build[TestUser]()
	user2 := gofab.Build[TestUser]()
	user3 := gofab.Build[TestUser]()
	
	if user2.ID != user1.ID+1 {
		t.Errorf("Expected user2.ID to be %d, got %d", user1.ID+1, user2.ID)
	}
	if user3.ID != user2.ID+1 {
		t.Errorf("Expected user3.ID to be %d, got %d", user2.ID+1, user3.ID)
	}
}

type InvalidRangeUser struct {
	BadRange1 int `gofab:"range:invalid"`
	BadRange2 int `gofab:"range:10"`
	BadRange3 int `gofab:"range:100,50"` // max < min
	GoodRange int `gofab:"range:1,10"`
}

func TestInvalidTags(t *testing.T) {
	user := gofab.Build[InvalidRangeUser]()
	
	// Invalid range tags should fallback to default range (1-100)
	if user.BadRange1 < 1 || user.BadRange1 > 100 {
		t.Errorf("BadRange1 should fallback to 1-100, got %d", user.BadRange1)
	}
	if user.BadRange2 < 1 || user.BadRange2 > 100 {
		t.Errorf("BadRange2 should fallback to 1-100, got %d", user.BadRange2)
	}
	if user.BadRange3 < 1 || user.BadRange3 > 100 {
		t.Errorf("BadRange3 should fallback to 1-100, got %d", user.BadRange3)
	}
	
	// Good range should work correctly
	if user.GoodRange < 1 || user.GoodRange > 10 {
		t.Errorf("GoodRange should be 1-10, got %d", user.GoodRange)
	}
}