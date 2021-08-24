package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func dataSourceLocation() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to access information about a Red Hat Satellite location.",

		ReadContext: dataSourceLocationRead,

		Schema: map[string]*schema.Schema{
			"search": {
				Description:  "A search filter for the Location search. The search must only return 1 Location.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"created_at": {
				Description: "Timestamp of when the location was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "A description of the location.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the location.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"parent_id": {
				Description: "The ID of the parent for this location.  If not set, the location is a top level location.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"parent_name": {
				Description: "The name of the parent for this location.  If not set, the location is a top level location.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"title": {
				Description: "The title of the location.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "Timestamp of when the location was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	searchString := d.Get("search").(string)

	locSearch := new(gosatellite.LocationsListOptions)
	locSearch.Search = searchString

	location, _, err := client.Locations.List(context.Background(), locSearch)
	if err != nil {
		return diag.FromErr(err)
	}

	locationList := *location.Results

	if len(locationList) == 0 {
		return diag.Errorf("No locations found for search string %s", searchString)
	}

	if len(locationList) > 1 {
		return diag.Errorf("%d locations found for search string %s", len(locationList), searchString)
	}

	d.SetId(strconv.Itoa(*locationList[0].ID))
	d.Set("created_at", locationList[0].CreatedAt)
	d.Set("description", locationList[0].Description)
	d.Set("name", locationList[0].Name)
	d.Set("parent_id", locationList[0].ParentID)
	d.Set("parent_name", locationList[0].ParentName)
	d.Set("title", locationList[0].Title)
	d.Set("updated_at", locationList[0].UpdatedAt)

	return nil
}
