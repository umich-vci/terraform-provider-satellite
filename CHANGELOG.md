## 0.3.0 (not yet released)

FEATURES:

* **New Data Source:** `satellite_location`

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
