# gofab

A simple and flexible factory pattern implementation for Go testing, inspired by factory_bot.

## Features

- Type-safe factory definitions using generics
- Automatic field population with struct tags
- Build single or multiple instances
- Sequence support for unique IDs
- Customizable field overrides
- Integration with faker for realistic test data

## Installation

```bash
go get github.com/sivchari/gofab
```

## Quick Start

### Basic Usage with Auto-population

```go
type User struct {
    ID    int    `gofab:"sequence"`
    Name  string `gofab:"name"`
    Email string `gofab:"email"`
}

// Build a single user
user := gofab.Build[User]()

// Build with custom values
admin := gofab.Build[User](func(u *User) {
    u.Name = "Admin User"
})
```

### Building Multiple Instances

```go
// Create 5 users with auto-generated data
users := gofab.BuildList[User](5)

// Create 3 admins with custom role
admins := gofab.BuildList[User](3, func(u *User) {
    u.Role = "admin"
})
```

### Using Factory Definitions

```go
// Define a factory with defaults
userFactory := gofab.Define[User](func(u *User) {
    u.Role = "user"
    u.Active = true
})

// Build single instance
user := userFactory.Build()

// Build multiple instances
users := userFactory.BuildList(10)

// Override defaults
admin := userFactory.Build(func(u *User) {
    u.Role = "admin"
})
```

## Struct Tags

gofab supports various struct tags for automatic field population:

- `gofab:"sequence"` - Auto-incrementing ID
- `gofab:"name"` - Random person name
- `gofab:"email"` - Random email address
- `gofab:"phone"` - Random phone number
- `gofab:"company"` - Random company name
- `gofab:"address"` - Random address
- `gofab:"sentence:N"` - Random sentence with N words
- `gofab:"word"` - Random word
- `gofab:"range:min,max"` - Random number in range
- `gofab:"-"` - Skip auto-population

## Advanced Examples

### Testing with Factories

```go
func TestUserService(t *testing.T) {
    // Create test data
    users := gofab.BuildList[User](3)
    
    // Test with specific scenarios
    inactiveUser := gofab.Build[User](func(u *User) {
        u.Active = false
    })
    
    // Your test logic here
}
```

### Sequence for Unique IDs

```go
type Order struct {
    ID     int `gofab:"sequence"`
    Number string
}

order1 := gofab.Build[Order]() // ID: 1
order2 := gofab.Build[Order]() // ID: 2
order3 := gofab.Build[Order]() // ID: 3
```

### Traits for Reusable Attribute Sets

Define named sets of attributes that can be reused:

```go
type User struct {
    ID          int      `gofab:"sequence"`
    Name        string   `gofab:"name"`
    Email       string   `gofab:"email"`
    Role        string
    Permissions []string
    Active      bool
}

userFactory := gofab.Define[User]().
    Trait("admin", func(u *User) {
        u.Role = "admin"
        u.Permissions = []string{"read", "write", "delete"}
        u.Active = true
    }).
    Trait("inactive", func(u *User) {
        u.Active = false
    })

// Use single trait
admin := userFactory.Build(userFactory.WithTrait("admin")...)

// Combine multiple traits
inactiveAdmin := userFactory.Build(
    userFactory.WithTraits("admin", "inactive")...
)

// Traits with custom overrides
customAdmin := userFactory.Build(append(
    userFactory.WithTrait("admin"),
    func(u *User) {
        u.Name = "Super Admin"
    },
)...)
```

### AfterBuild Callbacks

Use AfterBuild to automatically process instances after building:

```go
type BlogPost struct {
    Title     string `gofab:"sentence:3"`
    Body      string `gofab:"sentence:10"`
    Slug      string // Generated in AfterBuild
    WordCount int    // Calculated in AfterBuild
}

postFactory := gofab.Define[BlogPost]().
    AfterBuild(func(p *BlogPost) {
        // Generate slug from title
        p.Slug = strings.ToLower(strings.ReplaceAll(p.Title, " ", "-"))
    }).
    AfterBuild(func(p *BlogPost) {
        // Calculate word count
        p.WordCount = len(strings.Fields(p.Body))
    })

// AfterBuild callbacks run automatically
post := postFactory.Build()
// post.Slug and post.WordCount are automatically set
```

Multiple callbacks are executed in order:

```go
factory := gofab.Define[User]().
    AfterBuild(func(u *User) {
        u.Username = strings.ToLower(u.Name)
    }).
    AfterBuild(func(u *User) {
        u.Email = fmt.Sprintf("%s@example.com", u.Username)
    })
```

## License

MIT