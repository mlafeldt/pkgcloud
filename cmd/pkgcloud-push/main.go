package main

import (
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
var token string

func init() {
	token = os.Getenv("PACKAGECLOUD_TOKEN")
}

func pushPackage(dest, source string) error {
	s := strings.Split(dest, "/")
	user, repo, distro, version := s[0], s[1], s[2], s[3]

	distroVersionID := debDistros[distro+"/"+version]

	fmt.Println(user, repo, distro, version, distroVersionID)

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

	fmt.Println(string(body))

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return nil
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("missing args")
	}

	dest := flag.Args()[0]
	source := flag.Args()[1]

	if err := pushPackage(dest, source); err != nil {
		log.Fatal(err)
	}
}
