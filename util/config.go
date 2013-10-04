package util

import (
	"time"
)

var tz, _ = time.LoadLocation("Europe/Bratislava")

var C = &Config{
	Tz:     tz,
	Title:  "TOGY",
	Footer: `&copy 2013, Icons by <a href="http://glyphicons.com/">Glyphicons</a>`,
}

type Config struct {
	Tz     *time.Location
	Title  string
	Footer string
}
