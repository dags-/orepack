package ore

const (
	snapshotTimeStamp = "20190101.23000"
	snapshotUpdated   = "2019010123000"
	SnapshotSuffix    = "-20190101.23000-1"
)

type Metadata struct {
	GroupID    string     `xml:"groupId"`
	ArtifactID string     `xml:"artifactId"`
	Version    string     `xml:"version"`
	Versioning Versioning `xml:"versioning"`
}

type Versioning struct {
	Snapshot         Snapshot          `xml:"snapshot"`
	LastUpdated      string            `xml:"lastUpdated"`
	SnapshotVersions []SnapshotVersion `xml:"snapshotVersions"`
}

type Snapshot struct {
	Timestamp   string `xml:"timestamp"`
	BuildNumber int    `xml:"buildNumber"`
}

type SnapshotVersion struct {
	Extension string `xml:"extension"`
	Value     string `xml:"value"`
	Updated   string `xml:"updated"`
}

func NewMetaData(owner, name, version string) *Metadata {
	return &Metadata{
		GroupID:    "com.orepack." + owner,
		ArtifactID: name,
		Version:    version,
		Versioning: Versioning{
			Snapshot: Snapshot{
				Timestamp:   snapshotTimeStamp,
				BuildNumber: 1,
			},
			LastUpdated: snapshotUpdated,
			SnapshotVersions: []SnapshotVersion{
				{
					Extension: "jar",
					Value:     version,
					Updated:   snapshotUpdated,
				},
				{
					Extension: "pom",
					Value:     version,
					Updated:   snapshotUpdated,
				},
			},
		},
	}
}
