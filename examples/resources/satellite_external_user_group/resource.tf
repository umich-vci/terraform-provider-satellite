resource "satellite_external_user_group" "my_ldap_group" {
    name = "My Group"
    role_ids = [ 1 ]
}
