package satellite

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceOrganizationCreate,
		Read:   resourceOrganizationRead,
		Update: resourceOrganizationUpdate,
		Delete: resourceOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"hosts_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	org, resp, err := client.Organizations.Get(context.Background(), orgID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("description", org.Description)
	d.Set("hosts_count", org.HostsCount)
	d.Set("label", org.Label)
	d.Set("name", org.Name)
	d.Set("title", org.Title)

	return nil
}

func resourceOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	createBody := new(gosatellite.OrganizationCreate)
	createBody.Organization.Name = d.Get("name").(string)

	org, _, err := client.Organizations.Create(context.Background(), *createBody)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(*org.ID))

	return resourceOrganizationRead(d, meta)
}

func resourceOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	updateBody := new(gosatellite.OrganizationUpdate)
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateBody.Organization.Description = &description
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateBody.Organization.Name = &name
	}

	_, _, err = client.Organizations.Update(context.Background(), orgID, *updateBody)
	if err != nil {
		return err
	}

	return resourceOrganizationRead(d, meta)
}

func resourceOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Organizations.Delete(context.Background(), orgID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
