package xsd

import (
	"encoding/xml"
)

type Restriction struct {
	XMLName          xml.Name        `xml:"http://www.w3.org/2001/XMLSchema restriction"`
	Base             reference       `xml:"base,attr"`
	AttributesDirect []Attribute     `xml:"attribute"`
	EnumsDirect      []Enumeration   `xml:"enumeration"`
	SimpleContent    *SimpleContent  `xml:"simpleContent"`
	WhiteSpace       *WhiteSpace     `xml:"whiteSpace"`
	Pattern          *Pattern        `xml:"pattern"`
	MinExclusive     *MinExclusive   `xml:"minExclusive"`
	MaxExclusive     *MaxExclusive   `xml:"maxExclusive"`
	MinInclusive     *MinInclusive   `xml:"minInclusive"`
	MaxInclusive     *MaxInclusive   `xml:"maxInclusive"`
	TotalDigits      *TotalDigits    `xml:"totalDigits"`
	FractionDigits   *FractionDigits `xml:"fractionDigits"`
	Length           *Length         `xml:"length"`
	MinLength        *MinLength      `xml:"minLength"`
	MaxLength        *MaxLength      `xml:"maxLength"`
	schema           *Schema         `xml:"-"`
	typ              Type
}

func (r *Restriction) compile(sch *Schema, parentElement *Element) {
	r.schema = sch
	for idx := range r.AttributesDirect {
		attribute := &r.AttributesDirect[idx]
		attribute.compile(sch)
	}
	if r.SimpleContent != nil {
		r.SimpleContent.compile(sch, parentElement)
	}

	if r.Base == "" {
		panic("Not implemented: xsd:extension/@base empty, cannot extend unknown type")
	}

	r.typ = sch.findReferencedType(r.Base)
	if r.typ == nil {
		panic("Cannot build xsd:extension: unknown type: " + string(r.Base))
	}
	r.typ.compile(sch, parentElement)
}

func (r *Restriction) Attributes() []Attribute {
	result := make([]Attribute, 0)
	if r.typ != nil {
		result = append(result, r.typ.Attributes()...)
	}
	if r.SimpleContent != nil {
		result = append(result, r.SimpleContent.Attributes()...)
	}
	result = deduplicateAttributes(append(result, r.AttributesDirect...))

	return injectSchemaIntoAttributes(r.schema, result)
}

func (r *Restriction) Enums() []Enumeration {
	return r.EnumsDirect
}

type WhiteSpace struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema whiteSpace"`
	Value   string   `xml:"value,attr"`
	Fixed   *bool    `xml:"fixed,attr"`
}

type Pattern struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema pattern"`
	Value   string   `xml:"value,attr"`
}

type MinExclusive struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema minExclusive"`
	Value   string   `xml:"value,attr"`
	Fixed   *bool    `xml:"fixed,attr"`
}

type MaxExclusive struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema maxExclusive"`
	Value   string   `xml:"value,attr"`
	Fixed   *bool    `xml:"fixed,attr"`
}

type MinInclusive struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema minInclusive"`
	Value   string   `xml:"value,attr"`
	Fixed   *bool    `xml:"fixed,attr"`
}

type MaxInclusive struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema maxInclusive"`
	Value   string   `xml:"value,attr"`
	Fixed   *bool    `xml:"fixed,attr"`
}

type TotalDigits struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema totalDigits"`
	Value   string   `xml:"value,attr"`
	Fixed   *bool    `xml:"fixed,attr"`
}

type FractionDigits struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema fractionDigits"`
	Value   string   `xml:"value,attr"`
	Fixed   *bool    `xml:"fixed,attr"`
}

type Length struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema length"`
	Value   string   `xml:"value,attr"`
	Fixed   *bool    `xml:"fixed,attr"`
}

type MinLength struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema minLength"`
	Value   string   `xml:"value,attr"`
	Fixed   *bool    `xml:"fixed,attr"`
}

type MaxLength struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema maxLength"`
	Value   string   `xml:"value,attr"`
	Fixed   *bool    `xml:"fixed,attr"`
}
