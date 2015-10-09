// List of supported distributions with their internal names and IDs.
//
// Run `make generate` to update the list according to the Packagecloud API. By
// embedding the returned data, we save an expensive API call.
//
// See https://packagecloud.io/docs/api#resource_distributions

package pkgcloud

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/

import (
	"encoding/json"
	"errors"
	"strings"
)

type distribution struct {
	DisplayName string `json:"display_name"`
	IndexName   string `json:"index_name"`
	Versions    []version
}

type version struct {
	ID            int    `json:"id"`
	DisplayName   string `json:"display_name"`
	IndexName     string `json:"index_name"`
	VersionNumber string `json:"version_number"`
}

func makeMap(distros []distribution) map[string]int {
	m := make(map[string]int)
	for _, d := range distros {
		for _, v := range d.Versions {
			k := strings.Join([]string{d.IndexName, v.IndexName}, "/")
			m[k] = v.ID
		}
	}
	return m
}

var debDistroIDs, rpmDistroIDs map[string]int

func init() {
	data, err := Asset("assets/distributions.json")
	if err != nil {
		panic(err)
	}

	var pkgTypes map[string][]distribution
	if err := json.Unmarshal(data, &pkgTypes); err != nil {
		panic(err)
	}

	debDistroIDs = makeMap(pkgTypes["deb"])
	rpmDistroIDs = makeMap(pkgTypes["rpm"])
}

const (
	extensionDeb = ".deb"
	extensionDsc = ".dsc"
	extensionRpm = ".rpm"
	extensionGem = ".gem"
)

func distroID(ext, name string) (int, error) {
	switch strings.ToLower(ext) {
	case extensionDeb, extensionDsc:
		if id, ok := debDistroIDs[name]; ok {
			return id, nil
		}
	case extensionRpm:
		if id, ok := rpmDistroIDs[name]; ok {
			return id, nil
		}
	case extensionGem:
		return 0, errors.New("RubyGem packages have no distribution")
	default:
		return 0, errors.New("invalid file extension: " + ext)
	}
	return 0, errors.New("invalid distro name: " + name)
}
