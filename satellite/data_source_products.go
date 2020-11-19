package satellite

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gosatellite"
)

func dataSourceProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProductsRead,
		Schema: map[string]*schema.Schema{
			"organization_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"red_hat_only": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"product_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"products": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cp_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpg_key_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_sync": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_sync_words": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"repository_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceProductsRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}

	pSearch := new(gosatellite.ProductSearch)

	if oID, ok := d.GetOk("organization_id"); ok {
		orgID := oID.(int)
		pSearch.OrganizationID = &orgID
	}

	if rhOnly, ok := d.GetOk("red_hat_only"); ok {
		redHatOnly := rhOnly.(bool)
		pSearch.RedHatOnly = &redHatOnly
	}

	if pName, ok := d.GetOk("product_name"); ok {
		productName := pName.(string)
		pSearch.Name = &productName
	}

	products, _, err := client.Products.ListProducts(context.Background(), *pSearch)
	if err != nil {
		return err
	}

	//d.SetId(strconv.Itoa(orgID))

	productList := []map[string]interface{}{}
	for _, product := range *products.Results {
		prod := map[string]interface{}{}
		prod["cp_id"] = product.CpID
		prod["description"] = product.Description
		prod["gpg_key_id"] = product.GPGKeyID
		prod["id"] = product.ID
		prod["label"] = product.Label
		prod["last_sync"] = product.LastSync
		prod["last_sync_words"] = product.LastSyncWords
		prod["name"] = product.Name
		prod["provider_id"] = product.ProviderID
		prod["repository_count"] = product.RepositoryCount
		productList = append(productList, prod)
	}

	d.Set("products", productList)

	return nil
}
