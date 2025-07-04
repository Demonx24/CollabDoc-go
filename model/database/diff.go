package database

import "encoding/xml"

type DocumentDoc struct {
	XMLName xml.Name `xml:"document"`
	Body    Body     `xml:"body"`
}

type Body struct {
	Paragraphs []Paragraph `xml:"p"`
}

type Paragraph struct {
	Runs []Run `xml:"r"`
}

type Run struct {
	Texts []Text `xml:"t"`
}

type Text struct {
	Content string `xml:",chardata"`
}
