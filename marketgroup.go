package main

import (
	"github.com/THUNDERGROOVE/SDETool/sde"
	"log"
)

var MarketGroupCache *MarketGroup

func init() {
	log.Println("Starting MarketGroupCache collection")
	go func() {
		MarketGroupCache = GenerateMarketGroupTree()
		log.Println("Ending MarketGroupCache collection")
	}()
}

type MarketGroup struct {
	SubGroups []*MarketGroup
	Parent    *MarketGroup // If Parent is nil we assume that MarketGroup is the Master
	Types     []*sde.SDEType
	Name      string
	TypeID    int
}

func NewMarketGroup() *MarketGroup {
	g := new(MarketGroup)
	g.SubGroups = make([]*MarketGroup, 0)
	g.Types = make([]*sde.SDEType, 0)
	return g
}

func GenerateMarketGroupTree() *MarketGroup {
	if SDE == nil {
		log.Fatalf("GenerateMarketGroupTree cannot be called when the SDE is not initialized")
	}
	master := NewMarketGroup()
	master.Name = "Master"
	for _, t := range SDE.Types {
		if v, ok := t.Attributes["parentCategoryID"]; ok {
			pid := v.(int)
			cname, _ := t.Attributes["categoryName"].(string)
			if pid == (-1) { // Found first class MarketGroups
				a := NewMarketGroup()
				a.Name = cname
				a.Parent = master
				a.TypeID = t.TypeID
				master.SubGroups = append(master.SubGroups, a)
				fillNextTree(a)
			}
		}
	}
	return master
}

func fillNextTree(m *MarketGroup) {
	for _, t := range SDE.Types {
		if v, ok := t.Attributes["categoryID"]; ok {
			pid := v.(int)
			if pid == m.TypeID {
				m.Types = append(m.Types, t)
			}
		}

		if v, ok := t.Attributes["parentCategoryID"]; ok {
			pid := v.(int)
			cname, _ := t.Attributes["categoryName"].(string)
			if pid == m.TypeID {
				a := NewMarketGroup()
				a.Name = cname
				a.Parent = m
				a.TypeID = t.TypeID
				m.SubGroups = append(m.SubGroups, a)
				fillNextTree(a)
			}
		}
	}
}
