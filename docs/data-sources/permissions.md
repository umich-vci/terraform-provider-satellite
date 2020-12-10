# satellite\_permissions

Use this data source to access information about available Red Hat Satellite permissions.

## Example Usage

```hcl
data "satellite_permissions" "host_collection" {
    search = "resource_type=Katello::HostCollection"
}
```

## Argument Reference

* `search` - (Optional) A search filter for the permission search. If not specified all permissions are returned.

## Attributes Reference

* `permissions`

  * `id` - The ID of the permission.

  * `name` - The name of the permission.

  * `resource_type` - The resource type the permission applies to.
