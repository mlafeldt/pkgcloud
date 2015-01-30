package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var serviceURL = "https://packagecloud.io"
var token = os.Getenv("PACKAGECLOUD_TOKEN")

func pushPackage(dest, source string) error {
	s := strings.Split(dest, "/")
	user, repo, distro, version := s[0], s[1], s[2], s[3]

	// TODO: support other formats like .rpm
	distroVersionName := distro + "/" + version
	distroVersionID, ok := debDistros[distroVersionName]
	if !ok {
		return errors.New("Unknown distro version: " + distroVersionName)
	}

	endpoint := fmt.Sprintf("%s/api/v1/repos/%s/%s/packages.json",
		serviceURL, user, repo)
	extraParams := map[string]string{
		"package[distro_version_id]": strconv.Itoa(distroVersionID),
	}
	request, err := NewFileUploadRequest(endpoint, extraParams,
		"package[package_file]", source)
	if err != nil {
		return err
	}
	request.SetBasicAuth(token, "")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// TODO: unmarshal JSON to get error
	if resp.StatusCode != http.StatusCreated {
		msg := fmt.Sprintf("HTTP error: %s\nHTTP body: %s\n",
			resp.Status, body)
		return errors.New(msg)
	}

	return nil
}

func main() {
	log.SetFlags(0)

	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("Usage: pkgcloud-push user/repo/distro/version /path/to/packages")
	}

	dest := flag.Args()[0]
	// TODO: upload packages in parallel
	for _, source := range flag.Args()[1:] {
		fmt.Printf("Pushing %s to %s ...\n", source, dest)
		if err := pushPackage(dest, source); err != nil {
			log.Fatal(err)
		}
	}
}
