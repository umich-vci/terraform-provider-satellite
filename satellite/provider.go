package satellite

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
func Provider() terraform.ResourceProvider {
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
			"satellite_filter":       resourceFilter(),
			"satellite_organization": resourceOrganization(),
			"satellite_role":         resourceRole(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"satellite_organization": dataSourceOrganization(),
		},
		ConfigureFunc: providerConfigure,
	}
}
