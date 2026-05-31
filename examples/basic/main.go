// Command basic prints a few facts from the live models.dev catalog.
package main

import (
	"context"
	"fmt"
	"sort"

	modelsdev "github.com/promptrails/modelsdotdev-go"
)

func main() {
	cat, err := modelsdev.New().Catalog(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("catalog: %d providers, %d models\n\n",
		len(cat.Providers), len(cat.Models()))

	if m, err := cat.Model("anthropic:claude-opus-4-8"); err == nil {
		fmt.Printf("%s — reasoning=%t, tools=%t", m.Name, m.Reasoning, m.ToolCall)
		if m.Cost.Input != nil {
			fmt.Printf(", input=$%.2f/1M", *m.Cost.Input)
		}
		if m.Cost.CacheRead != nil {
			fmt.Printf(", cache-read=$%.2f/1M", *m.Cost.CacheRead)
		}
		fmt.Println()
	}

	// Cheapest reasoning models by input price.
	type row struct {
		id    string
		price float64
	}
	var rows []row
	for _, m := range cat.Filter(func(m modelsdev.Model) bool {
		return m.Reasoning && m.Cost.Input != nil
	}) {
		rows = append(rows, row{m.Provider + ":" + m.ID, *m.Cost.Input})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].price < rows[j].price })

	fmt.Println("\ncheapest reasoning models:")
	for _, r := range rows[:min(5, len(rows))] {
		fmt.Printf("  $%6.3f/1M  %s\n", r.price, r.id)
	}
}
