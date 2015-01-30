package pkgcloud

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var serviceURL = "https://packagecloud.io"

type Client struct {
	token string
}

func NewClient(token string) *Client {
	if token == "" {
		token = os.Getenv("PACKAGECLOUD_TOKEN")
	}
	return &Client{token}
}

func (c Client) PushPackage(target, name string) error {
	s := strings.Split(target, "/")
	user, repo, distro, version := s[0], s[1], s[2], s[3]

	id, err := distroID(filepath.Ext(name), distro+"/"+version)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/api/v1/repos/%s/%s/packages.json",
		serviceURL, user, repo)
	extraParams := map[string]string{
		"package[distro_version_id]": strconv.Itoa(id),
	}
	request, err := NewFileUploadRequest(endpoint, extraParams,
		"package[package_file]", name)
	if err != nil {
		return err
	}
	request.SetBasicAuth(c.token, "")

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
