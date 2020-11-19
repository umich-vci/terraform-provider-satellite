# satellite\_permissions

Use this data source to access information about available Red Hat Satellite permissions.

## Example Usage

```hcl
data "satellite_permissions" "host_collection" {
    resource_type = "Katello::HostCollection"
}
```

## Argument Reference

* `name` - (Optional) A permission name to filter the permission search on.

* `resource_type` - (Optional) A resource type to filter the permission search on.

## Attributes Reference

* `permissions`

  * `id` - The ID of the permission.

  * `name` - The name of the permission.

  * `resource_type` - The resource type the permission applies to.
