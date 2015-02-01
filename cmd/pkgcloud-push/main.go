package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mlafeldt/pkgcloud"
)

func main() {
	log.SetFlags(0)

	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("Usage: pkgcloud-push user/repo/distro/version /path/to/packages")
	}

	target := flag.Args()[0]
	packages := flag.Args()[1:]

	var user, repo, distro string
	elems := strings.Split(target, "/")
	if len(elems) == 2 {
		user = elems[0]
		repo = elems[1]
	} else if len(elems) == 4 {
		user = elems[0]
		repo = elems[1]
		distro = elems[2] + "/" + elems[3]
	} else {
		log.Fatal("invalid target: " + target)
	}

	client := pkgcloud.NewClient("")
	resc := make(chan string)
	errc := make(chan error)

	fmt.Printf("Pushing %d package(s) to %s ...\n", len(packages), target)
	for _, pkg := range packages {
		go func(pkg string) {
			if err := client.CreatePackage(user, repo, distro, pkg); err != nil {
				errc <- fmt.Errorf("%s ... %s", pkg, err)
				return
			}
			resc <- fmt.Sprintf("%s ... OK", pkg)
		}(pkg)
	}

	failure := false
	for range packages {
		select {
		case res := <-resc:
			log.Println(res)
		case err := <-errc:
			log.Println(err)
			failure = true
		}
	}
	if failure {
		os.Exit(1)
	}
}
