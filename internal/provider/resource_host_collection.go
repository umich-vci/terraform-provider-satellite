package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func resourceHostCollection() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage a Red Hat Satellite Host Collection.",

		CreateContext: resourceHostCollectionCreate,
		ReadContext:   resourceHostCollectionRead,
		UpdateContext: resourceHostCollectionUpdate,
		DeleteContext: resourceHostCollectionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the host collection.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"organization_id": {
				Description: "The ID of organization that the host collection should be created in. Once set, it cannot be changed without recreating the resource.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "A description of the host collection.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"max_hosts": {
				Description: "The maximum number of hosts allowed to be in the host collection. Should not be set if `unlimited_hosts` is set to `true`.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"unlimited_hosts": {
				Description: "A boolean that controls if an unlimited number of members are allowed in the host collection.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"created_at": {
				Description: "A timestamp containing when the host collection was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "A timestamp containing when the host collection was last changed.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceHostCollectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	hcID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	hc, resp, err := client.HostCollections.Get(context.Background(), hcID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	d.Set("name", hc.Name)
	d.Set("organization_id", hc.OrganizationID)
	d.Set("description", hc.Description)
	d.Set("max_hosts", hc.MaxHosts)
	d.Set("unlimited_hosts", hc.UnlimitedHosts)
	d.Set("created_at", hc.CreatedAt)
	d.Set("updated_at", hc.UpdatedAt)

	return nil
}

func resourceHostCollectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	orgID := d.Get("organization_id").(int)
	unlimited := d.Get("unlimited_hosts").(bool)

	createBody := new(gosatellite.HostCollectionCreate)
	createBody.Name = d.Get("name").(string)
	createBody.UnlimitedHosts = &unlimited

	if _, ok := d.GetOk("description"); ok {
		description := d.Get("description").(string)
		createBody.Description = &description
	}

	if _, ok := d.GetOk("max_hosts"); ok {
		maxHosts := d.Get("max_hosts").(int)
		createBody.MaxHosts = &maxHosts
	}

	hc, _, err := client.HostCollections.Create(context.Background(), orgID, *createBody)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(*hc.ID))

	return resourceHostCollectionRead(ctx, d, meta)
}

func resourceHostCollectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	hcID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updateBody := new(gosatellite.HostCollectionUpdate)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateBody.Name = &name
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateBody.Description = &description
	}

	if d.HasChange("max_hosts") {
		maxHosts := d.Get("max_hosts").(int)
		updateBody.MaxHosts = &maxHosts
	}

	if d.HasChange("unlimited_hosts") {
		unlimited := d.Get("unlimited_hosts").(bool)
		updateBody.UnlimitedHosts = &unlimited
	}

	_, _, err = client.HostCollections.Update(context.Background(), hcID, *updateBody)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceHostCollectionRead(ctx, d, meta)
}

func resourceHostCollectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	hcID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.HostCollections.Delete(context.Background(), hcID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
