package bizflycloud

import (
	"context"
	"github.com/bizflycloud/gobizfly"
	"log"
	"time"
)

type Config struct {
	AuthMethod          string
	Email               string
	Password            string
	RegionName          string
	AppCredentialID     string
	AppCredentialSecret string
	APIEndpoint         string
	TerraformVersion    string
}

type CombinedConfig struct {
	client *gobizfly.Client
}

func (c *CombinedConfig) gobizflyClient() *gobizfly.Client { return c.client }

func (c *Config) Client() (*CombinedConfig, error) {
	client, err := gobizfly.NewClient(gobizfly.WithTenantName(c.Email), gobizfly.WithRegionName(c.RegionName), gobizfly.WithAPIUrl(c.APIEndpoint))

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	tok, err := client.Token.Create(ctx, &gobizfly.TokenCreateRequest{
		AuthMethod: c.AuthMethod, Username: c.Email, Password: c.Password, AppCredID: c.AppCredentialID, AppCredSecret: c.AppCredentialSecret})
	if err != nil {
		return nil, err
	}

	client.SetKeystoneToken(tok.KeystoneToken)

	return &CombinedConfig{
		client: client,
	}, nil
}
