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
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func datasourceBizflyCloudDatabaseDatastore() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizflyCloudDatabaseDatastoreRead,
		Schema: dataCloudDatabaseDatastoreSchema(),
	}
}

func dataSourceBizflyCloudDatabaseDatastoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	log.Println("[DEBUG] Reading list datastore")
	engines, err := client.CloudDatabase.Engines().List(context.Background())
	if err != nil {
		return fmt.Errorf("error describing datastore: %w", err)
	}

	dsType := d.Get("type").(string)
	dsName := d.Get("name").(string)

	for _, engine := range engines {
		if engine.Name != dsType {
			continue
		}

		for _, version := range engine.Versions {
			if version.Name != dsName {
				continue
			}

			d.SetId(engine.ID)
			_ = d.Set("version_id", version.ID)
			break
		}
		break
	}

	return nil
}
