data "satellite_organization" "default" {
  search = "name=default"
}

output "default_org_id" {
  value = data.satellite_organization.default.id
}
