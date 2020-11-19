# satellite\_subscription\_manifest

Use this resource attach a subscription manifest to a Red Hat Satellite organization.

## Example Usage

```hcl

```

## Argument Reference

* `organization_id` - (Required) The organization ID you want to attach the manifest to.

* `manifest` - (Required) A Base64 encoded string of a manifest zip file downloaded from
  Red Hat Subscription Management. Most easily used in conjunction with [`rhsm_allocation_manifest` resource from the RHSM provider](https://registry.terraform.io/providers/umich-vci/rhsm/latest/docs/resources/allocation_manifest).

## Attributes Reference

* `history` - A list of objects containing information on operations peformed on the manifest.

  * `created` - 

  * `id` - 

  * `status` - 

  * `status_message` - 
