package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func dataSourceAuthSourceLDAP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAuthSourceLDAPRead,
		Schema: map[string]*schema.Schema{
			"search": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attr_firstname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attr_lastname": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"attr_login": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attr_mail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attr_photo": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"base_dn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// "external_usergroups": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeMap,
			// 	},
			// },
			"groups_base": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ldap_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// "locations": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeMap,
			// 	},
			// },
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"onthefly_register": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			// "organizations": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeMap,
			// 	},
			// },
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tls": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use_netgroups": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"usergroup_sync": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceAuthSourceLDAPRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	searchString := d.Get("search").(string)

	opt := new(gosatellite.AuthSourceLDAPsListOptions)
	opt.Search = searchString

	authSourcesResp, _, err := client.AuthSourceLDAPs.List(context.Background(), opt)
	if err != nil {
		return err
	}

	authSources := *authSourcesResp.Results

	if len(authSources) == 0 {
		return fmt.Errorf("No LDAP Auth Sources found for search string %s", searchString)
	}

	if len(authSources) > 1 {
		return fmt.Errorf("%d LDAP Auth Sources found for search string %s", len(authSources), searchString)
	}

	d.SetId(strconv.Itoa(*authSources[0].ID))

	d.Set("account", *authSources[0].Account)
	d.Set("attr_firstname", *authSources[0].AttrFirstName)
	d.Set("attr_lastname", *authSources[0].AttrLastName)
	d.Set("attr_login", *authSources[0].AttrLogin)
	d.Set("attr_mail", *authSources[0].AttrMail)
	d.Set("attr_photo", *authSources[0].AttrPhoto)
	d.Set("base_dn", *authSources[0].BaseDN)
	d.Set("created_at", *authSources[0].CreatedAt)
	d.Set("groups_base", *authSources[0].GroupsBase)
	d.Set("host", *authSources[0].Host)
	d.Set("ldap_filter", *authSources[0].LDAPFilter)
	d.Set("name", *authSources[0].Name)
	d.Set("onthefly_register", *authSources[0].OnTheFlyRegister)
	d.Set("port", *authSources[0].Port)
	d.Set("server_type", *authSources[0].ServerType)
	d.Set("tls", *authSources[0].TLS)
	d.Set("type", *authSources[0].Type)
	d.Set("use_netgroups", *authSources[0].UseNetGroups)
	d.Set("updated_at", *authSources[0].UpdatedAt)
	d.Set("usergroup_sync", *authSources[0].UserGroupSync)

	return nil
}
