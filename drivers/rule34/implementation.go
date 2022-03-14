package rule34

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"spiderden.net/go/libgallery"
	"spiderden.net/go/libgallery/drivers/internal"
)

type implementation struct {
	client *http.Client
}

func New() libgallery.Driver {
	client := retryablehttp.NewClient()
	client.Logger = &internal.NoLogger{}
	return &implementation{
		client: client.StandardClient(),
	}
}

func (i *implementation) Search(query string, page uint64) ([]libgallery.Post, error) {
	if page > 2000 {
		return []libgallery.Post{}, nil
	}

	const reqbase = "https://api.rule34.xxx/index.php?page=dapi&s=post&q=index&tags=%s&pid=%v"
	url := fmt.Sprintf(reqbase, url.QueryEscape(query), page)

	var response searchResponse
	err := internal.GetXML(url, i.client, &response)
	if err != nil {
		return []libgallery.Post{}, err
	}

	/* The rule34.xxx API only has a success
	   value if there was an error. */
	if response.Success != nil {
		return []libgallery.Post{}, response.Error
	}

	var posts []libgallery.Post

	for _, v := range response.Posts {
		ptime, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", v.CreatedAt)
		if err != nil {
			return []libgallery.Post{}, err
		}

		var source []string
		if v.Source != "" {
			source = append(source, strings.TrimSpace(v.Source))
		}

		score, err := strconv.ParseInt(v.Score, 10, 64)
		if err != nil {
			return []libgallery.Post{}, err
		}

		posts = append(posts, libgallery.Post{
			Tags:   strings.TrimSpace(v.Tags),
			Date:   ptime,
			Source: source,
			ID:     v.ID,
			NSFW:   true,
			Score:  score,
		})
	}

	return posts, err
}

func (i *implementation) File(id string) (libgallery.Files, error) {
	const reqbase = "https://api.rule34.xxx/index.php?page=dapi&s=post&q=index&id="

	var response searchResponse
	err := internal.GetXML(reqbase+id, i.client, &response)
	if err != nil {
		return []io.ReadCloser{}, err
	}

	// Same deal as before.
	if response.Success != nil {
		return []io.ReadCloser{}, response.Error
	}

	rc, err := internal.GetReadCloser(response.Posts[0].FileURL, i.client)
	if err != nil {
		return []io.ReadCloser{}, err
	}

	return []io.ReadCloser{rc}, nil
}

func (i *implementation) Name() string {
	return "rule34.xxx"
}

// Comments are broken on the API, I tried asking on their Discord
// about it twice with no response.
func (i *implementation) Comments(id string) ([]libgallery.Comment, error) {
	return []libgallery.Comment{}, nil
}
