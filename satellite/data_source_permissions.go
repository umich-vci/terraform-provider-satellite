package satellite

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func dataSourcePermissions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePermissionsRead,
		Schema: map[string]*schema.Schema{
			"search": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePermissionsRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	searchOpt := new(gosatellite.PermissionsListOptions)

	if n, ok := d.GetOk("search"); ok {
		searchOpt.Search = n.(string)
	}

	perms, _, err := client.Permissions.List(context.Background(), *searchOpt)
	if err != nil {
		return err
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
