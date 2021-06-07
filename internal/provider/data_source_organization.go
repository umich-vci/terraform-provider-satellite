package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func dataSourceOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to access information about a Red Hat Satellite organization.",

		ReadContext: dataSourceOrganizationRead,

		Schema: map[string]*schema.Schema{
			"search": {
				Description:  "A search filter for the Location search. The search must only return 1 Organization.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"created_at": {
				Description: "Timestamp of when the organization was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The description of the organization.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"label": {
				Description: "The label of the organization.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the organization.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"title": {
				Description: "The title of the organization.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "Timestamp of when the organization was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	searchString := d.Get("search").(string)

	opt := new(gosatellite.OrganizationsListOptions)
	opt.Search = searchString

	org, _, err := client.Organizations.List(context.Background(), opt)
	if err != nil {
		return diag.FromErr(err)
	}

	orgList := *org.Results

	if len(orgList) == 0 {
		return diag.Errorf("No organizations found for search string %s", searchString)
	}

	if len(orgList) > 1 {
		return diag.Errorf("%d organizations found for search string %s", len(orgList), searchString)
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
