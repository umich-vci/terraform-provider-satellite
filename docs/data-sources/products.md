# satellite\_products

Use this data source to access information about available Red Hat Satellite products.

## Example Usage

```hcl
data "satellite_products" "rh_products" {
    organization_id = 2
    red_hat_only = true
}
```

## Argument Reference

* `organization_id` - (Optional) An Organization ID to filter the product search on.

* `red_hat_only` - (Optional) A boolean that controls if the search should only return Red Hat products.

* `product_name` - (Optional) A product name to filter the product search on.

## Attributes Reference

* `products` - A list of objects containing information on the products. The attributes of the objects are listed below.

  * `cp_id` -

  * `description` -

  * `gpg_key_id` -

  * `id` -

  * `label` -

  * `last_sync` -

  * `last_sync_words` -

  * `name` -

  * `provider_id` -

  * `repository_count` -
