# gofab

A simple and flexible factory pattern implementation for Go testing, inspired by factory_bot.

## Features

- Type-safe factory definitions using generics
- Automatic field population with struct tags
- Build single or multiple instances
- Sequence builder for unique IDs
- Traits for reusable attribute sets
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
    Name  string `gofab:"name"`
    Email string `gofab:"email"`
    Age   int    `gofab:"range:18,65"`
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

### Sequence Builder for Unique IDs

Use the `Sequence` builder for auto-incrementing values:

```go
type Order struct {
    ID     int
    Number string
}

orderFactory := gofab.Define[Order](
    gofab.Sequence(
        func(o *Order, n int64) { o.ID = int(n) },
        func(n int64) int64 { return n + 1 },
    ),
)

order1 := orderFactory.Build() // ID: 1
order2 := orderFactory.Build() // ID: 2
order3 := orderFactory.Build() // ID: 3
```

### Traits for Reusable Attribute Sets

Define named sets of attributes that can be reused:

```go
type User struct {
    Email  string
    Role   string
    Active bool
}

userFactory := gofab.Define[User]().
    Trait("admin", func(u *User) {
        u.Role = "admin"
        u.Active = true
    }).
    Trait("suspended", func(u *User) {
        u.Active = false
    })

// Use single trait
admin := userFactory.Build(
    append(userFactory.WithTrait("admin"),
        func(u *User) {
            u.Email = "alice@company.com"
        })...
)

// Combine multiple traits
suspendedAdmin := userFactory.Build(
    userFactory.WithTraits("admin", "suspended")...
)
```

## Examples

See `example_test.go` for practical examples including:
- Basic usage with auto-generated data
- Using traits for test scenarios
- E-commerce product testing with discounts

## License

MIT