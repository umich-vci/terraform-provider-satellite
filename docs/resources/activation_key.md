# satellite\_activation\_key

Use this resource to create a Red Hat Satellite Activation Key.

## Example Usage

```hcl
resource "satellite_activation_key" "key" {
    name = "foo"
}
```

## Argument Reference

* `name` - (Required) The name of the activation key.  This is the value of the key that clients
  use to activate.

* `organization_id` - (Required) The ID of the organization to associate with the activation key.

* `content_view_id` - (Optional) The ID of the content view to associate with the activation key.

* `description` - (Optional) A description of the activation key.

* `environment_id` - (Optional) The ID of the environment that contains the `content_view_id`.

* `host_collection_ids` - (Optional) A list of host collection IDs to associate with the activation key.
  Machines activated with the key will be added to these host collections.

* `max_hosts` - (Optional) The maximum number of hosts allowed to use the activation key. Should not be set
  if `unlimited_hosts` is set to `true`.

* `unlimited_hosts` - (Optional) Should an unlimited number of hosts be allowed to use the activation key?
  Defaults to `true`.
