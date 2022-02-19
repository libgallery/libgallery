package e621

import "strings"

func (t *tags) toTagString() string {
	var converted []string
	f := func(array []string, prefix string) {
		for _, v := range array {
			converted = append(converted, prefix+":"+v)
		}
	}

	f(t.General, "general")
	f(t.Species, "species")
	f(t.Character, "character")
	f(t.Artist, "artist")
	f(t.Lore, "lore")
	f(t.Meta, "meta")

	return strings.Join(converted, " ")
}
