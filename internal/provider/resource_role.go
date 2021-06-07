package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage a role in Red Hat Satellite.",

		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "A name for the role.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Description: "A description of the role.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"location_ids": {
				Description: "A list of IDs of locations to associate with the role.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"organization_ids": {
				Description: "A list of IDs of organizations to associate with the role.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"builtin": {
				Description: "A boolean that indicates if the role is a default/builtin role.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"cloned_from_id": {
				Description: "If this role was cloned using the API, the ID number of the source role.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"filters": {
				Description: "A list of filter IDs associated with the role.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"locations": {
				Description: "A list of objects containing the locations the role applies to.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Description: "The description of a location the filter applies to.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "The ID of a location the filter applies to.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "The name of a location the filter applies to.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"title": {
							Description: "The title of a location the filter applies to.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"organizations": {
				Description: "A list of objects containing the organizations the role applies to.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Description: "The description of an organization the filter applies to.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "The ID of an organization the filter applies to.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "The name of an organization the filter applies to.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"title": {
							Description: "The title of an organization the filter applies to.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"origin": {
				Description: "TODO",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	roleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	role, resp, err := client.Roles.Get(context.Background(), roleID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
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

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

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
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(*role.ID))

	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	roleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}

	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	roleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Roles.Delete(context.Background(), roleID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
