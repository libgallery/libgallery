package e621

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/time/rate"
	"spiderden.net/go/libgallery"
	"spiderden.net/go/libgallery/drivers/internal"
)

func New() libgallery.Driver {
	client := retryablehttp.NewClient()
	client.Logger = &internal.NoLogger{}
	return &implementation{
		client:  client.StandardClient(),
		limiter: *rate.NewLimiter(2, 1),
	}
}

type implementation struct {
	client  *http.Client
	limiter rate.Limiter
}

func (i *implementation) getJSON(url string, h *http.Client, target interface{}) error {
	i.limiter.Wait(context.Background())
	return internal.GetJSON(url, h, target)
}

func (i *implementation) Search(query string, page uint64) ([]libgallery.Post, error) {
	const reqbase = "https://e621.net/posts.json?tags=%s&page=%v"
	url := fmt.Sprintf(reqbase, url.QueryEscape(query), page+1)

	var response struct {
		Posts []post `json:"posts"`
	}

	err := i.getJSON(url, i.client, &response)
	if err != nil {
		if herr, ok := err.(*internal.HTTPError); ok {
			if herr.Code() == http.StatusGone {
				return []libgallery.Post{}, nil
			} else {
				return []libgallery.Post{}, err
			}
		} else {
			return []libgallery.Post{}, err
		}
	}

	var libposts []libgallery.Post

	for _, v := range response.Posts {
		ptime, err := time.Parse(time.RFC3339, v.CreatedAt)
		if err != nil {
			return libposts, err
		}
		libposts = append(libposts, libgallery.Post{
			ID:          strconv.FormatUint(v.ID, 10),
			Date:        ptime,
			NSFW:        (v.Rating != "s"),
			Description: v.Description,
			Score:       v.Score.Total,
			Tags:        v.Tags.toTagString(),
			Uploader:    strconv.FormatUint(v.UploaderID, 10),
			Source:      v.Source,
		})
	}

	return libposts, err

}

func (i *implementation) File(id string) ([]io.ReadCloser, error) {
	const reqbase = "https://e621.net/posts/%s.json"
	url := fmt.Sprintf(reqbase, id)

	var response struct {
		Post post `json:"post"`
	}
	err := i.getJSON(url, i.client, &response)
	if err != nil {
		return []io.ReadCloser{}, err
	}

	filereader, err := internal.GetReadCloser(response.Post.File.URL, i.client)
	if err != nil {
		return []io.ReadCloser{}, err
	}

	return []io.ReadCloser{filereader}, nil

}

// No API access, will need to implement scraping.
func (i *implementation) Comments(id string) ([]libgallery.Comment, error) {
	return []libgallery.Comment{}, nil
}

func (i *implementation) Name() string {
	return "e621"
}
