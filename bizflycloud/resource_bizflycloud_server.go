package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudServer() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizFlyCloudServerCreate,
		Read:          resourceBizFlyCloudServerRead,
		Update:        resourceBizFlyCloudServerUpdate,
		Delete:        resourceBizFlyCloudServerDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"os_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"root_disk_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"root_disk_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBizFlyCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()

	// Build up creation options
	scr := &gobizfly.ServerCreateRequest{
		Name:       d.Get("name").(string),
		FlavorName: d.Get("flavor_name").(string),
		SSHKey:     d.Get("ssh_key").(string),
		Type:       d.Get("category").(string),
		OS: &gobizfly.ServerOS{
			Type: d.Get("os_type").(string),
			ID:   d.Get("os_id").(string),
		},
		AvailabilityZone: d.Get("availability_zone").(string),
		Password:         d.Get("password").(bool),
		RootDisk: &gobizfly.ServerDisk{
			Type: d.Get("root_disk_type").(string),
			Size: d.Get("root_disk_size").(int),
		},
	}
	log.Printf("[DEBUG] Create Cloud Server configuration: %#v", scr)

	tasks, err := client.Server.Create(context.Background(), scr)
	if err != nil {
		return fmt.Errorf("Error creating server: %s", err)
	}
	// Set ID of server with task ID, we need to change to the real ID after server is created
	d.SetId(tasks.Task[0])
	log.Printf("[INFO] Server is creating with task ID: %s", d.Id())
	// wait for cloud server to become active
	_, err = waitForServerCreate(d, meta)
	if err != nil {
		return fmt.Errorf("Error creating cloud server with task id (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceBizFlyCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	server, err := client.Server.Get(context.Background(), d.Id())
	if err != nil {
		if errors.Is(err, gobizfly.ErrNotFound) {
			log.Printf("[WARN] BizFly Cloud Server (%s) is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving server: %v", err)
	}
	d.Set("name", server.Name)
	d.Set("key_name", server.KeyName)
	d.Set("status", server.Status)
	d.Set("flavor_name", formatFlavor(server.Flavor.Name))
	d.Set("category", server.Category)
	d.Set("user_id", server.UserID)
	d.Set("project_id", server.ProjectID)
	d.Set("availability_zone", server.AvailabilityZone)
	d.Set("created_at", server.CreatedAt)
	d.Set("updated_at", server.UpdatedAt)
	return nil
}

func resourceBizFlyCloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	id := d.Id()
	if d.HasChange("flavor_name") {
		// Resize server to new flavor
		task, err := client.Server.Resize(context.Background(), d.Id(), d.Get("flavor_name").(string))
		if err != nil {
			return fmt.Errorf("Error when resize server [%s]: %v", id, err)
		}
		// wait for server is active again
		_, err = waitForServerUpdate(d, meta, task.TaskID)
		if err != nil {
			return fmt.Errorf("Error updating cloud server with task id (%s): %s", d.Id(), err)
		}
	}
	return resourceBizFlyCloudServerRead(d, meta)
}

func resourceBizFlyCloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	// destroy the cloud server
	err := client.Server.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error delete cloud server %v", err)
	}
	// TODO check server is deleted
	return nil
}

func waitForServerCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	log.Printf("[INFO] Waiting for server with task id (%s) to be created", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    newServerStateRefreshfunc(d, "status", meta),
		Timeout:    600 * time.Second,
		Delay:      20 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func waitForServerUpdate(d *schema.ResourceData, meta interface{}, taskID string) (interface{}, error) {
	log.Printf("[INFO] Waiting for server with task id (%s) to be created", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"HARD_REBOOT", "MIGRATING", "REBUILD", "RESIZE"},
		Target:     []string{"ACTIVE"},
		Refresh:    updateServerStateRefreshfunc(d, "status", meta, taskID),
		Timeout:    600 * time.Second,
		Delay:      20 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func newServerStateRefreshfunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()
	return func() (interface{}, string, error) {
		// Get task result from cloud server API
		resp, err := client.Server.GetTask(context.Background(), d.Id())
		if err != nil {
			return nil, "", err
		}
		// if the task is not ready, we need to wait for a moment
		if !resp.Ready {
			log.Println("[DEBUG] Cloud Server is not ready")
			return nil, "", nil
		}
		// server is ready now, set ID for resourceData
		d.SetId(resp.Result.Server.ID)
		err = resourceBizFlyCloudServerRead(d, meta)
		if err != nil {
			return nil, "", err
		}
		if attr, ok := d.GetOkExists(attribute); ok {
			server, err := client.Server.Get(context.Background(), resp.Result.Server.ID)
			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving cloud server: %v", err)
			}
			switch attr.(type) {
			case bool:
				return &server, strconv.FormatBool(attr.(bool)), nil
			default:
				return &server, attr.(string), nil
			}
		}
		return nil, "", nil
	}
}

func updateServerStateRefreshfunc(d *schema.ResourceData, attribute string, meta interface{}, taskID string) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).gobizflyClient()
	return func() (interface{}, string, error) {
		// Get task result from cloud server API
		resp, err := client.Server.GetTask(context.Background(), taskID)
		if err != nil {
			return nil, "", err
		}
		// if the task is not ready, we need to wait for a moment
		if !resp.Ready {
			log.Println("[DEBUG] Cloud Server is not ready")
			return nil, "", nil
		}
		err = resourceBizFlyCloudServerRead(d, meta)
		if err != nil {
			return nil, "", err
		}
		if attr, ok := d.GetOkExists(attribute); ok {
			server, err := client.Server.Get(context.Background(), d.Id())
			if err != nil {
				return nil, "", fmt.Errorf("Error retrieving cloud server: %v", err)
			}
			switch attr.(type) {
			case bool:
				return &server, strconv.FormatBool(attr.(bool)), nil
			default:
				return &server, attr.(string), nil
			}
		}
		return nil, "", nil
	}
}
func formatFlavor(s string) string {
	// This function will be removed in the near future when the API format for us
	if strings.Contains(s, ".") {
		return strings.Split(s, ".")[1]
	}
	return strings.Join(strings.Split(s, "_")[:2], "_")
}
