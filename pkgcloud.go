// Package pkgcloud allows you to talk to the packagecloud API.
// See https://packagecloud.io/docs/api
package pkgcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/mlafeldt/pkgcloud/upload"
	"github.com/tomnomnom/linkheader"
)

//go:generate bash -c "./gendistros.py supportedDistros | gofmt > distros.go"

// ServiceURL is the URL of packagecloud's API.
const ServiceURL = "https://packagecloud.io/api/v1"

const UserAgent = "pkgcloud Go client"

// A Client is a packagecloud client.
type Client struct {
	token string
}

// NewClient creates a packagecloud client. API requests are authenticated
// using an API token. If no token is passed, it will be read from the
// PACKAGECLOUD_TOKEN environment variable.
func NewClient(token string) (*Client, error) {
	if token == "" {
		token = os.Getenv("PACKAGECLOUD_TOKEN")
		if token == "" {
			return nil, errors.New("PACKAGECLOUD_TOKEN unset")
		}
	}
	return &Client{token}, nil
}

// decodeResponse checks http status code and tries to decode json body
func decodeResponse(resp *http.Response, respJson interface{}) error {
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return json.NewDecoder(resp.Body).Decode(respJson)
	case http.StatusUnauthorized, http.StatusNotFound:
		return fmt.Errorf("HTTP status: %s", http.StatusText(resp.StatusCode))
	case 422: // Unprocessable Entity
		var v map[string][]string
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return err
		}
		for _, messages := range v {
			for _, msg := range messages {
				// Only return the very first error message
				return errors.New(msg)
			}
			break
		}
		return fmt.Errorf("invalid HTTP body")
	default:
		return fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}
}

// CreatePackage pushes a new package to packagecloud.
func (c Client) CreatePackage(repo, distro, pkgFile string) error {
	var extraParams map[string]string
	if distro != "" {
		distID, ok := supportedDistros[distro]
		if !ok {
			return fmt.Errorf("invalid distro name: %s", distro)
		}
		extraParams = map[string]string{
			"package[distro_version_id]": strconv.Itoa(distID),
		}
	}

	endpoint := fmt.Sprintf("%s/repos/%s/packages.json", ServiceURL, repo)
	request, err := upload.NewRequest(endpoint, extraParams, "package[package_file]", pkgFile)
	if err != nil {
		return err
	}
	request.SetBasicAuth(c.token, "")
	request.Header.Add("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return decodeResponse(resp, &struct{}{})
}

type Package struct {
	Name               string    `json:"name"`
	CreatedAt          time.Time `json:"created_at,string"`
	Epoch              int       `json:"epoch"`
	Scope              string    `json:"scope"`
	Private            bool      `json:"private"`
	UploaderName       string    `json:"uploader_name"`
	Indexed            bool      `json:"indexed"`
	RepositoryHtmlUrl  string    `json:"repository_html_url,string"`
	DownloadDetailsUrl string    `json:"downloads_detail_url"`
	DownloadSeriesUrl  string    `json:"downloads_series_url"`
	DownloadCountUrl   string    `json:"downloads_count_url"`
	PromoteUrl         string    `json:"promote_url"`
	DestroyUrl         string    `json:"destroy_url"`
	Filename           string    `json:"filename"`
	DistroVersion      string    `json:"distro_version"`
	Version            string    `json:"version"`
	Release            string    `json:"release"`
	Type               string    `json:"type"`
	PackageUrl         string    `json:"package_url"`
	PackageHtmlUrl     string    `json:"package_html_url"`
}

// All list all packages in repository
func (c Client) All(repo string) ([]Package, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/packages.json", ServiceURL, repo)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.token, "")
	req.Header.Add("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var packages []Package
	err = decodeResponse(resp, &packages)
	return packages, err
}

// Destroy removes package from repository.
//
// repo should be full path to repository
// (e.g. youruser/repository/ubuntu/xenial).
func (c Client) Destroy(repo, packageFilename string) error {
	endpoint := fmt.Sprintf("%s/repos/%s/%s", ServiceURL, repo, packageFilename)

	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.token, "")
	req.Header.Add("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return decodeResponse(resp, &struct{}{})
}

// Search searches packages from repository.
// repo should be full path to repository
// (e.g. youruser/repository/ubuntu/xenial).
// q: The query string to search for package filename. If empty string is passed, all packages are returned
// filter: Search by package type (RPMs, Debs, DSCs, Gem, Python) - Ignored when dist != ""
// dist: The name of the distribution the package is in. (i.e. ubuntu, el/6) - Overrides filter.
// perPage: The number of packages to return from the results set. If nothing passed, default is 30
func (c Client) Search(repo, q, filter, dist string, perPage int) ([]Package, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/search.json?q=%s&filter=%s&dist=%s&per_page=%d", ServiceURL, repo, q, filter, dist, perPage)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.token, "")
	req.Header.Add("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var packages []Package
	err = decodeResponse(resp, &packages)
	return packages, err
}

type Paginated struct {
	Total      int
	PerPage    int
	MaxPerPage int
}

type PaginatedPackages struct {
	Packages []Package
	Next     func() (*PaginatedPackages, error)
	Paginated
}

func (c *Client) GetPaginatedPackages(endpoint string) (*PaginatedPackages, error) {
	rv := &PaginatedPackages{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Add("User-Agent", UserAgent)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = decodeResponse(resp, &rv.Packages)
	if err != nil {
		return nil, err
	}
	err = ExtractPaginationHeaders(&resp.Header, &rv.Paginated)
	if err != nil {
		return nil, err
	}
	header := resp.Header.Get("Link")
	links := linkheader.Parse(header)
	for _, link := range links {
		if link.Rel == "next" {
			if rv.MaxPerPage > 0 {
				u, err := url.Parse(link.URL)
				if err != nil {
					return nil, err
				}
				q := u.Query()
				q.Set("per_page", resp.Header.Get("Max-Per-Page"))
				u.RawQuery = q.Encode()
				link.URL = u.String()

			}
			rv.Next = func() (*PaginatedPackages, error) {
				return c.GetPaginatedPackages(link.URL)
			}
		}
	}
	return rv, nil
}

func ExtractPaginationHeaders(h *http.Header, p *Paginated) error {
	header := h.Get("Total")
	total, err := strconv.Atoi(header)
	if err != nil {
		return err
	}
	p.Total = total

	header = h.Get("Per-Page")
	perPage, err := strconv.Atoi(header)
	if err != nil {
		return err
	}
	p.PerPage = perPage

	header = h.Get("Max-Per-Page")
	maxPerPage, err := strconv.Atoi(header)
	if err != nil {
		return err
	}
	p.MaxPerPage = maxPerPage

	return nil
}

func (c *Client) PaginatedAll(repo string) (*PaginatedPackages, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/packages.json", ServiceURL, repo)
	return c.GetPaginatedPackages(endpoint)
}
