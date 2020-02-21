package satellite

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrganizationRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"hosts_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			//"locations"
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orgID := d.Get("id").(int)

	org, _, err := client.Organizations.GetOrganizationByID(context.Background(), orgID)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(int(*org.ID)))

	d.Set("ancestry", org.Ancestry)
	d.Set("hosts_count", org.HostsCount)
	d.Set("name", org.Name)
	d.Set("parent_id", org.ParentID)
	d.Set("parent_name", org.ParentName)
	d.Set("title", org.Title)

	return nil
}
