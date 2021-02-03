# satellite\_location

Use this data source to access information about a Red Hat Satellite location.

## Example Usage

```hcl
resource "satellite_location" "Tatooine" {
    search = "name=Tatooine"
}
```

## Argument Reference

* `search` - (Required) A search filter for the Location search. The search must only return 1 Location=.

## Attributes Reference

* `created_at` - Timestamp of when the location was created.

* `description` - A description of the location.

* `name` - The name of the location.

* `parent_id` - The ID of the parent for this location.  If not set, the location is a top level location.

* `parent_name` - The name of the parent for this location.  If not set, the location is a top level location.

* `title` - The title of the location.

* `updated_at` - Timestamp of when the location was last updated
