package pkgcloud

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mlafeldt/pkgcloud/upload"
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

func (c Client) CreatePackage(repo, distro, pkgFile string) error {
	distID, err := distroID(filepath.Ext(pkgFile), distro)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/api/v1/repos/%s/packages.json",
		serviceURL, repo)
	extraParams := map[string]string{
		"package[distro_version_id]": strconv.Itoa(distID),
	}
	request, err := upload.NewRequest(endpoint, extraParams,
		"package[package_file]", pkgFile)
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
		return fmt.Errorf("HTTP error (%d): %s", resp.StatusCode, body)
	}

	return nil
}
