package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type authMetadata struct {
	realm   string
	service string
	scope   string
}

type DockerRegistryClient struct {
	apiBaseUrl string
	token      string
	httpClient http.Client
}

func NewDockerRegistryClient() *DockerRegistryClient {
	d := &DockerRegistryClient{
		apiBaseUrl: "https://registry.hub.docker.com",
		httpClient: http.Client{
			Timeout: 30 * time.Second,
		},
	}
	return d
}

func (d *DockerRegistryClient) PullImage(name, reference string) ([]byte, error) {
	pullImageUrl := fmt.Sprintf("%s/v2/library/%s/manifests/%s", d.apiBaseUrl, name, reference)

	req, err := http.NewRequest(http.MethodGet, pullImageUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.token))

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		err = d.requestToken(d.getAuthMetadataFromHeader(resp.Header.Get("Www-Authenticate")))
		if err != nil {
			return nil, err
		}
		return d.PullImage(name, reference)
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (d *DockerRegistryClient) PullLayer(name, digest string) ([]byte, error) {
	pullLayerUrl := fmt.Sprintf("%s/v2/%s/blobs/%s", d.apiBaseUrl, name, digest)

	req, err := http.NewRequest(http.MethodGet, pullLayerUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.token))

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

func (d *DockerRegistryClient) requestToken(data authMetadata) error {
	authUrl := fmt.Sprintf("https://auth.docker.io/token?realm=%s&service=%s&scope=%s", data.realm, data.service, data.scope)
	resp, err := d.httpClient.Get(authUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	val := make(map[string]interface{})
	err = json.Unmarshal(bytes, &val)
	if err != nil {
		return err
	}

	d.token = val["token"].(string)
	return nil
}

func (d *DockerRegistryClient) getAuthMetadataFromHeader(authHeader string) authMetadata {
	var authData string
	fmt.Sscanf(authHeader, "Bearer realm=%s", &authData)

	parts := strings.Split(authData, ",")

	for i, p := range parts {
		if strings.Contains(p, "=") {
			parts[i] = strings.Split(p, "=")[1]
		}
		parts[i] = parts[i][1 : len(parts[i])-1]
	}

	return authMetadata{
		realm:   parts[0],
		service: parts[1],
		scope:   parts[2],
	}
}
