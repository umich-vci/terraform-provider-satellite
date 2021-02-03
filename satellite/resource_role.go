package satellite

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"organization_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"builtin": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cloned_from_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"filters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"locations": {
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
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"organizations": {
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
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"origin": {
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

	role, resp, err := client.Roles.Get(context.Background(), roleID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	var locationIDs []int
	var organizationIDs []int

	locationsList := []map[string]interface{}{}
	for _, x := range *role.Locations {
		locationIDs = append(locationIDs, *x.ID)
		location := make(map[string]interface{})
		location["description"] = x.Description
		location["id"] = x.ID
		location["name"] = x.Name
		location["title"] = x.Title
		locationsList = append(locationsList, location)
	}

	organizationsList := []map[string]interface{}{}
	for _, x := range *role.Organizations {
		organizationIDs = append(organizationIDs, *x.ID)
		organization := make(map[string]interface{})
		organization["description"] = x.Description
		organization["id"] = x.ID
		organization["name"] = x.Name
		organization["title"] = x.Title
		organizationsList = append(organizationsList, organization)
	}

	filtersList := []int{}
	for _, x := range *role.Filters {
		filtersList = append(filtersList, *x.ID)
	}

	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("location_ids", locationIDs)
	d.Set("organization_ids", organizationIDs)
	d.Set("builtin", role.Builtin)
	d.Set("cloned_from_id", role.ClonedFromID)
	d.Set("filters", filtersList)
	d.Set("locations", locationsList)
	d.Set("organizations", organizationsList)
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

	if desc, ok := d.GetOk("description"); ok {
		description := desc.(string)
		createBody.Role.Description = &description
	}

	if loc, ok := d.GetOk("location_ids"); ok {
		//rawLocationIDs := d.Get("location_ids").(*schema.Set).List()
		rawLocationIDs := loc.(*schema.Set).List()
		locationIDs := []int{}
		for x := range rawLocationIDs {
			locationIDs = append(locationIDs, rawLocationIDs[x].(int))
		}
		createBody.Role.LocationIDs = &locationIDs
	}

	if org, ok := d.GetOk("organization_ids"); ok {
		//rawOrganizationIDs := d.Get("organization_ids").(*schema.Set).List()
		rawOrganizationIDs := org.(*schema.Set).List()
		organizationIDs := []int{}
		for x := range rawOrganizationIDs {
			organizationIDs = append(organizationIDs, rawOrganizationIDs[x].(int))
		}
		createBody.Role.OrganizationIDs = &organizationIDs
	}

	role, _, err := client.Roles.Create(context.Background(), *createBody)
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

	_, _, err = client.Roles.Update(context.Background(), roleID, *updateBody)
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

	_, err = client.Roles.Delete(context.Background(), roleID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
