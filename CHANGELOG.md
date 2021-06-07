## 0.4.0 (unreleased)

ENHANCEMENTS:

* Updated [terraform-plugin-sdk](https://github.com/umich-vci/gosatellite) to 2.6.1.

* Reworked code to model the approach in
  [terraform-provider-scaffoling](https://github.com/hashicorp/terraform-provider-scaffolding).

* Added descriptions to resources and data sources to allow for usage in documentation
  generation and in the language server.

## 0.3.3 (June 3, 2021)

IMPROVEMENTS:

* resource/satellite_filter: Support `InsightsHit` filter type added in Satellite 6.9.

## 0.3.2 (February 18, 2021)

BUG FIXES:

* resource/satellite_external_user_group: Fixed `auth_source_ldap` computed value not being set properly.

## 0.3.1 (February 5, 2021)

BUG FIXES:

* resource/satellite_activation_key: Fixed setting optional values at creation time.

## 0.3.0 (February 4, 2021)

BREAKING CHANGES:

* datasource/satellite_location: `name` changed to computed value and `search` added as an argument.
* datasource/satellite_organization: `id` changed to computed value and `search` added as an argument.

FEATURES:

* **New Data Source:** `satellite_auth_source_ldap`
* **New Data Source:** `satellite_location`
* **New Resource:** `satellite_external_user_group`
* **New Resource:** `satellite_user_group`

## 0.2.0 (December 10, 2020)

BREAKING CHANGES:

* datasource/satellite_permissions: Removed `name` and `resource_type` arguments.
  Added `search` argument which is compatible with older and newer versions of Satellite.

IMPROVEMENTS:

* resource/satellite_filter: Resolved issue with permissions name validation not working correctly
  on Satellite 6.8.

* Updated [gosatellite](https://github.com/umich-vci/gosatellite) to 20201210181146-c8a049d1e6ab
  which removes some parameters incompatible with Satellite 6.8 from permissions list API call.

* Upgraded from Terraform SDK 2.0.3 to 2.3.0

## 0.1.0 (November 19, 2020)

Initial release of provider to Terraform Registry.
