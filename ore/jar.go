package ore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	URL string `json:"url"`
}

var download = `https://ore.spongepowered.org/api/projects/%s/versions/%s/download`

func GetJar(id, version string) (io.ReadCloser, error) {
	url := fmt.Sprintf(download, id, version)
	return getJar(id, version, url, 0)
}

func getJar(id, version, url string, flag int) (io.ReadCloser, error) {
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}

	content := r.Header.Get("Content-Type")
	if content == "application/json" {
		defer r.Body.Close()
		if flag != 0 {
			return nil, fmt.Errorf("jar not found: %s:%s", id, version)
		}

		var message Message
		e := json.NewDecoder(r.Body).Decode(&message)
		if e != nil {
			return nil, e
		}

		if message.URL == "" {
			return nil, fmt.Errorf("jar not found: %s:%s", id, version)
		}

		return getJar(id, version, message.URL, 1)
	}

	if content != "application/octet-stream" {
		defer r.Body.Close()
		return nil, fmt.Errorf("jar not found: %s:%s", id, version)
	}

	return r.Body, nil
}
