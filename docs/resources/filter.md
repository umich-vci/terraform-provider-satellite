# satellite\_filter

Use this resource to create a new permission filter for a role in Red Hat Satellite.

## Example Usage

```hcl
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
```

## Argument Reference

* `permission_names` - (Required) A list of permission names that should be enabled in the filter.
  The permission names must be valid for the role specified in `resource_type`.

* `role_id` - (Required) The ID of the role that the filter should be created under.

* `resource_type` - (Required) The resource type of the filter.  Once this is set, it cannot be
  changed without recreating the filter.  As of Satellite 6.7, valid resource types are:
  `""` (empty string - this corresponds with "Miscellaneous" in the GUI), `AnsibleRole`,
  `AnsibleVariable`, `Architecture`, `Audit`, `AuthSource`, `Bookmark`, `ComputeProfile`,
  `ComputeResource`, `ConfigGroup`, `ConfigReport`, `DiscoveryRule`, `Domain`, `Environment`,
  `ExternalUsergroup`, `FactValue`, `Filter`, `ForemanOpenscap::ArfReport`,
  `ForemanOpenscap::Policy`, `ForemanOpenscap::ScapContent`, `ForemanOpenscap::TailoringFile`,
  `ForemanTasks::RecurringLogic`, `ForemanTasks::Task`, `ForemanVirtWhoConfigure::Config`, `Host`,
  `HostClass`, `Hostgroup`, `HttpProxy`, `Image`, `JobInvocation`, `JobTemplate`,
  `Katello::ActivationKey`, `Katello::ContentView`, `Katello::GpgKey`, `Katello::HostCollection`,
  `Katello::KTEnvironment`, `Katello::Product`, `Katello::Subscription`, `Katello::SyncPlan`,
  `KeyPair`, `Location`, `MailNotification`, `Medium`, `Model`, `Operatingsystem`, `Organization`,
  `Parameter`, `PersonalAccessToken`, `ProvisioningTemplate`, `Ptable`, `Puppetclass`,
  `PuppetclassLookupKey`, `Realm`, `RemoteExecutionFeature`, `Report`, `ReportTemplate`, `Role`,
  `Setting`, `SmartProxy`, `SshKey`, `Subnet`, `Template`, `TemplateInvocation`, `Trend`, `User`,
  `Usergroup`, or `VariableLookupKey`.

* `location_ids` - (Optional) A list of IDs of locations to associate with the filter. Unless `override` is
  set to `true` this should generally contain the `location_ids` that the parent role is associated with. It may
  also need to be set to an empty list if you desire the permission to be `unlimited`.

* `organization_ids` - (Optional) A list of IDs of organizations to associate with the filter. Unless `override` is
  set to `true` this should generally contain the `organization_ids` that the parent role is associated with. It may also need to be set to an empty list if you desire the permission to be `unlimited`.

* `override` - (Optional) When set to true, you can specify `location_ids` and `organization_ids` to allow
  the role to access the `resource_type` in the specified locations and organizations.

* `search` - (Optional)  If this is not set, then the filter will apply to all objects of the specified resource type.
  This means the value of `unlimited` will be true.  You can specify a search which can be used to limit the resources
  that the permission applies to. This will result in the value of `unlimited` being false. For more information see
  the [Red Hat documentation](https://access.redhat.com/documentation/en-us/red_hat_satellite/6.8/html/administering_red_hat_satellite/chap-Red_Hat_Satellite-Administering_Red_Hat_Satellite-Users_and_Roles#sect-Red_Hat_Satellite-Administering_Red_Hat_Satellite-Users_and_Roles-Granular_Permission_Filtering).

## Attributes Reference

* `created_at` - A timestamp of when the filter was created.

* `locations` - A list of objects containing the locations the filter applies to.

* `organizations` - A list of objects containing the organizations the filter applies to.

* `permissions` - A list of objects containing the permissions enabled in the filter.

* `permission_ids` - A list of permission IDs that match the list of permissions supplied in `permission_names`.

* `role` - An object containing information about the role the filter is associated with.

* `unlimited` - A boolean that indicates if a filter applies to all resources of the `resource_type` or just
  a subset of resources specified in `search`.

* `updated_at` - A timestamp of when the filter was last updated.

## Locations

The keys of elements of the `locations` attribute.

* `id` - The ID of a location the filter applies to.

* `name` - The name of a location the filter applies to.

* `title` - The title of a location the filter applies to.

## Organizations

The keys of elements of the `organizations` attribute.

* `id` - The ID of an organization the filter applies to.

* `name` - The name of an organization the filter applies to.

* `title` - The title of an organization the filter applies to.

## Permissions

The keys of elements of the `permissions` attribute.

* `id` - The ID of a permission enabled in the filter.

* `name` - The name of the permission enabled in the filter.

* `resource_type` - The resource type the permission applies to.

## Role

The keys of the `role` object.

* `description` - The description of the role the filter is attached to.

* `id` - The ID of the role the filter is attached to.

* `name` - The name of the role the filter is attached to.
