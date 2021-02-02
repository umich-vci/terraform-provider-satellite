# satellite\_user\_group

Use this resource to create an external user group in Red Hat Satellite.

## Example Usage

```hcl
resource "satellite_external_user_group" "my_ldap_group" {
    name = "My Group"
    role_ids = [ 1 ]
}
```

## Argument Reference

* `name` - (Required) The name of the external user group.

* `auth_source_id` - (Required) The ID of the authentication source that contains the external user group.

* `user_group_id` - (Required) The ID of the user group that the external user group should be associated with.

## Attributes Reference

* `auth_source_ldap` - A list of objects containing the authentication source the associated with the external user group.

## auth_source_ldap

The keys of elements of the `auth_source_ldap` attribute.

* `id` - The id of the associated authentication source.

* `name` - The name of the associated authentication source.

* `type` - The type of the associated authentication source.
