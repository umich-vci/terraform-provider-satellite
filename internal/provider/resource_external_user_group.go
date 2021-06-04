package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func resourceExternalUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceExternalUserGroupCreate,
		Read:   resourceExternalUserGroupRead,
		Update: resourceExternalUserGroupUpdate,
		Delete: resourceExternalUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"auth_source_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"user_group_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"auth_source_ldap": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceExternalUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	ugID := d.Get("user_group_id").(int)

	eugID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	eug, resp, err := client.ExternalUserGroups.Get(context.Background(), ugID, eugID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("name", eug.Name)
	d.Set("auth_source_id", eug.AuthSourceLDAP.ID)

	authSourceLDAP := make(map[string]interface{})
	authSourceLDAP["id"] = eug.AuthSourceLDAP.ID
	authSourceLDAP["name"] = eug.AuthSourceLDAP.Name
	authSourceLDAP["type"] = eug.AuthSourceLDAP.Type
	d.Set("auth_source_ldap", authSourceLDAP)

	return nil
}

func resourceExternalUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	ugID := d.Get("user_group_id").(int)
	asID := d.Get("auth_source_id").(int)

	createBody := new(gosatellite.ExternalUserGroupCreate)
	createBody.ExternalUserGroup.AuthSourceID = &asID
	createBody.ExternalUserGroup.Name = &name

	eug, _, err := client.ExternalUserGroups.Create(context.Background(), ugID, *createBody)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(*eug.ID))

	return resourceExternalUserGroupRead(d, meta)
}

func resourceExternalUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	eugID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	ugID := d.Get("user_group_id").(int)

	updateBody := new(gosatellite.ExternalUserGroupUpdate)
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateBody.ExternalUserGroup.Name = &name
	}
	if d.HasChange("auth_source_id") {
		asID := d.Get("admin").(int)
		updateBody.ExternalUserGroup.AuthSourceID = &asID
	}

	_, _, err = client.ExternalUserGroups.Update(context.Background(), ugID, eugID, *updateBody)
	if err != nil {
		return err
	}
	return resourceExternalUserGroupRead(d, meta)
}

func resourceExternalUserGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	eugID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	ugID := d.Get("user_group_id").(int)

	_, _, err = client.ExternalUserGroups.Delete(context.Background(), ugID, eugID)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
