package satellite

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"location_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"organization_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"builtin": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cloned_from_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"filters": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"locations": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"organizations": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"origin": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRoleRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	roleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	role, _, err := client.Roles.GetRoleByID(context.Background(), roleID)
	if err != nil {
		return err
	}

	var locationIDs []int
	var organizationIDs []int

	for _, x := range *role.Locations {
		locationIDs = append(locationIDs, *x.ID)
	}

	for _, x := range *role.Organizations {
		organizationIDs = append(organizationIDs, *x.ID)
	}

	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("location_ids", locationIDs)
	d.Set("organization_ids", organizationIDs)
	d.Set("builtin", role.Builtin)
	d.Set("cloned_from_id", role.ClonedFromID)
	d.Set("filters", role.Filters)
	d.Set("locations", role.Locations)
	d.Set("organizations", role.Organizations)
	d.Set("origin", role.Origin)

	return nil
}

func resourceRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	createBody := new(gosatellite.RoleCreate)
	createBody.Role.Name = &name

	role, _, err := client.Roles.CreateRole(context.Background(), *createBody)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(*role.ID))

	return resourceRoleRead(d, meta)
}

func resourceRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	roleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	updateBody := new(gosatellite.RoleUpdate)
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateBody.Role.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateBody.Role.Description = &description
	}
	if d.HasChange("location_ids") {
		rawLocationIDs := d.Get("location_ids").(*schema.Set).List()
		locationIDs := []int{}
		for x := range rawLocationIDs {
			locationIDs = append(locationIDs, rawLocationIDs[x].(int))
		}
		updateBody.Role.LocationIDs = &locationIDs
	}
	if d.HasChange("organization_ids") {
		rawOrganizationIDs := d.Get("organization_ids").(*schema.Set).List()
		organizationIDs := []int{}
		for x := range rawOrganizationIDs {
			organizationIDs = append(organizationIDs, rawOrganizationIDs[x].(int))
		}
		updateBody.Role.OrganizationIDs = &organizationIDs
	}

	_, _, err = client.Roles.UpdateRole(context.Background(), roleID, *updateBody)
	if err != nil {
		return err
	}

	return resourceRoleRead(d, meta)
}

func resourceRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	roleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Roles.DeleteRole(context.Background(), roleID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
