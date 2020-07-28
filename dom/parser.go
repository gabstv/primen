package dom

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

// ParseXMLString parses a XLM string.
func ParseXMLString(xmlstr string) ([]Node, error) {
	rdr := strings.NewReader(xmlstr)
	return parsexml(rdr)
}

func parsexml(rdr io.Reader) ([]Node, error) {
	d := xml.NewDecoder(rdr)
	vroot := Element("_root", nil)
	err := parsenode(d, vroot)
	if err != nil {
		return nil, err
	}
	return vroot.Children(), nil
}

//TODO: maybe consider namespaces (?)

func parsenode(d *xml.Decoder, cur ElementNode) error {
	for {
		tkn, err := d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch tkn.(type) {
		case xml.StartElement:
			re := tkn.(xml.StartElement)
			tagname := re.Name.Local
			attrs := make(map[string]string)
			for _, v := range re.Attr {
				attrs[v.Name.Local] = v.Value
			}
			enode := Element(tagname, attrs)
			cur.Append(enode)
			if err := parsenode(d, enode); err != nil {
				return err
			}
		case xml.EndElement:
			re := tkn.(xml.EndElement)
			if re.Name.Local != cur.TagName() {
				// what do?
				return errors.New("closing a tag that was never opened: " + re.Name.Local)
			}
			// this tag closed
			return nil
		case xml.CharData:
			re := tkn.(xml.CharData)
			// discard if empty
			txt := strings.TrimSpace(string(re))
			if txt != "" {
				cur.Append(Text(txt))
			}
		default:
			// ignore xml.Comment
			// ignore xml.ProcInst
			// ignore xml.Directive
		}
	}
}
