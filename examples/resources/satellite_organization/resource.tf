resource "satellite_organization" "foo" {
  name  = "foo"
  label = "foo"
}

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
