# satellite\_user\_group

Use this resource to create a user group in Red Hat Satellite.

## Example Usage

```hcl
resource "satellite_user_group" "my_group" {
    name = "My Group"
    role_ids = [ 1 ]
}
```

## Argument Reference

* `name` - (Required) A name for the user group.

* `admin` - (Optional) If set to true, then the group will grant administrator privileges.

* `role_ids` - (Optional) A list of IDs of roles to associate with the group.

## Attributes Reference

* `created_at` - A timestamp containing when the user group was created.

* `updated_at` - A timestamp containing when the user group was last changed.

* `roles` - A list of objects containing the roles the associated with the user group.

## Roles

The keys of elements of the `roles` attribute.

* `id` - The ID of a role associated with the user group.

* `name` - The name of a role associated with the user group.

* `description` - A description of a role associated with the user group.

* `origin` - The origin of a role associated with the user group.
