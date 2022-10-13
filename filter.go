package main

import (
	"regexp"
	"strings"
)

var Slovar = []string{"lesbian", "xxx", "mylf", "rarbg", "sex", "uncensored", "anal", "480p", "porn", "erotic", "tushyraw",
	"tushy", "suck", "big butt", "cock", "onlyfans", "hot mom", "hardcore", "creampie", "virtualtaboo.com", "penthouse", "playboy",
	"wowgirls", "marc dorcel", "mygirlfriendsbustyfriend", "massagegirls18", "hentai", ".dmg", "x-art", "секс", "pure18", "[jav]",
	"masturbat", ".mpg", "[sunshine]", "girlsoutwest.com", "herlimit", " ass ", "teenpies.com", "セックス", "69av", ".wmv", "ssis",
	"所偷拍", "theav.cc", "cumshots", "youiv.net", ".flv", "rartv", "avmans.com", "thz.la", "lolly", "erai-raws", "macos", "aniua",
	"backroomcastingcouch.com", "shirasudon", "mywife"}

func filtr(name string) bool {
	name = strings.ToLower(name)
	fl := false

	pattern := `[a-z]{2,4}(-|)\d{3}`
	matched, _ := regexp.Match(pattern, []byte(name))

	if matched {
		fl = true
		return fl
	}
	for _, v := range Slovar {
		if strings.Contains(name, v) {
			fl = true
			break
		}
	}
	return fl
}
