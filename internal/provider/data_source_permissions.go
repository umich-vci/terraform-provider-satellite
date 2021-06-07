package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func dataSourcePermissions() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to access information about available Red Hat Satellite permissions.",

		ReadContext: dataSourcePermissionsRead,
		Schema: map[string]*schema.Schema{
			"search": {
				Description: "A search filter for the permission search. If not specified all permissions are returned.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"permissions": {
				Description: "A list of permissions.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of the permission.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "The name of the permission.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"resource_type": {
							Description: "The resource type the permission applies to.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePermissionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	searchOpt := new(gosatellite.PermissionsListOptions)

	if n, ok := d.GetOk("search"); ok {
		searchOpt.Search = n.(string)
	}

	perms, _, err := client.Permissions.List(context.Background(), *searchOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("-")

	permList := make([]map[string]interface{}, 0, len(*perms.Results))

	for _, x := range *perms.Results {
		perm := map[string]interface{}{
			"id":            x.ID,
			"name":          x.Name,
			"resource_type": x.ResourceType,
		}
		permList = append(permList, perm)
	}

	d.Set("permissions", permList)

	return nil
}
