package ore

import (
	"encoding/json"
	"fmt"
	"strings"
)

const owners = "https://ore.spongepowered.org/api/v1/users/%s"

type Owner struct {
	Projects []Project `json:"projects"`
}

func GetOwner(owner string) (*Owner, error) {
	url := fmt.Sprintf(owners, strings.ToLower(owner))
	r, e := HttpGet(url)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	var o Owner
	e = json.NewDecoder(r.Body).Decode(&o)
	if e != nil {
		return nil, e
	}
	return &o, nil
}
