package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func resourceLocation() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage a Red Hat Satellite location.",

		CreateContext: resourceLocationCreate,
		ReadContext:   resourceLocationRead,
		UpdateContext: resourceLocationUpdate,
		DeleteContext: resourceLocationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "A name for the location.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Description: "A description of the location.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"parent_id": {
				Description: "The ID of a parent for this location. This allows you to nest locations. If not set, a top level location is created.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
}

func resourceLocationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	locationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	location, resp, err := client.Locations.Get(context.Background(), locationID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	d.Set("name", location.Name)
	d.Set("description", location.Description)
	d.Set("parent_id", location.ParentID)

	return nil
}

func resourceLocationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	name := d.Get("name").(string)
	createBody := new(gosatellite.LocationCreate)
	createBody.Location.Name = &name

	if d, ok := d.GetOk("description"); ok {
		description := d.(string)
		createBody.Location.Description = &description
	}

	if p, ok := d.GetOk("parent_id"); ok {
		parentID := p.(int)
		createBody.Location.ParentID = &parentID
	}

	location, _, err := client.Locations.Create(context.Background(), *createBody)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(*location.ID))

	return resourceLocationRead(ctx, d, meta)
}

func resourceLocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	locationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updateBody := new(gosatellite.LocationUpdate)
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateBody.Location.Description = &description
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateBody.Location.Name = &name
	}
	if d.HasChange("parent_id") {
		parentID := d.Get("parent_id").(int)
		updateBody.Location.ParentID = &parentID
	}

	_, _, err = client.Locations.Update(context.Background(), locationID, *updateBody)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLocationRead(ctx, d, meta)
}

func resourceLocationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	locationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Locations.Delete(context.Background(), locationID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
