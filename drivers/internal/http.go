package internal

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

func Get(url string, h *http.Client) (*[]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "libgallery/0.0.0")

	response, err := h.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, NewHTTPError(response.StatusCode)
	}

	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &b, err
}

func GetReadCloser(url string, h *http.Client) (io.ReadCloser, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "libgallery/0.0.0")

	response, err := h.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, NewHTTPError(response.StatusCode)
	}

	return response.Body, err
}

// Does a request and unmarshals the request
// into target.
func GetJSON(url string, h *http.Client, target interface{}) error {
	response, err := GetReadCloser(url, h)
	if err != nil {
		return err
	}

	defer response.Close()

	err = json.NewDecoder(response).Decode(target)
	if err != nil {
		return err
	}

	return err
}

func GetXML(url string, h *http.Client, target interface{}) error {
	response, err := GetReadCloser(url, h)
	if err != nil {
		return err
	}

	defer response.Close()

	dc := xml.NewDecoder(response)
	dc.Strict = false

	dc.Decode(target)
	if err != nil {
		return err
	}

	return nil
}

type NoLogger struct{}

func (n *NoLogger) Printf(_ string, _ ...interface{}) {}
