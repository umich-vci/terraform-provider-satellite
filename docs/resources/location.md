# satellite\_location

Use this resource to create a Red Hat Satellite location.

## Example Usage

```hcl
resource "satellite_location" "Tatooine" {
    name = "Tatooine"
    description = "Well, if there's a bright center of the universe, you're on the planet that it's farthest from."
}

resource "satellite_location" "Mos_Eisley" {
    name = "Mos Eisley"
    description = "You will never find a more wretched hive of scum and villainy."
    parent_id = satellite_location.Tatooine.id
}
```

## Argument Reference

* `name` - (Required) A name for the location.

* `description` - (Optional) A description of the location.

* `parent_id` - (Optional) The ID of a parent for this location. This allows you to nest locations.
  If not set, a top level location is created.
