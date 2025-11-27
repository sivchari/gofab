package gofab_test

import (
	"testing"

	"github.com/sivchari/gofab"
)

// Test types for association tests.
type Author struct {
	ID   int    `gofab:"sequence"`
	Name string `gofab:"name"`
}

type Category struct {
	ID   int    `gofab:"sequence"`
	Name string `gofab:"word"`
}

type Post struct {
	ID         int    `gofab:"sequence"`
	Title      string `gofab:"sentence:3"`
	AuthorID   int
	CategoryID int
}

type PostRelations struct {
	Author   Author
	Category Category
}

func TestAssociationFactory_Build(t *testing.T) {
	authorFactory := gofab.Define[Author]()
	categoryFactory := gofab.Define[Category]()

	postFactory := gofab.DefineWithAssociations[Post, PostRelations](
		func() PostRelations {
			return PostRelations{
				Author:   authorFactory.Build(),
				Category: categoryFactory.Build(),
			}
		},
		func(p *Post, r PostRelations) {
			p.AuthorID = r.Author.ID
			p.CategoryID = r.Category.ID
		},
	)

	post := postFactory.Build()

	if post.ID == 0 {
		t.Error("Post ID should be auto-generated")
	}

	if post.Title == "" {
		t.Error("Post Title should be auto-generated")
	}

	if post.AuthorID == 0 {
		t.Error("Post AuthorID should be set from association")
	}

	if post.CategoryID == 0 {
		t.Error("Post CategoryID should be set from association")
	}
}

func TestAssociationFactory_BuildWithAssociations(t *testing.T) {
	authorFactory := gofab.Define[Author]()
	categoryFactory := gofab.Define[Category]()

	postFactory := gofab.DefineWithAssociations[Post, PostRelations](
		func() PostRelations {
			return PostRelations{
				Author:   authorFactory.Build(),
				Category: categoryFactory.Build(),
			}
		},
		func(p *Post, r PostRelations) {
			p.AuthorID = r.Author.ID
			p.CategoryID = r.Category.ID
		},
	)

	post, relations := postFactory.BuildWithAssociations()

	if post.AuthorID != relations.Author.ID {
		t.Errorf("Post AuthorID (%d) should match Author ID (%d)", post.AuthorID, relations.Author.ID)
	}

	if post.CategoryID != relations.Category.ID {
		t.Errorf("Post CategoryID (%d) should match Category ID (%d)", post.CategoryID, relations.Category.ID)
	}

	if relations.Author.Name == "" {
		t.Error("Author Name should be auto-generated")
	}

	if relations.Category.Name == "" {
		t.Error("Category Name should be auto-generated")
	}
}

func TestAssociationFactory_BuildList(t *testing.T) {
	authorFactory := gofab.Define[Author]()
	categoryFactory := gofab.Define[Category]()

	postFactory := gofab.DefineWithAssociations[Post, PostRelations](
		func() PostRelations {
			return PostRelations{
				Author:   authorFactory.Build(),
				Category: categoryFactory.Build(),
			}
		},
		func(p *Post, r PostRelations) {
			p.AuthorID = r.Author.ID
			p.CategoryID = r.Category.ID
		},
	)

	posts := postFactory.BuildList(3)

	if len(posts) != 3 {
		t.Errorf("Expected 3 posts, got %d", len(posts))
	}

	// Each post should have different associations
	authorIDs := make(map[int]bool)

	for _, post := range posts {
		if post.AuthorID == 0 {
			t.Error("Post AuthorID should be set")
		}

		authorIDs[post.AuthorID] = true
	}

	// All posts should have different author IDs (different associations created each time)
	if len(authorIDs) != 3 {
		t.Errorf("Expected 3 different author IDs, got %d", len(authorIDs))
	}
}

func TestAssociationFactory_WithCustomAssociations(t *testing.T) {
	authorFactory := gofab.Define[Author]()
	categoryFactory := gofab.Define[Category]()

	postFactory := gofab.DefineWithAssociations[Post, PostRelations](
		func() PostRelations {
			return PostRelations{
				Author:   authorFactory.Build(),
				Category: categoryFactory.Build(),
			}
		},
		func(p *Post, r PostRelations) {
			p.AuthorID = r.Author.ID
			p.CategoryID = r.Category.ID
		},
	)

	// Create custom associations
	customAuthor := authorFactory.Build(func(a *Author) {
		a.Name = "Custom Author"
	})
	customCategory := categoryFactory.Build(func(c *Category) {
		c.Name = "Custom Category"
	})

	post := postFactory.WithCustomAssociations(PostRelations{
		Author:   customAuthor,
		Category: customCategory,
	})

	if post.AuthorID != customAuthor.ID {
		t.Errorf("Post AuthorID (%d) should match custom Author ID (%d)", post.AuthorID, customAuthor.ID)
	}

	if post.CategoryID != customCategory.ID {
		t.Errorf("Post CategoryID (%d) should match custom Category ID (%d)", post.CategoryID, customCategory.ID)
	}
}

func TestAssociationFactory_WithTrait(t *testing.T) {
	authorFactory := gofab.Define[Author]()
	categoryFactory := gofab.Define[Category]()

	postFactory := gofab.DefineWithAssociations[Post, PostRelations](
		func() PostRelations {
			return PostRelations{
				Author:   authorFactory.Build(),
				Category: categoryFactory.Build(),
			}
		},
		func(p *Post, r PostRelations) {
			p.AuthorID = r.Author.ID
			p.CategoryID = r.Category.ID
		},
	).Trait("featured", func(p *Post) {
		p.Title = "Featured Post"
	})

	post := postFactory.Build(postFactory.WithTrait("featured")...)

	if post.Title != "Featured Post" {
		t.Errorf("Post Title should be 'Featured Post', got '%s'", post.Title)
	}

	if post.AuthorID == 0 {
		t.Error("Post AuthorID should still be set from association")
	}
}
