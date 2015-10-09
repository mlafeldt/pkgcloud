package pkgcloud

import "testing"

func TestDebDistroIDs(t *testing.T) {
	var tests = map[string]int{
		"ubuntu/warty":    1,
		"ubuntu/trusty":   20,
		"debian/jessie":   25,
		"any/any":         35,
		"linuxmint/petra": 157,
		"raspbian/buster": 156,
	}
	for name, id := range tests {
		if debDistroIDs[name] != id {
			t.Errorf("distro id of %s != %d", name, id)
		}
	}
}

func TestRpmDistroIDs(t *testing.T) {
	var tests = map[string]int{
		"el/7":         140,
		"fedora/22":    147,
		"scientific/5": 138,
		"ol/7":         146,
	}
	for name, id := range tests {
		if rpmDistroIDs[name] != id {
			t.Errorf("distro id of %s != %d", name, id)
		}
	}
}
