package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/sde"
	"testing"
)

func TestSetupStuff(t *testing.T) {
	var err error
	SDE, err = sde.Load("dust.sde")
	if err != nil {
		t.Fatalf("Couldn't open SDE file:( %v", err.Error())
	}
}

func TestDisplayAttributes(t *testing.T) {
	a, err := SDE.GetType(364022)
	if err != nil {
		t.Fatalf("Type not resolved: %v", err.Error())
	}
	w := SDETypeToWraped(a)
	d := w.GetDisplayAttributes()
	for _, v := range d {
		/*		if v.GetValue() == nil {
				fmt.Printf("DisplayAttribute: %v was nil.  Why? %#v\n", v)
			}*/
		fmt.Printf("DisplayAttribute: %v with value '%v'\n", v.DisplayName, v.GetValue())
		if v.Err != nil {
			fmt.Printf("\n===================================================")
			fmt.Printf("\nError in DisplayAttirbute struct: %v\n", v.Err.Error())
			fmt.Printf("====================================================\n")
		}
	}
}
