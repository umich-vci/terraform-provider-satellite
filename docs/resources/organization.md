# satellite\_organization

Use this resource to create a Red Hat Satellite organization.

## Example Usage

```hcl
resource "satellite_organization" "foo" {
    name = "foo"
    label = "foo"
}
```

```hcl
resource "satellite_organization" "foo" {
    name = "foo"

    lifecycle {
      ignore_changes = [
        // label defaults to the same as name in the API
        // However since we didn't explicitly set it it needs to be ignored
        // so that the resource isn't constantly destroyed and recreated.
        label,
      ]
    }
}
```

## Argument Reference

* `name` - (Required) The name of the organization.

* `description` - (Optional) A description of the organization.

* `label` - (Optional) The label of the organization. If not set, Satellite will
  use the `name` as the label.  This field can only be set at creation time. If not
  being set explicitly you will probably want to use `ignore_changes` on this
  in the lifecycle block.

## Attributes Reference

* `hosts_count` - A count of how many hosts are registered to the organization.

* `title` - The title of the organization.
