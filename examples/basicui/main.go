package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("ui.xml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	xd := xml.NewDecoder(f)
	for {
		tkn, err := xd.Token()
		if err != nil {
			println(err.Error())
			return
		}
		switch tkn.(type) {
		case xml.StartElement:
			v := tkn.(xml.StartElement)
			fmt.Println("xml.StartElement", v)
		case xml.EndElement:
			v := tkn.(xml.EndElement)
			fmt.Println("xml.EndElement", v)
		case xml.CharData:
			v := tkn.(xml.CharData)
			fmt.Println("xml.CharData", v)
		case xml.Comment:
			v := tkn.(xml.Comment)
			fmt.Println("xml.Comment", v)
		case xml.ProcInst:
			v := tkn.(xml.ProcInst)
			fmt.Println("xml.ProcInst", v)
		case xml.Directive:
			v := tkn.(xml.Directive)
			fmt.Println("xml.Directive", v)
		}
		//fmt.Println(tkn)
	}
}
