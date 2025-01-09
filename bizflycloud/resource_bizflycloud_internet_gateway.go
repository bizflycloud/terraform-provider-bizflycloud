package bizflycloud

import (
	"context"
	"fmt"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceInternetGateway() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: resourceInternetGatewaySchema(),
		Read:   resourceInternetGatewayRead,
		Create: resourceInternetGatewayCreate,
		Update: resourceInternetGatewayUpdate,
		Delete: resourceInternetGatewayDelete,
	}
}

func createIGWBuilder(d *schema.ResourceData) gobizfly.CreateInternetGatewayPayload {
	opts := gobizfly.CreateInternetGatewayPayload{}
	if v, ok := d.GetOk("name"); ok {
		opts.Name = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		opts.Description = &desc
	}
	if v, ok := d.GetOk("vpc_network_id"); ok {
		vpcNetworkID := v.(string)
		if vpcNetworkID != "" {
			netIDs := []string{vpcNetworkID}
			opts.NetworkIDs = &netIDs
		}
	}
	return opts
}

func resourceInternetGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	otps := createIGWBuilder(d)
	createdIGW, err := client.CloudServer.InternetGateways().Create(context.Background(), otps)
	if err != nil {
		return fmt.Errorf("Error when creating internet gateway: %v", err)
	}
	d.SetId(createdIGW.ID)
	return resourceInternetGatewayRead(d, meta)
}

func resourceInternetGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	igw, err := client.CloudServer.InternetGateways().Get(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when retrieving internet gateway: %v", err)
	}
	_ = d.Set("name", igw.Name)
	_ = d.Set("description", igw.Description)
	_ = d.Set("status", igw.Status)
	_ = d.Set("project_id", igw.ProjectID)
	_ = d.Set("availability_zones", igw.AvailabilityZones)
	_ = d.Set("tags", igw.Tags)
	_ = d.Set("created_at", igw.CreatedAt)
	_ = d.Set("updated_at", igw.UpdatedAt)
	if len(igw.InterfacesInfo) > 0 {
		interfaceInfo := igw.InterfacesInfo[0]
		vpcNetworkID := interfaceInfo.NetworkInfo.ID
		vpcNetworkName := interfaceInfo.NetworkInfo.Name
		_ = d.Set("vpc_network_id", vpcNetworkID)
		_ = d.Set("vpc_network_name", vpcNetworkName)
	} else {
		_ = d.Set("vpc_network_id", "")
		_ = d.Set("vpc_network_name", "")
	}
	return nil
}

func resourceInternetGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	if d.HasChanges("name", "description", "vpc_network_id") {
		opts := gobizfly.UpdateInternetGatewayPayload{}
		if v, ok := d.GetOk("name"); ok {
			opts.Name = v.(string)
		}
		if v, ok := d.GetOk("description"); ok {
			desc := v.(string)
			opts.Description = desc
		}
		if v, ok := d.GetOk("vpc_network_id"); ok {
			vpcNetworkID := v.(string)
			if vpcNetworkID != "" {
				netIDs := []string{vpcNetworkID}
				opts.NetworkIDs = netIDs
			}
		}
		_, err := client.CloudServer.InternetGateways().Update(context.Background(), d.Id(), opts)
		if err != nil {
			return fmt.Errorf("Error when updating internet gateway: %v", err)
		}
	}
	return resourceInternetGatewayRead(d, meta)
}

func resourceInternetGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.CloudServer.InternetGateways().Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error when deleting internet gateway: %v", err)
	}
	return nil
}
