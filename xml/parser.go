package xml

import (
	"encoding/xml"
	"jtso/logger"
	"regexp"
	"strings"
)

type RawData struct {
	RtrName    string
	Family     string
	IfDesc     *Ifdesc
	IfList     *Iflist
	HwInfo     *Hw
	LacpInfo   *Lacp
	LacpDigest *LacpDigest
}

// Struct for unmarshalling version
type Version struct {
	XMLName xml.Name `xml:"software-information"`
	Model   string   `xml:"product-model"`
	Ver     string   `xml:"junos-version"`
}

// Structs for unmarshalling interfaces descriptions
type Ifdesc struct {
	XMLName   xml.Name `xml:"interface-information"`
	Physicals []Phy    `xml:"physical-interface"`
	Logicals  []Log    `xml:"logical-interface"`
}

type Phy struct {
	XMLName xml.Name `xml:"physical-interface"`
	Name    string   `xml:"name"`
	Desc    string   `xml:"description"`
}

type Log struct {
	XMLName xml.Name `xml:"logical-interface"`
	Name    string   `xml:"name"`
	Desc    string   `xml:"description"`
}

// Structs for unmarshalling interfaces descriptions
type Iflist struct {
	XMLName   xml.Name  `xml:"interface-information"`
	Physicals []PhyList `xml:"physical-interface"`
}

type PhyList struct {
	XMLName  xml.Name  `xml:"physical-interface"`
	Name     string    `xml:"name"`
	Logicals []LogList `xml:"logical-interface"`
}

type LogList struct {
	XMLName xml.Name `xml:"logical-interface"`
	Name    string   `xml:"name"`
}

// structs for umarshalling chassis hw
type Hw struct {
	XMLName xml.Name `xml:"chassis-inventory"`
	Chassis Chassis  `xml:"chassis"`
}

type Chassis struct {
	XMLName xml.Name `xml:"chassis"`
	Desc    string   `xml:"description"`
	Modules []Module `xml:"chassis-module"`
}

type Module struct {
	XMLName xml.Name `xml:"chassis-module"`
	Name    string   `xml:"name"`
	Desc    string   `xml:"description"`
	SubMods []SubMod `xml:"chassis-sub-module"`
}

type SubMod struct {
	XMLName    xml.Name    `xml:"chassis-sub-module"`
	Name       string      `xml:"name"`
	Desc       string      `xml:"description"`
	SubSubMods []SubSubMod `xml:"chassis-sub-sub-module"`
}

type SubSubMod struct {
	XMLName       xml.Name       `xml:"chassis-sub-sub-module"`
	Name          string         `xml:"name"`
	Desc          string         `xml:"description"`
	SubSubSubMods []SubSubSubMod `xml:"chassis-sub-sub-sub-module"`
}

type SubSubSubMod struct {
	XMLName xml.Name `xml:"chassis-sub-sub-sub-module"`
	Name    string   `xml:"name"`
	Desc    string   `xml:"description"`
}

type Lacp struct {
	XMLName xml.Name  `xml:"lacp-interface-information-list"`
	LacpInt []LacpInt `xml:"lacp-interface-information"`
}

type LacpInt struct {
	XMLName   xml.Name    `xml:"lacp-interface-information"`
	LacpHead  LacpHead    `xml:"lag-lacp-header"`
	LacpProto []LacpProto `xml:"lag-lacp-protocol"`
}

type LacpHead struct {
	XMLName xml.Name `xml:"lag-lacp-header"`
	LagName string   `xml:"aggregate-name"`
}

type LacpProto struct {
	XMLName xml.Name `xml:"lag-lacp-protocol"`
	Name    string   `xml:"name"`
}

// easy to parse LACP data
type LacpDigest struct {
	LacpMap map[string]string
}

// Parsing function for version
func ParseVersion(s string) (*Version, error) {
	logger.HandlePanic()
	var i Version

	// Convert XML string to byte array
	b := []byte(s)

	// voor debug om de family te vinden
	logger.Log.Infof("debug response, %s", s)

	// Unmarshal XML based on QFX or EX or ACX structure
	if strings.Contains(s, "<multi-routing-engine-results>") {
		// QFX/EX device format (with multi-routing-engine-results)
		var multiResult struct {
			Items []struct {
				SoftwareInformation struct {
					//HostName     string `xml:"host-name"`
					JunosVersion string `xml:"junos-version"`
					ProductName  string `xml:"product-name"`
				} `xml:"software-information"`
			} `xml:"multi-routing-engine-item"`
		}
		err := xml.Unmarshal(b, &multiResult)
		if err != nil {
			return nil, err
		}
		if len(multiResult.Items) > 0 {
			localModelRegexPattern := regexp.MustCompile(`^[aA-zZ]*`)
			i.Model = localModelRegexPattern.FindString(multiResult.Items[0].SoftwareInformation.ProductName)

			//i.Model = multiResult.Items[0].SoftwareInformation.HostName
			i.Ver = multiResult.Items[0].SoftwareInformation.JunosVersion
		}

	} else {
		// ACX device format (direct software-information)
		var singleResult struct {
			HostName     string `xml:"host-name"`
			JunosVersion string `xml:"junos-version"`
		}
		err := xml.Unmarshal(b, &singleResult)
		if err != nil {
			return nil, err
		}
		i.Model = singleResult.HostName
		i.Ver = singleResult.JunosVersion
	}

	return &i, nil
}

// Parsing function for interfaces description
func ParseIfdesc(s string) (*Ifdesc, error) {
	logger.HandlePanic()
	var i Ifdesc
	// convert in byte array
	b := []byte(s)
	// unmarshall xml string
	err := xml.Unmarshal(b, &i)
	return &i, err
}

// Parsing function for interfaces terse
func ParseIflist(s string) (*Iflist, error) {
	logger.HandlePanic()
	var i Iflist
	// convert in byte array
	b := []byte(s)
	// unmarshall xml string
	err := xml.Unmarshal(b, &i)
	return &i, err
}

// Parsing function for chassis hw
func ParseChassis(s string) (*Hw, error) {
	logger.HandlePanic()
	var i Hw
	// convert in byte array
	b := []byte(s)
	// unmarshall xml string
	err := xml.Unmarshal(b, &i)
	return &i, err
}

// Parsing function for Lacp interface
func ParseLacp(s string) (*Lacp, *LacpDigest, error) {
	logger.HandlePanic()
	var i Lacp
	data := new(LacpDigest)
	data.LacpMap = make(map[string]string)
	// convert in byte array
	b := []byte(s)
	// unmarshall xml string
	err := xml.Unmarshal(b, &i)
	// now parse reply
	if err == nil {
		for _, l := range i.LacpInt {
			for _, n := range l.LacpProto {
				data.LacpMap[n.Name] = l.LacpHead.LagName
			}
		}
	}

	return &i, data, err
}
