package danbooru

import "spiderden.net/go/libgallery"

func init() {
	libgallery.Register("danbooru", New("Danbooru", "danbooru.donmai.us"))
	libgallery.Register("thebub.club", New("The Bub Club", "thebub.club"))
}
