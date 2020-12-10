package satellite

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func resourceFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceFilterCreate,
		Read:   resourceFilterRead,
		Update: resourceFilterUpdate,
		Delete: resourceFilterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"permission_names": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"role_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"resource_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(resourceTypeList, false),
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
			"permissions": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"permission_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
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
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
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

func resourceFilterCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	roleID := d.Get("role_id").(int)
	resourceType := d.Get("resource_type").(string)
	createBody := new(gosatellite.FilterCreate)
	createBody.Filter.RoleID = &roleID

	// validate the permission names passed in
	permSearchBody := new(gosatellite.PermissionsSearch)

	// For the miscellaneous role_type which has a value of null,
	// I haven't figured out a way to search for resource_type=null
	// So just get all permissions and then go through them all
	if resourceType == "" {
		perPage := 400
		permSearchBody.PerPage = &perPage
	} else {
		search := "resource_type=" + resourceType
		permSearchBody.Search = &search
	}
	validPermissions, _, err := client.Permissions.ListPermissions(context.Background(), *permSearchBody)
	if err != nil {
		return err
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
			return fmt.Errorf("%s is not a valid permission for resource type %s", permNames[x].(string), resourceType)
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
			return fmt.Errorf("organization_ids cannot be specified for a resource_type of Location")
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
			return fmt.Errorf("organization_ids cannot be specified for a resource_type of Location")
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
		permSearchBody := new(gosatellite.PermissionsSearch)
		if resourceType == "" {
			perPage := 400
			permSearchBody.PerPage = &perPage
		} else {
			search := "resource_type=" + resourceType
			permSearchBody.Search = &search
		}

		validPermissions, _, err := client.Permissions.ListPermissions(context.Background(), *permSearchBody)
		if err != nil {
			return err
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
				return fmt.Errorf("%s is not a valid permission for resource type %s", permNames[x].(string), resourceType)
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
