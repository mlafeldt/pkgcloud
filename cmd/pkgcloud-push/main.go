package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	target, err := newTarget(flag.Args()[0])
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	packages := flag.Args()[1:]

	client := pkgcloud.NewClient("")
	resc := make(chan string)
	errc := make(chan error)

	fmt.Printf("Pushing %d package(s) to %s ...\n", len(packages), target)
	for _, pkg := range packages {
		go func(pkg string) {
			if err := client.CreatePackage(target.repo, target.distro, pkg); err != nil {
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
