package satellite

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func dataSourceOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrganizationRead,
		Schema: map[string]*schema.Schema{
			"search": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": {
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

func dataSourceOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	searchString := d.Get("search").(string)

	opt := new(gosatellite.OrganizationsListOptions)
	opt.Search = searchString

	org, _, err := client.Organizations.List(context.Background(), opt)
	if err != nil {
		return err
	}

	orgList := *org.Results

	if len(orgList) == 0 {
		return fmt.Errorf("No organizations found for search string %s", searchString)
	}

	if len(orgList) > 1 {
		return fmt.Errorf("%d organizations found for search string %s", len(orgList), searchString)
	}

	d.SetId(strconv.Itoa(int(*orgList[0].ID)))

	d.Set("created_at", orgList[0].CreatedAt)
	d.Set("description", orgList[0].Description)
	d.Set("label", orgList[0].Label)
	d.Set("name", orgList[0].Name)
	d.Set("title", orgList[0].Title)
	d.Set("updated_at", orgList[0].UpdatedAt)

	return nil
}
