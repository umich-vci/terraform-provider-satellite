# satellite\_role

Use this resource to create a new role in Red Hat Satellite.

## Example Usage

```hcl
resource "satellite_role" "my_role" {
    name = "My Role"
    organization_ids = [ 10 ]
    description = "Role granting access for someone to do something in one org"
}
```

## Argument Reference

* `name` - (Required) A name for the role.

* `description` - (Optional) A description of the role.

* `location_ids` - (Optional) A list of IDs of locations to associate with the role.

* `organization_ids` - (Optional) A list of IDs of organizations to associate with the role.

## Attributes Reference

* `builtin` - A boolean that indicates if the role is a default/builtin role.

* `cloned_from_id` - If this role was cloned using the API, the ID number of the source role.

* `filters` - A list of filter IDs associated with the role.

* `locations` - A list of objects containing the locations the role applies to.

* `organizations` - A list of objects containing the organizations the role applies to.

* `origin` - 

## Locations

The keys of elements of the `locations` attribute.

* `id` - The ID of a location the filter applies to.

* `name` - The name of a location the filter applies to.

* `title` - The title of a location the filter applies to.

## Organizations

The keys of elements of the `organizations` attribute.

* `id` - The ID of an organization the filter applies to.

* `name` - The name of an organization the filter applies to.

* `title` - The title of an organization the filter applies to.
