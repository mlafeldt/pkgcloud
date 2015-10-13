// List of supported distributions with their internal names and IDs.
//
// Run `make generate` to update the list according to the Packagecloud API. By
// embedding the returned data, we save an expensive API call.
//
// See https://packagecloud.io/docs/api#resource_distributions

package pkgcloud

//go:generate bash -c "./gendistros.py | gofmt > gendistros.go"

import (
	"errors"
	"strings"
)

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
