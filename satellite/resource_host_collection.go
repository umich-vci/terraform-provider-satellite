package satellite

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func resourceHostCollection() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostCollectionCreate,
		Read:   resourceHostCollectionRead,
		Update: resourceHostCollectionUpdate,
		Delete: resourceHostCollectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_hosts": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"unlimited_hosts": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceHostCollectionRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	hcID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	hc, resp, err := client.HostCollections.GetHostCollectionByID(context.Background(), hcID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
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

func resourceHostCollectionCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

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

	hc, _, err := client.HostCollections.CreateHostCollection(context.Background(), orgID, *createBody)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(*hc.ID))

	return resourceHostCollectionRead(d, meta)
}

func resourceHostCollectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	hcID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
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

	_, _, err = client.HostCollections.UpdateHostCollection(context.Background(), hcID, *updateBody)
	if err != nil {
		return err
	}

	return resourceHostCollectionRead(d, meta)
}

func resourceHostCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	hcID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.HostCollections.DeleteHostCollection(context.Background(), hcID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
