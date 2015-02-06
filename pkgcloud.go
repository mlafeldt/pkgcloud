package pkgcloud

import (
	"encoding/json"
	"errors"
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

func NewClient(token string) (*Client, error) {
	if token == "" {
		token = os.Getenv("PACKAGECLOUD_TOKEN")
		if token == "" {
			return nil, errors.New("PACKAGECLOUD_TOKEN unset")
		}
	}
	return &Client{token}, nil
}

func decodeResponse(status int, body []byte) error {
	switch status {
	case http.StatusOK, http.StatusCreated:
		return nil
	case http.StatusUnauthorized, http.StatusNotFound:
		return fmt.Errorf("HTTP status: %s", http.StatusText(status))
	case 422: // Unprocessable Entity
		var v map[string][]string
		if err := json.Unmarshal(body, &v); err != nil {
			return err
		}
		for _, messages := range v {
			for _, msg := range messages {
				// Only return the very first error message
				return errors.New(msg)
			}
			break
		}
		return fmt.Errorf("invalid HTTP body: %s", body)
	default:
		return fmt.Errorf("unexpected HTTP status: %d", status)
	}
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

	return decodeResponse(resp.StatusCode, body)
}
