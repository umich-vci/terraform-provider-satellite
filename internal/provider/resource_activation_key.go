package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func resourceActivationKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceActivationKeyCreate,
		Read:   resourceActivationKeyRead,
		Update: resourceActivationKeyUpdate,
		Delete: resourceActivationKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"content_view_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"host_collection_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"max_hosts": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"unlimited_hosts": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceActivationKeyRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	akID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	activationKey, resp, err := client.ActivationKeys.Get(context.Background(), akID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	// set values we can directly set from struct
	d.Set("organization_id", activationKey.OrganizationID)
	d.Set("name", activationKey.Name)
	d.Set("content_view_id", activationKey.ContentViewID)
	d.Set("description", activationKey.Description)
	d.Set("environment_id", activationKey.EnvironmentID)
	d.Set("max_hosts", activationKey.MaxHosts)
	d.Set("unlimited_hosts", activationKey.UnlimitedHosts)

	var hcIDs []int
	for _, x := range *activationKey.HostCollections {
		hcIDs = append(hcIDs, *x.ID)
	}
	d.Set("host_collection_ids", hcIDs)

	return nil
}

func resourceActivationKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	orgID := d.Get("organization_id").(int)
	name := d.Get("name").(string)
	unlimited := d.Get("unlimited_hosts").(bool)

	createBody := new(gosatellite.ActivationKeyCreate)
	createBody.OrganizationID = &orgID
	createBody.Name = &name
	createBody.UnlimitedHosts = &unlimited

	if c, ok := d.GetOk("content_view_id"); ok {
		cvID := c.(int)
		createBody.ContentViewID = &cvID
	}

	if d, ok := d.GetOk("description"); ok {
		desc := d.(string)
		createBody.Description = &desc
	}

	if e, ok := d.GetOk("environment_id"); ok {
		eID := e.(int)
		createBody.EnvironmentID = &eID
	}

	if m, ok := d.GetOk("max_hosts"); ok {
		max := m.(int)
		createBody.MaxHosts = &max
	}

	activationKey, _, err := client.ActivationKeys.Create(context.Background(), *createBody)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(*activationKey.ID))

	if hc, ok := d.GetOk("host_collection_ids"); ok {
		rawHCIDs := hc.(*schema.Set).List()
		hcIDs := []int{}
		for x := range rawHCIDs {
			hcIDs = append(hcIDs, rawHCIDs[x].(int))
		}
		_, _, err := client.ActivationKeys.AssociateHostCollections(context.Background(), *activationKey.ID, hcIDs)
		if err != nil {
			return err
		}
	}

	return resourceActivationKeyRead(d, meta)
}

func resourceActivationKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	akID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	update := false

	updateBody := new(gosatellite.ActivationKeyUpdate)
	if d.HasChange("organization_id") {
		orgID := d.Get("organization_id").(int)
		updateBody.OrganizationID = &orgID
		update = true
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateBody.Name = &name
		update = true
	}
	if d.HasChange("content_view_id") {
		cvID := d.Get("content_view_id").(int)
		updateBody.ContentViewID = &cvID
		update = true
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateBody.Description = &description
		update = true
	}
	if d.HasChange("environment_id") {
		eID := d.Get("environment_id").(int)
		updateBody.EnvironmentID = &eID
		update = true
	}
	if d.HasChange("max_hosts") {
		maxHosts := d.Get("max_hosts").(int)
		updateBody.MaxHosts = &maxHosts
		update = true
	}
	if d.HasChange("unlimited_hosts") {
		unlimited := d.Get("unlimited_hosts").(bool)
		updateBody.UnlimitedHosts = &unlimited
		update = true
	}

	if update {
		_, _, err = client.ActivationKeys.Update(context.Background(), akID, *updateBody)
		if err != nil {
			return err
		}
	}

	if d.HasChange("host_collection_ids") {
		oldHC, newHC := d.GetChange("host_collection_ids")
		rawOldHC := oldHC.(*schema.Set).List()
		rawNewHC := newHC.(*schema.Set).List()

		hcRemoveList := []int{}
		for x := range rawOldHC {
			found := false
			for y := range rawNewHC {
				if rawOldHC[x].(int) == rawNewHC[y].(int) {
					found = true
				}
			}
			if !found {
				hcRemoveList = append(hcRemoveList, rawOldHC[x].(int))
			}
		}

		hcAddList := []int{}
		for x := range rawNewHC {
			found := false
			for y := range rawOldHC {
				if rawNewHC[x].(int) == rawOldHC[y].(int) {
					found = true
				}
			}
			if !found {
				hcAddList = append(hcAddList, rawNewHC[x].(int))
			}
		}

		if len(hcAddList) > 0 {
			_, _, err := client.ActivationKeys.AssociateHostCollections(context.Background(), akID, hcAddList)
			if err != nil {
				return err
			}
		}

		if len(hcRemoveList) > 0 {
			_, _, err := client.ActivationKeys.DisassociateHostCollections(context.Background(), akID, hcRemoveList)
			if err != nil {
				return err
			}
		}

	}

	return resourceActivationKeyRead(d, meta)
}

func resourceActivationKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	akID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.ActivationKeys.Delete(context.Background(), akID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
