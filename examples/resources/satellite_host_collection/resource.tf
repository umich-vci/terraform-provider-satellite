resource "satellite_host_collection" "host_collection" {
  name            = "My Host Collection"
  organization_id = 10
  description     = "A host collection containing unlimited hosts for something"
}

resource "satellite_host_collection" "limited_host_collection" {
  name            = "My Limited Host Collection"
  organization_id = 10
  description     = "A host collection containing at most 10 hosts for something"
  max_hosts       = 10
  unlimited_hosts = false
}
