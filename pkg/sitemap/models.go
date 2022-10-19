package sitemap

import (
	"encoding/xml"
	"time"
)

type Page struct {
	XMLName         xml.Name   `xml:"url"`
	Location        string     `xml:"loc"`
	LastModified    *time.Time `xml:"lastmod,omitempty"`
	ChangeFrequency Frequency  `xml:"changefreq,omitempty"`
	Priority        float64    `xml:"priority,omitempty"`
	Depth           int        `xml:"-"`
	Links           []string   `xml:"-"`
}

var Extensions = []string{".png", ".jpg", ".jpeg", ".tiff", ".pdf", ".txt", ".gif", ".psd", ".ai", "dwg", ".bmp", ".zip", ".tar", ".gzip", ".svg", ".avi", ".mov", ".json", ".xml", ".mp3", ".wav", ".mid", ".ogg", ".acc", ".ac3", "mp4", ".ogm", ".cda", ".mpeg", ".avi", ".swf", ".acg", ".bat", ".ttf", ".msi", ".lnk", ".dll", ".db"}

var FalseUrls = []string{"mailto:", "javascript:", "tel:", "whatsapp:", "callto:", "wtai:", "sms:", "market:", "geopoint:", "ymsgr:", "msnim:", "gtalk:", "skype:"}

var PriorityMap = map[int]float64{
	1: 1.0,
	2: 0.9,
	3: 0.8,
	4: 0.7,
	5: 0.6,
	6: 0.5,
	7: 0.4,
	8: 0.3,
	9: 0.2,
}
