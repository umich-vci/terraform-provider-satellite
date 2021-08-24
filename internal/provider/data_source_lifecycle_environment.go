package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func dataSourceLifecycleEnvironment() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to access information about a Red Hat Satellite Lifecycle Environment.",

		ReadContext: dataSourceLifecycleEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the Lifecycle Environment.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"organization_id": {
				Description: "The ID of the organization that contains the Lifecycle Environment.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"search": {
				Description: "A search filter for the Lifecycle Environment search. The search must only return 1 Lifecycle Environment.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"counts": {
				Description: "A map of various counts of attributes of the Lifecycle Environment.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"created_at": {
				Description: "Timestamp of when the Lifecycle Environment was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "A description of the Lifecycle Environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"label": {
				Description: "A label for the Lifecycle Environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"library": {
				Description: "Is the Lifecycle Environment a base Library?",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"organization": {
				Description: "The organization that contains the Lifecycle Environment.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"permissions": {
				Description: "The current Satellite user's permissions for the Lifecycle Environment.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
			},
			"prior": {
				Description: "The Lifecycle Environment directly before this one.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"registry_name_pattern": {
				Description: "TODO",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"registry_unauthenticated_pull": {
				Description: "TODO",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"successor": {
				Description: "The Lifecycle Environment directly after this one.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"updated_at": {
				Description: "Timestamp of when the Lifecycle Environment was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceLifecycleEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	opt := new(gosatellite.LifecycleEnvironmentsListOptions)

	if name, ok := d.GetOk("name"); ok {
		opt.Name = name.(string)
	}

	if orgID, ok := d.GetOk("organization_id"); ok {
		opt.OrganizationID = orgID.(int)
	}

	if search, ok := d.GetOk("search"); ok {
		opt.Search = search.(string)
	}

	le, _, err := client.LifecycleEnvironments.List(context.Background(), opt)
	if err != nil {
		return diag.FromErr(err)
	}

	leList := *le.Results

	if len(leList) == 0 {
		return diag.Errorf("No Lifecyle Environment found")
	}

	if len(leList) > 1 {
		return diag.Errorf("%d Lifecyle Environments found, adjust arguments so only 1 is returned", len(leList))
	}

	d.SetId(strconv.Itoa(int(*leList[0].ID)))

	counts := make(map[string]interface{})
	if leList[0].Counts != nil {
		counts["content_hosts"] = leList[0].Counts.ContentHosts
		counts["content_views"] = leList[0].Counts.ContentViews
		counts["docker_repositories"] = leList[0].Counts.DockerRepositories
		if leList[0].Counts.Errata != nil {
			counts["errata_bugfix"] = leList[0].Counts.Errata.Bugfix
			counts["errata_enhancement"] = leList[0].Counts.Errata.Enhancement
			counts["errata_security"] = leList[0].Counts.Errata.Security
			counts["errata_total"] = leList[0].Counts.Errata.Total
		}
		counts["module_streams"] = leList[0].Counts.ModuleStreams
		counts["os_tree_repositories"] = leList[0].Counts.OSTreeRepositories
		counts["packages"] = leList[0].Counts.Packages
		counts["products"] = leList[0].Counts.Products
		counts["puppet_modules"] = leList[0].Counts.PuppetModules
		counts["yum_repositories"] = leList[0].Counts.YumRepositories
	}

	organization := make(map[string]interface{})
	if leList[0].Organization != nil {
		organization["id"] = strconv.Itoa(*leList[0].Organization.ID)
		organization["name"] = leList[0].Organization.Name
		organization["label"] = leList[0].Organization.Label
	}

	permissions := make(map[string]interface{})
	if leList[0].Permissions != nil {
		permissions["create_lifecycle_environments"] = leList[0].Permissions.CreateLifecycleEnvironments
		permissions["destroy_lifecycle_environments"] = leList[0].Permissions.DestroyLifecycleEnvironments
		permissions["edit_lifecycle_environments"] = leList[0].Permissions.EditLifecycleEnvironments
		permissions["promote_or_remove_content_views_to_environments"] = leList[0].Permissions.PromoteOrRemoveContentViewsToEnvironments
		permissions["view_lifecycle_environments"] = leList[0].Permissions.ViewLifecycleEnvironments
	}

	prior := make(map[string]interface{})
	if leList[0].Prior != nil {
		prior["id"] = strconv.Itoa(*leList[0].Prior.ID)
		prior["name"] = leList[0].Prior.Name
	}

	successor := make(map[string]interface{})
	if leList[0].Successor != nil {
		successor["id"] = strconv.Itoa(*leList[0].Successor.ID)
		successor["name"] = leList[0].Successor.Name
	}

	d.Set("counts", counts)
	d.Set("created_at", leList[0].CreatedAt)
	d.Set("description", leList[0].Description)
	d.Set("label", leList[0].Label)
	d.Set("library", leList[0].Library)
	d.Set("name", leList[0].Name)
	d.Set("organization", organization)
	d.Set("organization_id", leList[0].OrganizationID)
	d.Set("permissions", permissions)
	d.Set("registry_name_pattern", leList[0].RegistryNamePattern)
	d.Set("registry_unauthenticated_pull", leList[0].RegistryUnauthenticatedPull)
	d.Set("prior", prior)
	d.Set("successor", successor)
	d.Set("updated_at", leList[0].UpdatedAt)

	return nil
}
