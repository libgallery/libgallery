// Driver for Danbooru and compatible APIs.
package danbooru

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"spiderden.net/go/libgallery"
	"spiderden.net/go/libgallery/drivers/internal"
)

type implementation struct {
	host   string
	name   string
	client *http.Client
}

// Creates a new Driver with a given configuration and
// host. You can use this to make wrappers around
// this driver for other Danbooru-compatible sites.

// past1000 sets whether the driver should look past
// a thousand pages.
func New(name string, host string) libgallery.Driver {
	client := retryablehttp.NewClient()
	client.Logger = &internal.NoLogger{}
	return &implementation{
		client: client.StandardClient(),
		name:   name,
		host:   host,
	}
}

func (i *implementation) Search(query string, page uint64) ([]libgallery.Post, error) {
	const reqbase = "https://%s/posts.json?tags=%s&page=%v"
	url := fmt.Sprintf(reqbase, i.host, url.QueryEscape(query), page+1)

	var response []post
	err := internal.GetJSON(url, i.client, &response)
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

	var posts []libgallery.Post

	for _, v := range response {
		ptime, err := time.Parse(time.RFC3339, v.CreatedAt)
		if err != nil {
			return []libgallery.Post{}, err
		}
		posts = append(posts, libgallery.Post{
			Uploader: strconv.FormatUint(uint64(v.ID), 10),
			Tags:     v.Tags,
			Date:     ptime,
			Source:   []string{v.Source},
			ID:       fmt.Sprintf("%v", v.ID),
			NSFW:     true,
			Score:    int64(v.Score),
		})
	}

	return posts, err
}

func (i *implementation) File(id string) (libgallery.Files, error) {
	const reqbase = "https://%s/posts/%v.json"
	url := fmt.Sprintf(reqbase, i.host, id)

	var response post
	err := internal.GetJSON(url, i.client, &response)
	if err != nil {
		return []io.ReadCloser{}, err
	}

	rc, err := internal.GetReadCloser(response.LargeFileURL, i.client)
	if err != nil {
		return []io.ReadCloser{}, err
	}

	return []io.ReadCloser{rc}, nil
}

func (i *implementation) Comments(id string) ([]libgallery.Comment, error) {
	const reqbase = "https://%s/comments.json?group_by=comment&post_id=%v"
	url := fmt.Sprintf(reqbase, i.host, id)

	var response comments
	err := internal.GetJSON(url, i.client, &response)
	if err != nil {
		return []libgallery.Comment{}, err
	}

	var comments []libgallery.Comment
	for _, v := range response {
		ptime, err := time.Parse(time.RFC3339, v.CreatedAt)
		if err != nil {
			return []libgallery.Comment{}, err
		}

		str := libgallery.Comment{
			Author: strconv.FormatUint(uint64(v.ID), 10),
			Body:   v.Body,
			Date:   ptime,
			Score:  v.Score,
		}

		comments = append(comments, str)
	}

	return comments, err
}

func (i *implementation) Name() string {
	return i.name
}
