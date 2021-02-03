# satellite\_organization

Use this data source to access information about a Red Hat Satellite organization.

## Example Usage

```hcl
data "satellite_organization" "default" {
    search = "name=default"
}

output "default_org_id" {
    value = data.satellite_organization.default.id
}
```

## Argument Reference

* `search` - (Required) A search filter for the Location search. The search must only return 1 Organization.

## Attributes Reference

* `created_at` - Timestamp of when the organization was created.

* `description` - The description of the organization.

* `label` - The label of the organization.

* `name` - The name of the organization.

* `title` - The title of the organization.

* `updated_at` - Timestamp of when the organization was last updated
