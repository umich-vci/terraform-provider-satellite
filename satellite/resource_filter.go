package satellite

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func resourceFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceFilterCreate,
		Read:   resourceFilterRead,
		Update: resourceFilterUpdate,
		Delete: resourceFilterDelete,

		Schema: map[string]*schema.Schema{
			"role_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
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
			"override": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"permission_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"search": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
			"resource_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"unlimited": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFilterRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	filterID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	filter, resp, err := client.Filters.GetFilterByID(context.Background(), filterID)
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	var locationIDs []int
	var organizationIDs []int

	for _, x := range *filter.Locations {
		locationIDs = append(locationIDs, *x.ID)
	}

	for _, x := range *filter.Organizations {
		organizationIDs = append(organizationIDs, *x.ID)
	}

	d.Set("role_id", filter.Role.ID)
	d.Set("location_ids", locationIDs)
	d.Set("organization_ids", organizationIDs)
	d.Set("override", filter.Override)
	d.Set("permission_ids", filter.Permissions)
	d.Set("search", filter.Search)
	d.Set("created_at", filter.CreatedAt)
	d.Set("locations", filter.Locations)
	d.Set("organizations", filter.Organizations)
	d.Set("resource_type", filter.ResourceType)
	d.Set("role", filter.Role)
	d.Set("unlimited", filter.Unlimited)
	d.Set("updated_at", filter.UpdatedAt)

	return nil
}

func resourceFilterCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	roleID := d.Get("role_id").(int)

	createBody := new(gosatellite.FilterCreate)
	createBody.Filter.RoleID = &roleID

	if loc, ok := d.GetOk("location_ids"); ok {
		rawLocationIDs := loc.(*schema.Set).List()
		locationIDs := []int{}
		for x := range rawLocationIDs {
			locationIDs = append(locationIDs, rawLocationIDs[x].(int))
		}
		createBody.Filter.LocationIDs = &locationIDs
	}

	if org, ok := d.GetOk("organization_ids"); ok {
		rawOrganizationIDs := org.(*schema.Set).List()
		organizationIDs := []int{}
		for x := range rawOrganizationIDs {
			organizationIDs = append(organizationIDs, rawOrganizationIDs[x].(int))
		}
		createBody.Filter.OrganizationIDs = &organizationIDs
	}

	if ovr, ok := d.GetOk("override"); ok {
		override := ovr.(bool)
		createBody.Filter.Override = &(override)
	}

	if perm, ok := d.GetOk("permission_ids"); ok {
		rawPermissionIDs := perm.(*schema.Set).List()
		permissionIDs := []int{}
		for x := range rawPermissionIDs {
			permissionIDs = append(permissionIDs, rawPermissionIDs[x].(int))
		}
		createBody.Filter.PermissionIDs = &permissionIDs
	}

	if srch, ok := d.GetOk("search"); ok {
		search := srch.(string)
		createBody.Filter.Search = &search
	}

	filter, _, err := client.Filters.CreateFilter(context.Background(), *createBody)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(*filter.ID))

	return resourceFilterRead(d, meta)
}

func resourceFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	filterID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	updateBody := new(gosatellite.FilterUpdate)

	if d.HasChange("role_id") {
		roleID := d.Get("role_id").(int)
		updateBody.Filter.RoleID = &roleID
	}

	if d.HasChange("location_ids") {
		if loc, ok := d.GetOk("location_ids"); ok {
			rawLocationIDs := loc.(*schema.Set).List()
			locationIDs := []int{}
			for x := range rawLocationIDs {
				locationIDs = append(locationIDs, rawLocationIDs[x].(int))
			}
			updateBody.Filter.LocationIDs = &locationIDs
		} else {
			updateBody.Filter.LocationIDs = new([]int)
		}
	}

	if d.HasChange("organization_ids") {
		if org, ok := d.GetOk("organization_ids"); ok {
			rawOrganizationIDs := org.(*schema.Set).List()
			organizationIDs := []int{}
			for x := range rawOrganizationIDs {
				organizationIDs = append(organizationIDs, rawOrganizationIDs[x].(int))
			}
			updateBody.Filter.OrganizationIDs = &organizationIDs
		} else {
			updateBody.Filter.OrganizationIDs = new([]int)
		}
	}

	if d.HasChange("override") {
		if ovr, ok := d.GetOk("override"); ok {
			override := ovr.(bool)
			updateBody.Filter.Override = &(override)
		}
	}

	if d.HasChange("permission_ids") {
		if perm, ok := d.GetOk("permission_ids"); ok {
			rawPermissionIDs := perm.(*schema.Set).List()
			permissionIDs := []int{}
			for x := range rawPermissionIDs {
				permissionIDs = append(permissionIDs, rawPermissionIDs[x].(int))
			}
			updateBody.Filter.PermissionIDs = &permissionIDs
		} else {
			updateBody.Filter.PermissionIDs = new([]int)
		}
	}

	if d.HasChange("search") {
		if srch, ok := d.GetOk("search"); ok {
			search := srch.(string)
			updateBody.Filter.Search = &search
		}
	}

	_, _, err = client.Filters.UpdateFilter(context.Background(), filterID, *updateBody)
	if err != nil {
		return err
	}

	return resourceFilterRead(d, meta)
}

func resourceFilterDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	filterID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Filters.DeleteFilter(context.Background(), filterID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
