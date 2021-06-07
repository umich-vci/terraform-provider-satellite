resource "satellite_role" "my_role" {
  name             = "My Role"
  organization_ids = [10]
  description      = "Role granting access for someone to do something in one org"
}
