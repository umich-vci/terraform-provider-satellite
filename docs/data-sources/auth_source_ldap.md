# satellite\_auth_source\_ldap

Use this data source to access information about a Red Hat Satellite LDAP Authentication Source.

## Example Usage

```hcl
resource "satellite_auth_source_ldap" "ActiveDirectory" {
    search = "name=AD"
}
```

## Argument Reference

* `search` - (Required) A search filter for the LDAP Authentication Source search. The search must only return 1 authentication source.

## Attributes Reference

* `permissions`

  * `id` - The ID of the LDAP Authentication Source

  * `account` - The DN of the LDAP Bind Account

  * `attr_login` - The LDAP attribute that maps to username

  * `attr_firstname` - The LDAP attribute that maps to first name

  * `attr_lastname` - The LDAP attribute that maps to last name

  * `attr_mail` - The LDAP attribute that maps to email address

  * `attr_photo` - The LDAP attribute that maps to a photo

  * `base_dn` - The base DN from which LDAP searches will be performed

  * `created_at` - Timestamp of when the LDAP authentication source was created.

  * `groups_base` - The base DN from which LDAP searches for groups will be performed

  * `host` - The hostname of the LDAP server

  * `ldap_filter` -

  * `name` - The name of the LDAP authentication source

  * `onthefly_register` -

  * `server_type` - The type of the LDAP server

  * `tls` - Is TLS enabled for the LDAP server?

  * `type` -

  * `use_netgroups` -

  * `updated_at` - Timestamp of when the LDAP authentication source was last updated

  * `usergroup_sync` -
