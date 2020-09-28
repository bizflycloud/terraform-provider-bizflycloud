// This file is part of gobizfly
//
// Copyright (C) 2020  BizFly Cloud
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>

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
	client, err := gobizfly.NewClient(gobizfly.WithTenantName(c.Email), gobizfly.WithRegionName(c.RegionName), gobizfly.WithAPIUrl(c.APIEndpoint)) // nolint

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
