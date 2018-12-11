package ore

import "encoding/xml"

type Pom struct {
	XMLName      xml.Name `xml:"project"`
	XMLNS        string   `xml:"xmlns,attr"`
	XMLNSXSI     string   `xml:"xmlns:xsi,attr"`
	XSI          string   `xml:"xsi:schemaLocation,attr"`
	ModelVersion string   `xml:"modelVersion"`
	GroupID      string   `xml:"groupId"`
	ArtifactID   string   `xml:"artifactId"`
	Version      string   `xml:"version"`
}

func NewPom(id string, version *Version) *Pom {
	return &Pom{
		XMLNS:        "http://maven.apache.org/POM/4.0.0",
		XMLNSXSI:     "http://www.w3.org/2001/XMLSchema-instance",
		XSI:          "http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd",
		ModelVersion: "4.0.0",
		GroupID:      "com.orepack",
		ArtifactID:   id,
		Version:      version.Name,
	}
}
