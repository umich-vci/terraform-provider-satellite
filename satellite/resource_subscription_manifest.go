package satellite

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSubscriptionManifest() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubscriptionManifestCreate,
		Read:   resourceSubscriptionManifestRead,
		Update: resourceSubscriptionManifestUpdate,
		Delete: resourceSubscriptionManifestDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"manifest": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"history": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceSubscriptionManifestRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	hist, resp, err := client.Manifests.GetHistory(context.Background(), orgID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
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

func resourceSubscriptionManifestCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orgID := d.Get("organization_id").(int)

	manifestString := d.Get("manifest").(string)
	manifest, err := base64.StdEncoding.DecodeString(manifestString)
	if err != nil {
		return err
	}

	_, _, err = client.Manifests.Upload(context.Background(), orgID, nil, manifest, "manifest.zip")
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(orgID))

	return resourceSubscriptionManifestRead(d, meta)
}

func resourceSubscriptionManifestUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Manifests.Refresh(context.Background(), orgID)
	if err != nil {
		return err
	}

	return resourceSubscriptionManifestRead(d, meta)
}

func resourceSubscriptionManifestDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orgID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Manifests.Delete(context.Background(), orgID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
