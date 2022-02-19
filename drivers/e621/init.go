package e621

import "spiderden.net/go/libgallery"

func init() {
	libgallery.Register("e621", New())
}
