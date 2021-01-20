package satellite

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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": &schema.Schema{
				Type:     schema.TypeInt,
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

	name := d.Get("name").(string)
	searchString := fmt.Sprintf("name=\"%s\"", name)

	locSearch := gosatellite.LocationsSearch{
		Search: &searchString,
	}

	location, _, err := client.Locations.ListLocations(context.Background(), locSearch)
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
	d.Set("description", *locationList[0].Description)
	d.Set("parent_id", *locationList[0].ParentID)

	return nil
}
