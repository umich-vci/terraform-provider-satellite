package satellite

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func dataSourcePermissions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePermissionsRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"permissions": &schema.Schema{
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

	searchBody := new(gosatellite.PermissionsSearch)
	id := "search="
	idList := []string{}

	if n, ok := d.GetOk("name"); ok {
		name := n.(string)
		searchBody.Name = &name
		idList = append(idList, "name="+name)
	}

	if r, ok := d.GetOk("resource_type"); ok {
		resourceType := r.(string)
		searchBody.ResourceType = &resourceType
		idList = append(idList, "resource_type="+resourceType)
	}

	id = id + strings.Join(idList, "&")

	perms, _, err := client.Permissions.ListPermissions(context.Background(), *searchBody)
	if err != nil {
		return err
	}

	d.SetId(id)

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
