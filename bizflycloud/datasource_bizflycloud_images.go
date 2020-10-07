package bizflycloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func datasourceBizFlyCloudImages() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBizFlyCloudImageRead,
		Schema: imageSchema(),
	}
}

func dataSourceBizFlyCloudImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	osImages, err := client.Server.ListOSImages(context.Background())
	if err != nil {
		return err
	}
	distribution, okDist := d.GetOk("distribution")
	version, okVer := d.GetOk("version")
	if okDist && okVer {
		for _, image := range osImages {
			if strings.ToLower(image.OSDistribution) != strings.ToLower(distribution.(string)) {
				continue
			}
			for _, v := range image.Version {
				if strings.ToLower(v.Name) != strings.ToLower(version.(string)) {
					continue
				}
				d.SetId(v.ID)
				break
			}
		}
	} else {
		return fmt.Errorf("Distribution and Version must be set")
	}
	return nil
}
