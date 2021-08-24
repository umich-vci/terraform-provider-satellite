package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func resourceExternalUserGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage an external user group in Red Hat Satellite.",

		CreateContext: resourceExternalUserGroupCreate,
		ReadContext:   resourceExternalUserGroupRead,
		UpdateContext: resourceExternalUserGroupUpdate,
		DeleteContext: resourceExternalUserGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "The name of the external user group.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"auth_source_id": {
				Description: "The ID of the authentication source that contains the external user group.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"user_group_id": {
				Description: "The ID of the user group that the external user group should be associated with.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"auth_source_ldap": {
				Description: "A list of objects containing the authentication source the associated with the external user group.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceExternalUserGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	ugID := d.Get("user_group_id").(int)

	eugID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	eug, resp, err := client.ExternalUserGroups.Get(context.Background(), ugID, eugID)
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	d.Set("name", eug.Name)
	d.Set("auth_source_id", eug.AuthSourceLDAP.ID)

	authSourceLDAP := make(map[string]interface{})
	if eug.AuthSourceLDAP != nil {
		authSourceLDAP["id"] = strconv.Itoa(*eug.AuthSourceLDAP.ID)
		authSourceLDAP["name"] = eug.AuthSourceLDAP.Name
		authSourceLDAP["type"] = eug.AuthSourceLDAP.Type
	}
	d.Set("auth_source_ldap", authSourceLDAP)

	return nil
}

func resourceExternalUserGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	name := d.Get("name").(string)
	ugID := d.Get("user_group_id").(int)
	asID := d.Get("auth_source_id").(int)

	createBody := new(gosatellite.ExternalUserGroupCreate)
	createBody.ExternalUserGroup.AuthSourceID = &asID
	createBody.ExternalUserGroup.Name = &name

	eug, _, err := client.ExternalUserGroups.Create(context.Background(), ugID, *createBody)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(*eug.ID))

	return resourceExternalUserGroupRead(ctx, d, meta)
}

func resourceExternalUserGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	eugID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}
	return resourceExternalUserGroupRead(ctx, d, meta)
}

func resourceExternalUserGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	eugID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ugID := d.Get("user_group_id").(int)

	_, _, err = client.ExternalUserGroups.Delete(context.Background(), ugID, eugID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
