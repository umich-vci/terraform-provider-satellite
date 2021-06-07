resource "satellite_subscription_manifest" "manifest" {
  organization_id = satellite_organization.org.id
  manifest        = rhsm_allocation_manifest.manifest.manifest
}
