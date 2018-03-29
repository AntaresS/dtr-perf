package sharedutils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func MakeRegistryAuth(username, password, refreshToken string) string {
	authCfg := types.AuthConfig{}
	if refreshToken != "" {
		authCfg.IdentityToken = refreshToken
	} else {
		authCfg.Username = username
		authCfg.Password = password
	}
	buf, _ := json.Marshal(authCfg)
	return base64.URLEncoding.EncodeToString(buf)
}

func MakeDockerClient(host, jwt string, httpClient *http.Client, hubUsername, hubPassword, refreshToken string) (*client.Client, error) {
	if !strings.HasPrefix(host, "unix://") {
		hostPort := strings.Split(host, ":")
		if len(hostPort) == 1 {
			host = fmt.Sprintf("tcp://%s:443", host)
		} else {
			host = fmt.Sprintf("tcp://%s:%s", hostPort[0], hostPort[1])
		}
	}
	headers := map[string]string{}
	if jwt != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", jwt)
	}
	if hubUsername != "" && hubPassword != "" {
		headers["X-Registry-Auth"] = MakeRegistryAuth(hubUsername, hubPassword, refreshToken)
	}
	return client.NewClient(host, "v1.22", httpClient, headers)
}

func ReadConfigFile(file string) ([]byte, error) {
	configBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return configBytes, err
}
