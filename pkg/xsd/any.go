package xsd

import (
	"encoding/xml"
)

type Any struct {
	XMLName         xml.Name `xml:"http://www.w3.org/2001/XMLSchema any"`
	Namespace       string   `xml:"namespace,attr"`
	MinOccurs       string   `xml:"minOccurs,attr"`
	MaxOccurs       string   `xml:"maxOccurs,attr"`
	ProcessContents string   `xml:"processContents,attr"`
	schema          *Schema  `xml:"-"`
}

func (a *Any) compile(sch *Schema, parentElement *Element) {
	a.schema = sch
}
