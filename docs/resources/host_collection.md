# satellite\_host\_collection

Use this resource to create a Red Hat Satellite Host Collection.

## Example Usage

```hcl
resource "satellite_host_collection" "host_collection" {
    name            = "My Host Collection"
    organization_id = 10
    description     = "A host collection containing unlimited hosts for something"
}
```

```hcl
resource "satellite_host_collection" "limited_host_collection" {
    name            = "My Limited Host Collection"
    organization_id = 10
    description     = "A host collection containing at most 10 hosts for something"
    max_hosts       = 10
    unlimited_hosts = false
}
```

## Argument Reference

* `name` - (Required) The name of the host collection.

* `organization_id` - (Required) The ID of organization that the host collection should be created in.
  Once set, it cannot be changed without recreating the resource.

* `description` - (Optional) A description of the host collection.

* `max_hosts` - (Optional) The maximum number of hosts allowed to be in the host collection. Should not be set
  if `unlimited_hosts` is set to `true`.

* `unlimited_hosts` - (Optional) A boolean that controls if an unlimited number of members are allowed in the
  host collection. Defaults to `true`.

## Attributes Reference

* `created_at` - A timestamp containing when the host collection was created.

* `updated_at` - A timestamp containing when the host collection was last changed.
