package main

import "strings"

var Slovar = []string{"lesbian", "xxx", "mylf", "rarbg", "sex", "uncensored", "anal", "480p", "porn", "erotic", "tushyraw", "tushy", "suck", "big butt", "cock", "onlyfans", "hot mom", "hardcore", "creampie", "virtualtaboo.com", "penthouse", "playboy", "wowgirls", "marc dorcel", "mygirlfriendsbustyfriend", "massagegirls18", "hentai", ".dmg", "x-art", "секс", "pure18", "[jav]", "masturbat", ".mpg", "[sunshine]"}

func filtr(name string) bool {
	name = strings.ToLower(name)
	fl := false
	for _, v := range Slovar {
		if strings.Contains(name, v) {
			fl = true
			break
		}
	}
	return fl
}
