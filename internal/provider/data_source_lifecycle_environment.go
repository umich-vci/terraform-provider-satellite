package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func dataSourceLifeCycleEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLifeCycleEnvironmentRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLifeCycleEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}

	leList := *le.Results

	if len(leList) == 0 {
		return fmt.Errorf("No Lifecyle Environment found")
	}

	if len(leList) > 1 {
		return fmt.Errorf("%d Lifecyle Environments found, adjust arguments so only 1 is returned", len(leList))
	}

	d.SetId(strconv.Itoa(int(*leList[0].ID)))

	//leList[0].Counts
	//leList[0].Library
	//leList[0].Permissions
	//leList[0].Prior
	//leList[0].RegistryNamePattern
	//leList[0].RegistryUnauthenticatedPull
	//leList[0].Successor

	organization := make(map[string]interface{})
	organization["id"] = leList[0].Organization.ID
	organization["name"] = leList[0].Organization.Name
	organization["label"] = leList[0].Organization.Label

	d.Set("created_at", leList[0].CreatedAt)
	d.Set("description", leList[0].Description)
	d.Set("label", leList[0].Label)
	d.Set("name", leList[0].Name)
	d.Set("organization", organization)
	d.Set("organization_id", leList[0].OrganizationID)
	d.Set("updated_at", leList[0].UpdatedAt)

	return nil
}
