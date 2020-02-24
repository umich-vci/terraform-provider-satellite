package satellite

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceOrganizationCreate,
		Read:   resourceOrganizationRead,
		Update: resourceOrganizationUpdate,
		Delete: resourceOrganizationDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"hosts_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"title": &schema.Schema{
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

	org, resp, err := client.Organizations.GetOrganizationByID(context.Background(), orgID)
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
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

	org, _, err := client.Organizations.CreateOrganization(context.Background(), *createBody)
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

	_, _, err = client.Organizations.UpdateOrganization(context.Background(), orgID, *updateBody)
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

	_, err = client.Organizations.DeleteOrganization(context.Background(), orgID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
