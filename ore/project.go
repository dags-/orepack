package ore

import (
	"fmt"
	"strings"
)

type Project struct {
	ID   string `json:"pluginId"`
	Name string `json:"name"`
	HREF string `json:"href"`
}

func GetProject(owner string, name string) (*Project, error) {
	o, e := GetOwner(owner)
	if e != nil {
		return nil, e
	}
	name = strings.ToLower(name)
	for _, p := range o.Projects {
		if p.ID == name || strings.ToLower(p.Name) == name {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("project not found: %s:%s", owner, name)
}
