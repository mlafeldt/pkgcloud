// List of supported distributions with their internal names and IDs.
// Generated from https://packagecloud.io/docs/api#resource_distributions

package pkgcloud

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	debDistroIDs, rpmDistroIDs map[string]int
)

func init() {
	data, err := Asset("assets/distributions.json")
	if err != nil {
		panic(err)
	}
	// TODO: parse json here
	fmt.Print(string(data))
	os.Exit(0)
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
