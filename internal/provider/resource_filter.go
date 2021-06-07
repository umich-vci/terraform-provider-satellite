package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func resourceFilter() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage a permission filter for a role in Red Hat Satellite.",

		CreateContext: resourceFilterCreate,
		ReadContext:   resourceFilterRead,
		UpdateContext: resourceFilterUpdate,
		DeleteContext: resourceFilterDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"permission_names": {
				Description: "A list of permission names that should be enabled in the filter. The permission names must be valid for the role specified in `resource_type`.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"role_id": {
				Description: "The ID of the role that the filter should be created under.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"resource_type": {
				Description:  "The resource type of the filter.  Once this is set, it cannot be changed without recreating the filter.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(resourceTypeList, false),
			},
			"location_ids": {
				Description: "A list of IDs of locations to associate with the filter. Unless `override` is set to `true` this should generally contain the `location_ids` that the parent role is associated with. It may also need to be set to an empty list if you desire the permission to be `unlimited`.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"organization_ids": {
				Description: "A list of IDs of organizations to associate with the filter. Unless `override` is set to `true` this should generally contain the `organization_ids` that the parent role is associated with. It may also need to be set to an empty list if you desire the permission to be `unlimited`.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"override": {
				Description: "When set to true, you can specify `location_ids` and `organization_ids` to allow the role to access the `resource_type` in the specified locations and organizations.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"search": {
				Description: "If this is not set, then the filter will apply to all objects of the specified resource type. This means the value of `unlimited` will be true.  You can specify a search which can be used to limit the resources that the permission applies to. This will result in the value of `unlimited` being false. For more information see the [Red Hat documentation](https://access.redhat.com/documentation/en-us/red_hat_satellite/6.8/html/administering_red_hat_satellite/chap-Red_Hat_Satellite-Administering_Red_Hat_Satellite-Users_and_Roles#sect-Red_Hat_Satellite-Administering_Red_Hat_Satellite-Users_and_Roles-Granular_Permission_Filtering).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"created_at": {
				Description: "A timestamp of when the filter was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"locations": {
				Description: "A list of objects containing the locations the filter applies to.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"organizations": {
				Description: "A list of objects containing the organizations the filter applies to.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"permissions": {
				Description: "A list of objects containing the permissions enabled in the filter.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"permission_ids": {
				Description: "A list of permission IDs that match the list of permissions supplied in `permission_names`.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"role": {
				Description: "An object containing information about the role the filter is associated with.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"unlimited": {
				Description: "A boolean that indicates if a filter applies to all resources of the `resource_type` or just a subset of resources specified in `search`.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"updated_at": {
				Description: "A timestamp of when the filter was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceFilterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	filterID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	filter, resp, err := client.Filters.Get(context.Background(), filterID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	var resourceType string
	if filter.ResourceType == nil {
		resourceType = ""
	} else {
		resourceType = *filter.ResourceType
	}

	// set values we can directly set from struct
	d.Set("role_id", filter.Role.ID)
	d.Set("search", filter.Search)
	d.Set("created_at", filter.CreatedAt)
	d.Set("override", filter.Override)
	d.Set("resource_type", resourceType)
	d.Set("unlimited", filter.Unlimited)
	d.Set("updated_at", filter.UpdatedAt)

	// set location_ids
	var locationIDs []int
	for _, x := range *filter.Locations {
		locationIDs = append(locationIDs, *x.ID)
	}
	d.Set("location_ids", locationIDs)

	// set organization_ids
	var organizationIDs []int
	for _, x := range *filter.Organizations {
		organizationIDs = append(organizationIDs, *x.ID)
	}
	d.Set("organization_ids", organizationIDs)

	//set permission_ids and permission_names
	var permNames []string
	var permIDs []int
	for _, x := range *filter.Permissions {
		permNames = append(permNames, *x.Name)
		permIDs = append(permIDs, *x.ID)
	}
	d.Set("permission_ids", permIDs)
	d.Set("permission_names", permNames)

	// set permissions
	permissionsList := []map[string]string{}
	for _, x := range *filter.Permissions {
		permission := map[string]string{}
		if x.ID != nil {
			permission["id"] = strconv.Itoa(*x.ID)
		}
		if x.Name != nil {
			permission["name"] = *x.Name
		}
		if x.ResourceType != nil {
			permission["resource_type"] = *x.ResourceType
		}
		permissionsList = append(permissionsList, permission)
	}
	d.Set("permissions", permissionsList)

	// set locations
	locationsList := []map[string]string{}
	for _, x := range *filter.Locations {
		location := map[string]string{}
		if x.Description != nil {
			location["description"] = *x.Description
		}
		if x.ID != nil {
			location["id"] = strconv.Itoa(*x.ID)
		}
		if x.Name != nil {
			location["name"] = *x.Name
		}
		if x.Title != nil {
			location["title"] = *x.Title
		}
		locationsList = append(locationsList, location)
	}
	d.Set("locations", locationsList)

	// set organizations
	organizationList := []map[string]string{}
	for _, x := range *filter.Organizations {
		organization := map[string]string{}
		if x.Description != nil {
			organization["description"] = *x.Description
		}
		if x.ID != nil {
			organization["id"] = strconv.Itoa(*x.ID)
		}
		if x.Name != nil {
			organization["name"] = *x.Name
		}
		if x.Title != nil {
			organization["title"] = *x.Title
		}
		organizationList = append(organizationList, organization)
	}
	d.Set("organizations", organizationList)

	//set role
	role := map[string]string{}
	if filter.Role.Description != nil {
		role["description"] = *filter.Role.Description
	}
	if filter.Role.ID != nil {
		role["id"] = strconv.Itoa(*filter.Role.ID)
	}
	if filter.Role.Name != nil {
		role["name"] = *filter.Role.Name
	}
	if filter.Role.Origin != nil {
		role["origin"] = *filter.Role.Origin
	}
	d.Set("role", role)

	return nil
}

func resourceFilterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	roleID := d.Get("role_id").(int)
	resourceType := d.Get("resource_type").(string)
	createBody := new(gosatellite.FilterCreate)
	createBody.Filter.RoleID = &roleID

	// validate the permission names passed in
	permSearchOpts := new(gosatellite.PermissionsListOptions)

	// For the miscellaneous role_type which has a value of null,
	// I haven't figured out a way to search for resource_type=null
	// So just get all permissions and then go through them all
	if resourceType == "" {
		permSearchOpts.PerPage = 400
	} else {
		permSearchOpts.Search = fmt.Sprintf("resource_type=%s", resourceType)
	}
	validPermissions, _, err := client.Permissions.List(context.Background(), *permSearchOpts)
	if err != nil {
		return diag.FromErr(err)
	}
	permNames := d.Get("permission_names").(*schema.Set).List()
	permMap := make(map[string]int)
	var permIDs []int
	for _, x := range *validPermissions.Results {
		if resourceType == "" {
			if x.ResourceType == nil {
				permMap[*x.Name] = *x.ID
			}
		} else {
			permMap[*x.Name] = *x.ID
		}
	}
	for x := range permNames {
		if permID, ok := permMap[permNames[x].(string)]; ok {
			permIDs = append(permIDs, permID)
		} else {
			return diag.Errorf("%s is not a valid permission for resource type %s", permNames[x].(string), resourceType)
		}
	}
	createBody.Filter.PermissionIDs = &permIDs

	if loc, ok := d.GetOk("location_ids"); ok {
		rawLocationIDs := loc.(*schema.Set).List()
		locationIDs := []int{}
		for x := range rawLocationIDs {
			locationIDs = append(locationIDs, rawLocationIDs[x].(int))
		}
		createBody.Filter.LocationIDs = &locationIDs
	}

	if org, ok := d.GetOk("organization_ids"); ok {
		if resourceType != "Location" {
			rawOrganizationIDs := org.(*schema.Set).List()
			organizationIDs := []int{}
			for x := range rawOrganizationIDs {
				organizationIDs = append(organizationIDs, rawOrganizationIDs[x].(int))
			}
			createBody.Filter.OrganizationIDs = &organizationIDs
		} else {
			return diag.Errorf("organization_ids cannot be specified for a resource_type of Location")
		}
	}

	if ovr, ok := d.GetOk("override"); ok {
		override := ovr.(bool)
		createBody.Filter.Override = &(override)
	}

	if srch, ok := d.GetOk("search"); ok {
		search := srch.(string)
		createBody.Filter.Search = &search
	}

	filter, _, err := client.Filters.Create(context.Background(), *createBody)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(*filter.ID))

	return resourceFilterRead(ctx, d, meta)
}

func resourceFilterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	filterID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resourceType := d.Get("resource_type").(string)

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
		if resourceType != "Location" {
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
		} else {
			return diag.Errorf("organization_ids cannot be specified for a resource_type of Location")
		}
	}

	if d.HasChange("override") {
		if ovr, ok := d.GetOk("override"); ok {
			override := ovr.(bool)
			updateBody.Filter.Override = &(override)
		}
	}

	if d.HasChange("permission_names") {
		// validate the permission names passed in

		// For the miscellaneous role_type which has a value of null,
		// I haven't figured out a way to search for resource_type=null
		// So just get all permissions and then go through them all
		permSearchOpts := new(gosatellite.PermissionsListOptions)
		if resourceType == "" {
			permSearchOpts.PerPage = 400
		} else {
			permSearchOpts.Search = fmt.Sprintf("resource_type=%s", resourceType)
		}

		validPermissions, _, err := client.Permissions.List(context.Background(), *permSearchOpts)
		if err != nil {
			return diag.FromErr(err)
		}
		permNames := d.Get("permission_names").(*schema.Set).List()
		permMap := make(map[string]int)
		var permIDs []int
		for _, x := range *validPermissions.Results {
			if resourceType == "" {
				if x.ResourceType == nil {
					permMap[*x.Name] = *x.ID
				}
			} else {
				permMap[*x.Name] = *x.ID
			}
		}
		for x := range permNames {
			if permID, ok := permMap[permNames[x].(string)]; ok {
				permIDs = append(permIDs, permID)
			} else {
				return diag.Errorf("%s is not a valid permission for resource type %s", permNames[x].(string), resourceType)
			}
		}
		updateBody.Filter.PermissionIDs = &permIDs
	}

	if d.HasChange("search") {
		if srch, ok := d.GetOk("search"); ok {
			search := srch.(string)
			updateBody.Filter.Search = &search
		}
	}

	_, _, err = client.Filters.Update(context.Background(), filterID, *updateBody)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFilterRead(ctx, d, meta)
}

func resourceFilterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	filterID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Filters.Delete(context.Background(), filterID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

var resourceTypeList = []string{
	"",
	"AnsibleRole",
	"AnsibleVariable",
	"Architecture",
	"Audit",
	"AuthSource",
	"Bookmark",
	"ComputeProfile",
	"ComputeResource",
	"ConfigGroup",
	"ConfigReport",
	"DiscoveryRule",
	"Domain",
	"Environment",
	"ExternalUsergroup",
	"FactValue",
	"Filter",
	"ForemanOpenscap::ArfReport",
	"ForemanOpenscap::Policy",
	"ForemanOpenscap::ScapContent",
	"ForemanOpenscap::TailoringFile",
	"ForemanTasks::RecurringLogic",
	"ForemanTasks::Task",
	"ForemanVirtWhoConfigure::Config",
	"Host",
	"HostClass",
	"Hostgroup",
	"HttpProxy",
	"InsightsHit",
	"Image",
	"JobInvocation",
	"JobTemplate",
	"Katello::ActivationKey",
	"Katello::ContentView",
	"Katello::GpgKey",
	"Katello::HostCollection",
	"Katello::KTEnvironment",
	"Katello::Product",
	"Katello::Subscription",
	"Katello::SyncPlan",
	"KeyPair",
	"Location",
	"MailNotification",
	"Medium",
	"Model",
	"Operatingsystem",
	"Organization",
	"Parameter",
	"PersonalAccessToken",
	"ProvisioningTemplate",
	"Ptable",
	"Puppetclass",
	"PuppetclassLookupKey",
	"Realm",
	"RemoteExecutionFeature",
	"Report",
	"ReportTemplate",
	"Role",
	"Setting",
	"SmartProxy",
	"SshKey",
	"Subnet",
	"Template",
	"TemplateInvocation",
	"Trend",
	"User",
	"Usergroup",
	"VariableLookupKey",
}
