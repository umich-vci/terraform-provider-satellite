# Satellite Provider

 The Satellite provider is used to interact with Red Hat Satellite 6.x.

## Example Usage

```hcl
// Configure the Satellite Provider
provider "satellite" {
    username = "username"
    password = "password123"
    satellite_host = satellite.example.com
}

// Create a new Organization
resource "satellite_organization" "foo" {
  name = "foo"
  label = "foo"
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `username` - (Optional) This is the username to use to access the Red Hat Satellite server.
  This must be provided in the config or in the environment variable `SATELLITE_USERNAME`.

* `password` - (Optional) This is the password to use to access the Red Hat Satellite server.
  This must be provided in the config or in the environment variable `SATELLITE_PASSWORD`.

* `satellite_host` - (Optional) This is the hostname or IP address of the Red Hat Satellite
  server. This must be provided in the config or in the environment variable
  `SATELLITE_HOST`.

* `ssl_verify` - (Optional) Should we validate the SSL certificate presented by the Satellite
  server. Defaults to `true`.
