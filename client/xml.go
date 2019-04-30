package client

import "encoding/xml"

type DeviceListInfo struct {
	XMLName xml.Name `xml:"devicelist"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Device  []struct {
		Text            string `xml:",chardata"`
		Identifier      string `xml:"identifier,attr"`
		ID              string `xml:"id,attr"`
		Functionbitmask string `xml:"functionbitmask,attr"`
		Fwversion       string `xml:"fwversion,attr"`
		Manufacturer    string `xml:"manufacturer,attr"`
		Productname     string `xml:"productname,attr"`
		Present         string `xml:"present"`
		Name            string `xml:"name"`
		Temperature     *struct {
			Text    string `xml:",chardata"`
			Celsius int    `xml:"celsius"`
			Offset  int    `xml:"offset"`
		} `xml:"temperature"`
		Hkr struct {
			Text            string `xml:",chardata"`
			Tist            int    `xml:"tist"`
			Tsoll           int    `xml:"tsoll"`
			Absenk          int    `xml:"absenk"`
			Komfort         int    `xml:"komfort"`
			Lock            string `xml:"lock"`
			Devicelock      string `xml:"devicelock"`
			Errorcode       string `xml:"errorcode"`
			Batterylow      string `xml:"batterylow"`
			Windowopenactiv string `xml:"windowopenactiv"`
			Battery         int    `xml:"battery"`
			Nextchange      struct {
				Text      string `xml:",chardata"`
				Endperiod string `xml:"endperiod"`
				Tchange   string `xml:"tchange"`
			} `xml:"nextchange"`
			Summeractive  string `xml:"summeractive"`
			Holidayactive string `xml:"holidayactive"`
		} `xml:"hkr"`
	} `xml:"device"`
	Group []struct {
		Text            string `xml:",chardata"`
		Identifier      string `xml:"identifier,attr"`
		ID              string `xml:"id,attr"`
		Functionbitmask string `xml:"functionbitmask,attr"`
		Fwversion       string `xml:"fwversion,attr"`
		Manufacturer    string `xml:"manufacturer,attr"`
		Productname     string `xml:"productname,attr"`
		Present         string `xml:"present"`
		Name            string `xml:"name"`
		Hkr             struct {
			Text            string `xml:",chardata"`
			Tist            string `xml:"tist"`
			Tsoll           string `xml:"tsoll"`
			Absenk          string `xml:"absenk"`
			Komfort         string `xml:"komfort"`
			Lock            string `xml:"lock"`
			Devicelock      string `xml:"devicelock"`
			Errorcode       string `xml:"errorcode"`
			Batterylow      string `xml:"batterylow"`
			Windowopenactiv string `xml:"windowopenactiv"`
			Battery         string `xml:"battery"`
			Nextchange      struct {
				Text      string `xml:",chardata"`
				Endperiod string `xml:"endperiod"`
				Tchange   string `xml:"tchange"`
			} `xml:"nextchange"`
			Summeractive  string `xml:"summeractive"`
			Holidayactive string `xml:"holidayactive"`
		} `xml:"hkr"`
		Groupinfo struct {
			Text           string `xml:",chardata"`
			Masterdeviceid string `xml:"masterdeviceid"`
			Members        string `xml:"members"`
		} `xml:"groupinfo"`
	} `xml:"group"`
}
