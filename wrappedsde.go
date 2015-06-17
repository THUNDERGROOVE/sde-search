package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/sde"
	"log"
	"reflect"
	"strings"
)

/*
TODO:
Certain things like flaylocks require the projectile to be looked up and searched as well for attributes
*/

var MissingAttributes = []string{
	"overHeatingInfo.cooldownTime",
	"mMaxEquipmentBandwidth",
}

type WrappedSDEType struct {
	*sde.SDEType
}

// DisplayAttribute represents the attributes that an attrib_adaptor_* has
type DisplayAttribute struct {
	Units         string
	AttributeName string
	DisplayName   string
	ValueConvert  string
	ValueSource   string
	TypeName      string
	TypeID        int
	Parent        *WrappedSDEType

	Err error
}

// SDETypeToWrapped takes an sde.SDEType and returns a WrappedSDEType
func SDETypeToWraped(t *sde.SDEType) *WrappedSDEType {
	o := new(WrappedSDEType)
	o.SDEType = t
	return o
}

func (w *WrappedSDEType) GetDisplayAttributes() map[string]*DisplayAttribute {
	out := make([]*DisplayAttribute, 0)
	outm := make(map[string]*DisplayAttribute)
	for k, v := range w.Attributes {
		if strings.Contains(k, "displayAttributes.") {
			dattr := parseAttr(v.(int), w)
			out = append(out, dattr)
		}
	}
	for _, v := range out {
		outm[v.AttributeName] = v
	}
	/*	keys := make([]string, 0)
		for _, v := range out {
			keys = append(keys, v.AttributeName)
		}

		sort.Strings(keys)*/

	return outm
}

func parseAttr(v int, parent *WrappedSDEType) *DisplayAttribute {
	out := new(DisplayAttribute)

	if t, err := SDE.GetType(v); err != nil {
		out.Err = err
		return out
	} else {
		out.AttributeName, _ = t.Attributes["attributeName"].(string)
		out.DisplayName, _ = t.Attributes["displayName"].(string)
		out.ValueConvert, _ = t.Attributes["valueConvert"].(string)
		out.ValueSource, _ = t.Attributes["valueSource"].(string)
		out.Units, _ = t.Attributes["Units"].(string)
		out.TypeName = t.TypeName
		out.TypeID = t.TypeID
	}

	out.Parent = parent
	return out
}

// GetValue parses the DisplayAttribute for a value to display.
func (d *DisplayAttribute) GetValue() interface{} {
	if within(d.AttributeName, MissingAttributes) {
		return nil
	}
	var origVal interface{}

	switch d.ValueSource {
	case "EAVS_FROM_TYPE":
		if d.AttributeName == "@hardcoded" {
			log.Printf("Encountered a @hardcoded attribute: %v#%v", d.TypeName, d.TypeID)
			return origVal
		}
		if v, ok := d.Parent.Attributes[d.AttributeName]; ok {
			origVal = v
		} else {
			d.Err = fmt.Errorf("parent had no type that matched DisplayAttribute's:%v#%v AttributeName %v", d.TypeName, d.TypeID, d.AttributeName)
			return nil
		}
	default:
		d.Err = fmt.Errorf("don't know how to handle value source %v", d.ValueSource)
		return nil
	}

	switch d.ValueConvert {
	case "EAVDC_NONE":
		return origVal
	case "EAVDC_CENTIMETER_TO_METER":
		switch t := origVal.(type) {
		case float64:
			origVal = interface{}(t * 0.01)
			return origVal
		default:
			d.Err = fmt.Errorf("EAVDC_CENTIMETER_TO_METER encountered non-supported type: %v", reflect.TypeOf(origVal))
			return origVal
		}
	case "EAVDC_SECONDS_INTERVAL_TO_RATE_PER_MINUTE":
		switch t := origVal.(type) {
		case float64:
			fm, _ := d.Parent.Attributes["mFireMode0.m_eFireMode"]
			return interface{}(calcROF(t, fm.(string)))
		default:
			d.Err = fmt.Errorf("EAVDC_SECONDS_INTERVAL_TO_RATE_PER_MINUTE encountered non-supported type: %v", reflect.TypeOf(origVal))
			return origVal
		}
	default:
		d.Err = fmt.Errorf("don't know how to handle convert value %v defaulting to none.", d.ValueConvert)
		return origVal
	}
	return origVal
}

func (d *DisplayAttribute) GetValueString() string {
	return fmt.Sprintf("%v", d.GetValue())
}

func within(v string, s []string) bool {
	for _, val := range s {
		if v == val {
			return true
		}
	}
	return false
}

func calcROF(i float64, firemode string) float64 {
	switch firemode {
	default:
		return (1 / i) * 60
	}
	return 0
}
