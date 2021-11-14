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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func imageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "Id of the image",
			Computed:    true,
		},
		"distribution": {
			Type:        schema.TypeString,
			Description: "OS Distribution",
			Required:    true,
		},
		"version": {
			Type:        schema.TypeString,
			Description: "OS Version",
			Required:    true,
		},
	}
}
