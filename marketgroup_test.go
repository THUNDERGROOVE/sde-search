package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMarketGroupTree(t *testing.T) {
	tree := GenerateMarketGroupTree()
	recurPrintTree(tree, 0)
}

func recurPrintTree(m *MarketGroup, d int) {
	fmt.Printf("%v> %v\n", strings.Repeat("=", d), m.Name)
	for _, t := range m.Types {
		fmt.Printf("%v> %v\n", strings.Repeat("-", d+1), t.GetName())
	}
	for _, v := range m.SubGroups {
		recurPrintTree(v, d+1)
	}
}
