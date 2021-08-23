package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func dataSourceContentView() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceContentViewRead,
		Schema: map[string]*schema.Schema{
			"composite": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"environment_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"noncomposite": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nondefault": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"search": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"without": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"activation_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"auto_publish": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"component_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"force_puppet_environment": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_published": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"repositories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"repository_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"solve_dependencies": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceContentViewRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient).Client

	opt := new(gosatellite.ContentViewsListOptions)

	if composite, ok := d.GetOk("composite"); ok {
		opt.Composite = composite.(bool)
	}

	if envID, ok := d.GetOk("environment_id"); ok {
		opt.EnvironmentID = envID.(int)
	}

	if name, ok := d.GetOk("name"); ok {
		opt.Name = name.(string)
	}

	if noncomposite, ok := d.GetOk("noncomposite"); ok {
		opt.Noncomposite = noncomposite.(bool)
	}

	if nondefault, ok := d.GetOk("nondefault"); ok {
		opt.Nondefault = nondefault.(bool)
	}

	if orgID, ok := d.GetOk("organization_id"); ok {
		opt.OrganizationID = orgID.(int)
	}

	if search, ok := d.GetOk("search"); ok {
		opt.Search = search.(string)
	}

	if wo, ok := d.GetOk("without"); ok {
		rawWithout := wo.(*schema.Set).List()
		without := []string{}
		for x := range rawWithout {
			without = append(without, rawWithout[x].(string))
		}
		opt.Without = without
	}

	cv, _, err := client.ContentViews.List(context.Background(), opt)
	if err != nil {
		return err
	}

	cvList := *cv.Results

	if len(cvList) == 0 {
		return fmt.Errorf("No Content Views found")
	}

	if len(cvList) > 1 {
		return fmt.Errorf("%d Content Views found, adjust arguments so only 1 is returned", len(cvList))
	}

	d.SetId(strconv.Itoa(int(*cvList[0].ID)))

	activationKeys := []map[string]interface{}{}
	for _, x := range *cvList[0].ActivationKeys {
		activationKey := make(map[string]interface{})
		activationKey["id"] = x.ID
		activationKey["name"] = x.Name
		activationKeys = append(activationKeys, activationKey)
	}

	environments := []map[string]interface{}{}
	for _, x := range *cvList[0].Environments {
		environment := make(map[string]interface{})
		environment["id"] = x.ID
		environment["name"] = x.Name
		environment["label"] = x.Label
		environments = append(environments, environment)
	}

	organization := make(map[string]interface{})
	organization["id"] = cvList[0].Organization.ID
	organization["name"] = cvList[0].Organization.Name
	organization["label"] = cvList[0].Organization.Label

	repositories := []map[string]interface{}{}
	for _, x := range *cvList[0].Repositories {
		repository := make(map[string]interface{})
		repository["id"] = x.ID
		repository["name"] = x.Name
		repository["label"] = x.Label
		repository["content_type"] = x.ContentType
		repositories = append(repositories, repository)
	}

	versions := []map[string]interface{}{}
	for _, x := range *cvList[0].Versions {
		version := make(map[string]interface{})
		version["id"] = x.ID
		version["name"] = x.EnvironmentIDs
		version["published"] = x.Published
		version["version"] = x.Version
		versions = append(versions, version)
	}

	d.Set("activation_keys", activationKeys)
	d.Set("auto_publish", cvList[0].AutoPublish)
	d.Set("component_ids", cvList[0].ComponentIDs)
	d.Set("composite", cvList[0].Composite)
	d.Set("created_at", cvList[0].CreatedAt)
	d.Set("default", cvList[0].Default)
	d.Set("description", cvList[0].Description)
	d.Set("environments", environments)
	d.Set("force_puppet_environment", cvList[0].ForcePuppetEnvironment)
	d.Set("label", cvList[0].Label)
	d.Set("last_published", cvList[0].LastPublished)
	d.Set("latest_version", cvList[0].LatestVersion)
	d.Set("name", cvList[0].Name)
	d.Set("next_version", cvList[0].NextVersion)
	d.Set("organization", organization)
	d.Set("organization_id", cvList[0].OrganizationID)
	d.Set("repositories", repositories)
	d.Set("repository_ids", cvList[0].RepositoryIDs)
	d.Set("solve_dependencies", cvList[0].SolveDependencies)
	d.Set("updated_at", cvList[0].UpdatedAt)
	d.Set("version_count", cvList[0].VersionCount)
	d.Set("versions", versions)

	return nil
}
