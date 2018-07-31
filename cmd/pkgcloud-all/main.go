package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mlafeldt/pkgcloud"
)

var usage = "Usage: pkgcloud-all user/repo[/distro/version]\n"

func main() {
	log.SetFlags(0)

	flag.Usage = func() { fmt.Fprintf(os.Stderr, usage) }
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal(usage)
	}

	repo := flag.Args()[0]

	client, err := pkgcloud.NewClient("")
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	next := func() (*pkgcloud.PaginatedPackages, error) {
		return client.PaginatedAll(repo)
	}
	var packages []pkgcloud.Package
	for next != nil {
		paginatedPackages, err := next()
		if err != nil {
			log.Fatalf("pagination error: %s\n", err)
		}
		packages = append(packages, paginatedPackages.Packages...)
		for _, p := range packages {
			fmt.Printf("%+v\n", p)
		}
		next = paginatedPackages.Next
	}
}
