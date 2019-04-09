package ore

type Metadata struct {
	ModelVersion string     `xml:"modelVersion,attr"`
	GroupId      string     `xml:"groupId"`
	ArtifactId   string     `xml:"artifactId"`
	Versioning   Versioning `xml:"versioning"`
}

type Versioning struct {
	Release  string   `xml:"release"`
	Versions []string `xml:"versions>version"`
}

func NewMeta(owner, name, version string) *Metadata {
	return &Metadata{
		ModelVersion: "1.0.0",
		GroupId:      "com.orepack" + owner,
		ArtifactId:   name,
		Versioning: Versioning{
			Release:  version,
			Versions: []string{version},
		},
	}
}
