package gofab_test

import (
	"fmt"
	"time"

	"github.com/sivchari/gofab"
)

type Product struct {
	ID        int64
	SKU       string
	Name      string
	Price     int
	CreatedAt time.Time
}

func ExampleSequence() {
	// Define a factory with sequential values
	productFactory := gofab.Define[Product]().
		Default(gofab.Sequence(
			func(p *Product, id int64) { p.ID = id },
			func(n int64) int64 { return n + 1000 }, // Start from 1000
		)).
		Default(gofab.Sequence(
			func(p *Product, sku string) { p.SKU = sku },
			func(n int64) string { return fmt.Sprintf("PROD-%04d", n) }, // PROD-0000, PROD-0001, ...
		)).
		Default(gofab.Sequence(
			func(p *Product, name string) { p.Name = name },
			func(n int64) string { return fmt.Sprintf("Product %d", n) },
		)).
		Default(gofab.Sequence(
			func(p *Product, price int) { p.Price = price },
			func(n int64) int { return int(1000 + n*100) }, // 1000, 1100, 1200, ...
		)).
		Default(func(p *Product) {
			p.CreatedAt = time.Now()
		})

	// Create multiple products
	products := make([]Product, 3)
	for i := range products {
		products[i] = productFactory.Build()
	}

	// Print products
	for _, p := range products {
		fmt.Printf("Product{ID:%d, SKU:%s, Name:%s, Price:%d}\n",
			p.ID, p.SKU, p.Name, p.Price)
	}

	// Output:
	// Product{ID:1000, SKU:PROD-0000, Name:Product 0, Price:1000}
	// Product{ID:1001, SKU:PROD-0001, Name:Product 1, Price:1100}
	// Product{ID:1002, SKU:PROD-0002, Name:Product 2, Price:1200}
}
