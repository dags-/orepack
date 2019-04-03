package ore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	URL string `json:"post"`
}

var download = `https://ore.spongepowered.org/api/projects/%s/versions/%s/download`

func GetJar(id, version string) (io.ReadCloser, error) {
	url := fmt.Sprintf(download, id, version)
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}

	content := r.Header.Get("Content-Type")
	if content == "application/octet-stream" {
		return r.Body, nil
	}

	defer r.Body.Close()

	if content == "application/json" {
		var message Message

		e = json.NewDecoder(r.Body).Decode(&message)
		if e != nil {
			return nil, e
		}

		p, e := http.Post(message.URL, "application/json", nil)
		if e != nil {
			return nil, e
		}

		return p.Body, nil
	}

	return nil, fmt.Errorf("jar not found: %s:%s", id, version)
}
