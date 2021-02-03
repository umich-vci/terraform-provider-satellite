package satellite

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupCreate,
		Read:   resourceUserGroupRead,
		Update: resourceUserGroupUpdate,
		Delete: resourceUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"admin": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// "user_ids": {
			// 	Type:     schema.TypeSet,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeInt,
			// 	},
			// },
			// "usergroup_ids": {
			// 	Type:     schema.TypeSet,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeInt,
			// 	},
			// },
			"role_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	ugID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	ug, resp, err := client.UserGroups.Get(context.Background(), ugID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	roleIDs := []int{}
	roleList := []map[string]interface{}{}
	for _, x := range *ug.Roles {
		roleIDs = append(roleIDs, *x.ID)
		role := make(map[string]interface{})
		role["description"] = x.Description
		role["id"] = x.ID
		role["name"] = x.Name
		role["origin"] = x.Origin
		roleList = append(roleList, role)
	}

	d.Set("name", ug.Name)
	d.Set("admin", ug.Admin)
	d.Set("role_ids", roleIDs)
	d.Set("roles", roleList)
	d.Set("created_at", ug.CreatedAt)
	d.Set("updated_at", ug.UpdatedAt)

	return nil
}

func resourceUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	createBody := new(gosatellite.UserGroupCreate)
	createBody.UserGroup.Name = &name

	if adm, ok := d.GetOk("admin"); ok {
		admin := adm.(bool)
		createBody.UserGroup.Admin = &admin
	}

	if ri, ok := d.GetOk("role_ids"); ok {
		rawRoleIDs := ri.(*schema.Set).List()
		roleIDs := []int{}
		for x := range rawRoleIDs {
			roleIDs = append(roleIDs, rawRoleIDs[x].(int))
		}
		createBody.UserGroup.RoleIDs = &roleIDs
	}

	ug, _, err := client.UserGroups.Create(context.Background(), *createBody)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(*ug.ID))

	return resourceUserGroupRead(d, meta)
}

func resourceUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	ugID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	updateBody := new(gosatellite.UserGroupUpdate)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateBody.UserGroup.Name = &name
	}
	if d.HasChange("admin") {
		admin := d.Get("admin").(bool)
		updateBody.UserGroup.Admin = &admin
	}
	if d.HasChange("role_ids") {
		rawRoleIDs := d.Get("role_ids").(*schema.Set).List()
		roleIDs := []int{}
		for x := range rawRoleIDs {
			roleIDs = append(roleIDs, rawRoleIDs[x].(int))
		}
		updateBody.UserGroup.RoleIDs = &roleIDs
	}

	_, _, err = client.UserGroups.Update(context.Background(), ugID, *updateBody)
	if err != nil {
		return err
	}
	return resourceUserGroupRead(d, meta)
}

func resourceUserGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	ugID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, _, err = client.UserGroups.Delete(context.Background(), ugID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
