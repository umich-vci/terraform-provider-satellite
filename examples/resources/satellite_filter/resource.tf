resource "satellite_filter" "miscellaneous" {
  role_id          = satellite_role.my_role.id
  resource_type    = ""
  permission_names = [
    "access_dashboard",
    "rh_telemetry_api",
    "rh_telemetry_view",
  ]
  organization_ids = [
    // unlimited == no organization set
  ]
}
