// List of supported distributions with their internal names and IDs.
// Generated from https://packagecloud.io/docs/api#resource_distributions

package pkgcloud

var debDistros = map[string]int{
	"ubuntu/warty":    1,
	"ubuntu/hoary":    2,
	"ubuntu/breezy":   3,
	"ubuntu/dapper":   4,
	"ubuntu/edgy":     5,
	"ubuntu/feisty":   6,
	"ubuntu/gutsy":    7,
	"ubuntu/hardy":    8,
	"ubuntu/intrepid": 9,
	"ubuntu/jaunty":   10,
	"ubuntu/karmic":   11,
	"ubuntu/lucid":    12,
	"ubuntu/maverick": 13,
	"ubuntu/natty":    14,
	"ubuntu/oneiric":  15,
	"ubuntu/precise":  16,
	"ubuntu/quantal":  17,
	"ubuntu/raring":   18,
	"ubuntu/saucy":    19,
	"ubuntu/trusty":   20,
	"ubuntu/utopic":   142,
	"debian/etch":     21,
	"debian/lenny":    22,
	"debian/squeeze":  23,
	"debian/wheezy":   24,
	"debian/jessie":   25,
	"any/any":         35,
}

var dscDistros = debDistros

var rpmDistros = map[string]int{
	"el/5":         26,
	"el/6":         27,
	"el/7":         140,
	"fedora/14":    28,
	"fedora/15":    29,
	"fedora/16":    30,
	"fedora/17":    31,
	"fedora/18":    32,
	"fedora/19":    33,
	"fedora/20":    34,
	"fedora/21":    143,
	"scientific/5": 138,
	"scientific/6": 139,
	"scientific/7": 141,
}
