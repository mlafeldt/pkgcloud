package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

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

	var wg sync.WaitGroup
	for _, name := range flag.Args()[1:] {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			fmt.Printf("Pushing %s to %s ...\n", name, target)
			if err := client.PushPackage(target, name); err != nil {
				log.Fatal(err)
			}
		}(name)
	}
	wg.Wait()
}
