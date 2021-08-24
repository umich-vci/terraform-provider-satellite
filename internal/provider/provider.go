package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"username": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("SATELLITE_USERNAME", nil),
					Description: "This is the username to use to access the Red Hat Satellite server. This must be provided in the config or in the environment variable `SATELLITE_USERNAME`.",
				},
				"password": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SATELLITE_PASSWORD", nil),
					Description: "This is the password to use to access the Red Hat Satellite server. This must be provided in the config or in the environment variable `SATELLITE_PASSWORD`.",
				},
				"satellite_host": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("SATELLITE_HOST", nil),
					Description: "This is the hostname or IP address of the Red Hat Satellite server. This must be provided in the config or in the environment variable `SATELLITE_HOST`.",
				},
				"ssl_verify": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Should we validate the SSL certificate presented by the Satellite server?.",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"satellite_auth_source_ldap":      dataSourceAuthSourceLDAP(),
				"satellite_content_view":          dataSourceContentView(),
				"satellite_lifecycle_environment": dataSourceLifeCycleEnvironment(),
				"satellite_location":              dataSourceLocation(),
				"satellite_organization":          dataSourceOrganization(),
				"satellite_permissions":           dataSourcePermissions(),
				"satellite_products":              dataSourceProducts(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"satellite_activation_key":        resourceActivationKey(),
				"satellite_external_user_group":   resourceExternalUserGroup(),
				"satellite_filter":                resourceFilter(),
				"satellite_host_collection":       resourceHostCollection(),
				"satellite_location":              resourceLocation(),
				"satellite_organization":          resourceOrganization(),
				"satellite_role":                  resourceRole(),
				"satellite_subscription_manifest": resourceSubscriptionManifest(),
				"satellite_user_group":            resourceUserGroup(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	Client gosatellite.Client
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		userAgent := p.UserAgent("terraform-provider-umich", version)
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		satelliteHost := d.Get("satellite_host").(string)
		sslVerify := d.Get("ssl_verify").(bool)

		config := new(gosatellite.Config)
		config.Username = username
		config.Password = password
		config.SatelliteHost = satelliteHost
		config.SSLVerify = sslVerify

		client, err := gosatellite.NewClient(config)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		client.UserAgent = userAgent

		return &apiClient{Client: *client}, nil
	}
}
