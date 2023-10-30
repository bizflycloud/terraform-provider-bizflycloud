// This file is part of terraform-provider-bizflycloud
//
// Copyright (C) 2021  Bizfly Cloud
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
	"log"
	"time"

	"github.com/bizflycloud/gobizfly"
)

// Config is define a client struct
type Config struct {
	AuthMethod          string
	Email               string
	Password            string
	RegionName          string
	AppCredentialID     string
	AppCredentialSecret string
	APIEndpoint         string
	TerraformVersion    string
	ProjectID           string
}

// CombinedConfig is ...
type CombinedConfig struct {
	client *gobizfly.Client
}

func (c *CombinedConfig) gobizflyClient() *gobizfly.Client { return c.client }

// Client is interface to connect plugin provider
func (c *Config) Client() (*CombinedConfig, error) {
	client, err := gobizfly.NewClient(gobizfly.WithProjectId(c.ProjectID),
		gobizfly.WithRegionName(c.RegionName),
		gobizfly.WithAPIUrl(c.APIEndpoint)) // nolint

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	log.Println("[INFO] Authenticating with Bizfly Cloud API")
	tok, err := client.Token.Create(ctx, &gobizfly.TokenCreateRequest{
		AuthMethod:    c.AuthMethod,
		Username:      c.Email,
		Password:      c.Password,
		AppCredID:     c.AppCredentialID,
		AppCredSecret: c.AppCredentialSecret,
		ProjectID:     c.ProjectID,
	},
	)
	if err != nil {
		return nil, err
	}

	client.SetKeystoneToken(tok)
	userInfo, err := client.Account.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	allowedToRegion := false
	log.Println("[DEBUG] Checking if you are allowed to access this region")
	log.Println("[DEBUG] Allowed Region: ", userInfo.UserRegions)
	for _, region := range userInfo.UserRegions {
		if c.RegionName == region.Code || c.RegionName == region.ShortName {
			allowedToRegion = true
		}
	}
	if !allowedToRegion {
		log.Fatal("You are not allowed to access this region")
	}

	return &CombinedConfig{
		client: client,
	}, nil
}
