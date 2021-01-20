package satellite

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{
		Username:      d.Get("username").(string),
		Password:      d.Get("password").(string),
		SatelliteHost: d.Get("satellite_host").(string),
		SSLVerify:     d.Get("ssl_verify").(bool),
	}

	return config, nil
}

// Provider returns a terraform resource provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SATELLITE_USERNAME", nil),
				Description: "A Satellite username.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SATELLITE_PASSWORD", nil),
				Description: "A Satellite password.",
			},
			"satellite_host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SATELLITE_HOST", nil),
				Description: "The Red Hat Satellite hostname",
			},
			"ssl_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Perform SSL verification",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"satellite_activation_key":        resourceActivationKey(),
			"satellite_filter":                resourceFilter(),
			"satellite_host_collection":       resourceHostCollection(),
			"satellite_location":              resourceLocation(),
			"satellite_organization":          resourceOrganization(),
			"satellite_role":                  resourceRole(),
			"satellite_subscription_manifest": resourceSubscriptionManifest(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"satellite_location":     dataSourceLocation(),
			"satellite_organization": dataSourceOrganization(),
			"satellite_permissions":  dataSourcePermissions(),
			"satellite_products":     dataSourceProducts(),
		},
		ConfigureFunc: providerConfigure,
	}
}
