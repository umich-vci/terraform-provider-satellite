package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gosatellite"
)

func dataSourceAuthSourceLDAP() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to access information about a Red Hat Satellite LDAP Authentication Source.",

		ReadContext: dataSourceAuthSourceLDAPRead,

		Schema: map[string]*schema.Schema{
			"search": {
				Description:  "A search filter for the LDAP Authentication Source search. The search must only return 1 authentication source.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"account": {
				Description: "The DN of the LDAP Bind Account.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"attr_firstname": {
				Description: "The LDAP attribute that maps to first name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"attr_lastname": {
				Description: "The LDAP attribute that maps to last name.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"attr_login": {
				Description: "The LDAP attribute that maps to username.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"attr_mail": {
				Description: "The LDAP attribute that maps to email address.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"attr_photo": {
				Description: "The LDAP attribute that maps to a photo.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"base_dn": {
				Description: "The base DN from which LDAP searches will be performed.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "Timestamp of when the LDAP authentication source was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			// "external_usergroups": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeMap,
			// 	},
			// },
			"groups_base": {
				Description: "The base DN from which LDAP searches for groups will be performed.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"host": {
				Description: "The hostname of the LDAP server.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ldap_filter": {
				Description: "TODO",
				Type:        schema.TypeString,
				Computed:    true,
			},
			// "locations": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeMap,
			// 	},
			// },
			"name": {
				Description: "The name of the LDAP authentication source.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"onthefly_register": {
				Description: "TODO",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			// "organizations": &schema.Schema{
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeMap,
			// 	},
			// },
			"port": {
				Description: "The port the LDAP server is listening on.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"server_type": {
				Description: "The type of the LDAP server.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tls": {
				Description: "Is TLS enabled for the LDAP server?",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"type": {
				Description: "TODO",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"use_netgroups": {
				Description: "TODO",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"updated_at": {
				Description: "Timestamp of when the LDAP authentication source was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"usergroup_sync": {
				Description: "TODO",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func dataSourceAuthSourceLDAPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	searchString := d.Get("search").(string)

	opt := new(gosatellite.AuthSourceLDAPsListOptions)
	opt.Search = searchString

	authSourcesResp, _, err := client.AuthSourceLDAPs.List(context.Background(), opt)
	if err != nil {
		return diag.FromErr(err)
	}

	authSources := *authSourcesResp.Results

	if len(authSources) == 0 {
		return diag.Errorf("No LDAP Auth Sources found for search string %s", searchString)
	}

	if len(authSources) > 1 {
		return diag.Errorf("%d LDAP Auth Sources found for search string %s", len(authSources), searchString)
	}

	d.SetId(strconv.Itoa(*authSources[0].ID))

	d.Set("account", authSources[0].Account)
	d.Set("attr_firstname", authSources[0].AttrFirstName)
	d.Set("attr_lastname", authSources[0].AttrLastName)
	d.Set("attr_login", authSources[0].AttrLogin)
	d.Set("attr_mail", authSources[0].AttrMail)
	d.Set("attr_photo", authSources[0].AttrPhoto)
	d.Set("base_dn", authSources[0].BaseDN)
	d.Set("created_at", authSources[0].CreatedAt)
	d.Set("groups_base", authSources[0].GroupsBase)
	d.Set("host", authSources[0].Host)
	d.Set("ldap_filter", authSources[0].LDAPFilter)
	d.Set("name", authSources[0].Name)
	d.Set("onthefly_register", authSources[0].OnTheFlyRegister)
	d.Set("port", authSources[0].Port)
	d.Set("server_type", authSources[0].ServerType)
	d.Set("tls", authSources[0].TLS)
	d.Set("type", authSources[0].Type)
	d.Set("use_netgroups", authSources[0].UseNetGroups)
	d.Set("updated_at", authSources[0].UpdatedAt)
	d.Set("usergroup_sync", authSources[0].UserGroupSync)

	return nil
}
