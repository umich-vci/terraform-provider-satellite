package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage a Red Hat Satellite organization.",

		CreateContext: resourceOrganizationCreate,
		ReadContext:   resourceOrganizationRead,
		UpdateContext: resourceOrganizationUpdate,
		DeleteContext: resourceOrganizationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the organization.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "A description of the organization.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"label": {
				Description: "The label of the organization. If not set, Satellite will use the `name` as the label.  This field can only be set at creation time. If not being set explicitly you will probably want to use `ignore_changes` on this in the lifecycle block.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"hosts_count": {
				Description: "A count of how many hosts are registered to the organization.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"title": {
				Description: "The title of the organization.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	org, resp, err := client.Organizations.Get(context.Background(), orgID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	d.Set("description", org.Description)
	d.Set("hosts_count", org.HostsCount)
	d.Set("label", org.Label)
	d.Set("name", org.Name)
	d.Set("title", org.Title)

	return nil
}

func resourceOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	createBody := new(gosatellite.OrganizationCreate)
	createBody.Organization.Name = d.Get("name").(string)

	org, _, err := client.Organizations.Create(context.Background(), *createBody)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(*org.ID))

	return resourceOrganizationRead(ctx, d, meta)
}

func resourceOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updateBody := new(gosatellite.OrganizationUpdate)
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateBody.Organization.Description = &description
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateBody.Organization.Name = &name
	}

	_, _, err = client.Organizations.Update(context.Background(), orgID, *updateBody)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOrganizationRead(ctx, d, meta)
}

func resourceOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Organizations.Delete(context.Background(), orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
