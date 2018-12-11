package ore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var versions = `https://ore.spongepowered.org/api/v1/projects/%s/versions`

type Version struct {
	Name     string `json:"name"`
	MD5      string `json:"md5"`
	HREF     string `json:"href"`
	Approved bool   `json:"staffApproved"`
}

func GetVersion(id, version string) (*Version, error) {
	r, e := http.Get(fmt.Sprintf(versions, id))
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, fmt.Errorf(r.Status)
	}

	var versions []Version
	e = json.NewDecoder(r.Body).Decode(&versions)
	if e != nil {
		return nil, e
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("version not found: %s:%s", id, version)
	}

	if strings.ToUpper(version) == "LATEST" {
		latest := &versions[0]
		return latest, nil
	}

	for _, v := range versions {
		if v.Name == version {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("version not found: %s:%s", id, version)
}
