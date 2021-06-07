package provider

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSubscriptionManifest() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage a subscription manifest attached to a Red Hat Satellite organization.",

		CreateContext: resourceSubscriptionManifestCreate,
		ReadContext:   resourceSubscriptionManifestRead,
		UpdateContext: resourceSubscriptionManifestUpdate,
		DeleteContext: resourceSubscriptionManifestDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description: "The organization ID you want to attach the manifest to.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"manifest": {
				Description: "A Base64 encoded string of a manifest zip file downloaded from Red Hat Subscription Management. Most easily used in conjunction with [`rhsm_allocation_manifest` resource from the RHSM provider](https://registry.terraform.io/providers/umich-vci/rhsm/latest/docs/resources/allocation_manifest).",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"history": {
				Description: "A list of objects containing information on operations peformed on the manifest.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created": {
							Description: "TODO",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "TODO",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"status": {
							Description: "TODO",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"status_message": {
							Description: "TODO",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceSubscriptionManifestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	hist, resp, err := client.Manifests.GetHistory(context.Background(), orgID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	histList := []map[string]interface{}{}

	for _, x := range *hist {
		histItem := make(map[string]interface{})
		histItem["created"] = x.Created
		histItem["id"] = x.ID
		histItem["status"] = x.Status
		histItem["status_message"] = x.StatusMessage
		histList = append(histList, histItem)
	}

	d.Set("organization_id", orgID)
	d.Set("history", histList)

	return nil
}

func resourceSubscriptionManifestCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	orgID := d.Get("organization_id").(int)

	manifestString := d.Get("manifest").(string)
	manifest, err := base64.StdEncoding.DecodeString(manifestString)
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = client.Manifests.Upload(context.Background(), orgID, nil, manifest, "manifest.zip")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(orgID))

	return resourceSubscriptionManifestRead(ctx, d, meta)
}

func resourceSubscriptionManifestUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Manifests.Refresh(context.Background(), orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSubscriptionManifestRead(ctx, d, meta)
}

func resourceSubscriptionManifestDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Manifests.Delete(context.Background(), orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
