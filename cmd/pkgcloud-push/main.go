package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mlafeldt/pkgcloud"
)

func main() {
	log.SetFlags(0)

	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("Usage: pkgcloud-push user/repo/distro/version /path/to/packages")
	}

	// TODO: upload packages in parallel
	client := pkgcloud.NewClient("")
	target := flag.Args()[0]
	for _, name := range flag.Args()[1:] {
		fmt.Printf("Pushing %s to %s ...\n", name, target)
		if err := client.PushPackage(target, name); err != nil {
			log.Fatal(err)
		}
	}
}
