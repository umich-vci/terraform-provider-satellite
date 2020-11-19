# satellite\_organization

Use this data source to access information about a Red Hat Satellite organization.

## Example Usage

```hcl
data "satellite_organization" "default" {
    id = 1
}

output "default_org_name" {
    value = data.satellite_organization.default.name
}
```

## Argument Reference

* `id` - (Required) The id of the organiztion to look up.

## Attributes Reference

* `description` - The description of the organization.

* `hosts_count` - A count of how many hosts are registered to the organization.

* `label` - The label of the organization.

* `name` - The name of the organization.

* `title` - The title of the organization.
