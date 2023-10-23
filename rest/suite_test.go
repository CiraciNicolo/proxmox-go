package rest

import (
	"crypto/tls"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	restclient RESTClient
}

func (s *TestSuite) SetupSuite() {
	url := os.Getenv("PROXMOX_URL")
	user := os.Getenv("PROXMOX_USERNAME")
	password := os.Getenv("PROXMOX_PASSWORD")
	tokenid := os.Getenv("PROXMOX_TOKENID")
	secret := os.Getenv("PROXMOX_SECRET")
	if url == "" {
		s.T().Fatal("url must not be empty")
	}

	var loginOption ClientOption
	if user != "" && password != "" {
		loginOption = WithUserPassword(user, password)
	} else if tokenid != "" && secret != "" {
		loginOption = WithAPIToken(tokenid, secret)
	} else {
		s.T().Logf("username=%s, password=%s, tokenid=%s, secret=%s", user, password, tokenid, secret)
		s.T().Fatal("username&password or tokenid&secret pair must be provided")
	}

	base := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	restclient, err := NewRESTClient(url, loginOption, WithClient(&base))
	if err != nil {
		s.T().Logf("username=%s, password=%s, tokenid=%s, secret=%s", user, password, tokenid, secret)
		s.T().Fatalf("failed to create rest client: %v", err)
	}

	s.restclient = *restclient
}

func TestSuiteIntegration(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
