package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func dataSourceLocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationRead,
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"parent_name": {
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

func dataSourceLocationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	searchString := d.Get("search").(string)

	locSearch := new(gosatellite.LocationsListOptions)
	locSearch.Search = searchString

	location, _, err := client.Locations.List(context.Background(), locSearch)
	if err != nil {
		return err
	}

	locationList := *location.Results

	if len(locationList) == 0 {
		return fmt.Errorf("No locations found for search string %s", searchString)
	}

	if len(locationList) > 1 {
		return fmt.Errorf("%d locations found for search string %s", len(locationList), searchString)
	}

	d.SetId(strconv.Itoa(*locationList[0].ID))
	d.Set("created_at", *locationList[0].CreatedAt)
	d.Set("description", *locationList[0].Description)
	d.Set("name", *locationList[0].Name)
	d.Set("parent_id", *locationList[0].ParentID)
	d.Set("parent_name", *locationList[0].ParentName)
	d.Set("title", *locationList[0].Title)
	d.Set("updated_at", *locationList[0].UpdatedAt)

	return nil
}
