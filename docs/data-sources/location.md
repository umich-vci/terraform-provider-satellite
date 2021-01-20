# satellite\_location

Use this data source to access information about a Red Hat Satellite location.

## Example Usage

```hcl
resource "satellite_location" "Tatooine" {
    name = "Tatooine"
}
```

## Attributes Reference

* `description` - A description of the location.

* `parent_id` - The ID of a parent for this location. This allows you to nest locations.
  If not set, the resource is a top level location.
