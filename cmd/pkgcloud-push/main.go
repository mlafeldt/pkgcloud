package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mlafeldt/pkgcloud"
)

func main() {
	log.SetFlags(0)

	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("Usage: pkgcloud-push user/repo/distro/version /path/to/packages")
	}

	client := pkgcloud.NewClient("")
	target := flag.Args()[0]
	packages := flag.Args()[1:]

	resc := make(chan string)
	errc := make(chan error)

	fmt.Printf("Pushing %d package(s) to %s ...\n", len(packages), target)
	for _, name := range packages {
		go func(name string) {
			if err := client.PushPackage(target, name); err != nil {
				errc <- fmt.Errorf("%s ... %s", name, err)
				return
			}
			resc <- fmt.Sprintf("%s ... OK", name)
		}(name)
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
