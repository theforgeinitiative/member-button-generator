package dymo

import "encoding/xml"

type LabelSet struct {
	XMLName     xml.Name      `xml:"LabelSet"`
	Text        string        `xml:",chardata"`
	LabelRecord []LabelRecord `xml:"LabelRecord"`
}

type LabelRecord struct {
	Text       string       `xml:",chardata"`
	ObjectData []LabelField `xml:"ObjectData"`
}

type LabelField struct {
	Text string `xml:",chardata"`
	Name string `xml:"Name,attr"`
}

type Printers struct {
	XMLName  xml.Name  `xml:"Printers"`
	Printers []Printer `xml:"LabelWriterPrinter"`
}

type Printer struct {
	XMLName   xml.Name `xml:"LabelWriterPrinter"`
	Name      string   `xml:"Name"`
	Model     string   `xml:"ModelName"`
	Connected bool     `xml:"IsConnected"`
}
