package ore

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	URL string `json:"url"`
}

var download = `https://ore.spongepowered.org/api/projects/%s/versions/%s/download`

func GetJarURL(id, version string) (string, error) {
	url := fmt.Sprintf(download, id, version)
	r, e := http.Get(url)
	if e != nil {
		return "", e
	}
	defer r.Body.Close()

	content := r.Header.Get("Content-Type")
	if content == "application/octet-stream" {
		return url, nil
	}

	if content == "application/json" {
		var message Message

		e = json.NewDecoder(r.Body).Decode(&message)
		if e != nil {
			return "", e
		}

		return message.URL, nil
	}

	return "", fmt.Errorf("jar not found: %s:%s", id, version)
}
