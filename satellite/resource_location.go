package satellite

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func resourceLocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceLocationCreate,
		Read:   resourceLocationRead,
		Update: resourceLocationUpdate,
		Delete: resourceLocationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceLocationRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	locationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	location, resp, err := client.Locations.Get(context.Background(), locationID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("name", location.Name)
	d.Set("description", location.Description)
	d.Set("parent_id", location.ParentID)

	return nil
}

func resourceLocationCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}
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
		return err
	}

	d.SetId(strconv.Itoa(*location.ID))

	return resourceLocationRead(d, meta)
}

func resourceLocationUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	locationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
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
		return err
	}

	return resourceLocationRead(d, meta)
}

func resourceLocationDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	locationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Locations.Delete(context.Background(), locationID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
