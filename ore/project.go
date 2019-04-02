package ore

import (
	"fmt"
)

const HOME = "https://ore.spongepowered.org"
const projects = "https://ore.spongepowered.org/api/v1/users/%s"

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
	for _, p := range o.Projects {
		if p.ID == name || p.Name == name {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("project not found: %s:%s", owner, name)
}
