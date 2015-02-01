package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mlafeldt/pkgcloud"
)

var usage = "Usage: pkgcloud-push user/repo/distro/version /path/to/packages\n"

func main() {
	log.SetFlags(0)

	flag.Usage = func() { fmt.Fprintf(os.Stderr, usage) }
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal(usage)
	}
	target := flag.Args()[0]
	packages := flag.Args()[1:]

	var repo, distro string
	elems := strings.Split(target, "/")
	if len(elems) == 2 {
		repo = strings.Join(elems[0:2], "/")
	} else if len(elems) == 4 {
		repo = strings.Join(elems[0:2], "/")
		distro = strings.Join(elems[2:4], "/")
	} else {
		log.Fatal("error: invalid target: " + target)
	}

	client := pkgcloud.NewClient("")
	resc := make(chan string)
	errc := make(chan error)

	fmt.Printf("Pushing %d package(s) to %s ...\n", len(packages), target)
	for _, pkg := range packages {
		go func(pkg string) {
			if err := client.CreatePackage(repo, distro, pkg); err != nil {
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
