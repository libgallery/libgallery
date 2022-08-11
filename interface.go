/*
libgallery provides an interface for accessing boorus,
and other resources in a booru-like fashion.
*/
package libgallery

import (
	"io"
	"time"
)

// Driver interface provides a method to access boorus, or
// a view of resources in a booru-like fashion (called virtual
// implementations)
type Driver interface {
	Name() string // The proper name of a site, i.e. "Danbooru"
	// Search() can only accept space-separated snake-case tags.
	// If a booru doesn't support it you will need to
	// transmogrify the query.
	Search(string, uint64) ([]Post, error) // Tags, and page number.
	File(string) (Files, error)            // Fetches a file with a given ID.
	Comments(string) ([]Comment, error)    // Fetches the comments from a given ID.
}

// Post has various fields for what a post may contain.
// All zero values aside from Score and NSFW should be
// considered not-implemented on the site or otherwise
// unfetchable. Uploader should be the ID, not the
// username (unless usernames are used as IDs like
// Reddit)
type Post struct {
	URL         string
	ID          string    `json:"id"`
	NSFW        bool      `json:"nsfw"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Tags        string    `json:"tags"`
	Uploader    string    `json:"uploader"`
	Source      []string  `json:"source"`
	Score       int64     `json:"score"`
}

type Files []io.ReadCloser

func (f *Files) Close() {
	for _, v := range *f {
		v.Close()
	}
}

type Comment struct {
	Author string    `json:"author"`
	Body   string    `json:"body"`
	Date   time.Time `json:"date"`
	Score  int64     `json:"score"`
}
