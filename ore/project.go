package ore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const HOME = "https://ore.spongepowered.org"
const projects = "https://ore.spongepowered.org/api/v1/projects/%s"

type Project struct {
	HREF string `json:"href"`
}

func GetProjectPage(id string) (string, error) {
	url := fmt.Sprintf(projects, strings.ToLower(id))
	r, e := http.Get(url)
	if e != nil {
		return "", e
	}
	defer r.Body.Close()

	var p Project
	e = json.NewDecoder(r.Body).Decode(&p)
	if e != nil {
		return "", e
	}
	return HOME + p.HREF, nil
}

func GetVersionPage(id, version string) (string, error) {
	v, e := GetVersion(strings.ToLower(id), version)
	if e != nil {
		return "", e
	}
	return HOME + v.HREF, nil
}
