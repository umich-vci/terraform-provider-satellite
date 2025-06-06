---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "satellite_lifecycle_environment Data Source - terraform-provider-satellite"
subcategory: ""
description: |-
  Data source to access information about a Red Hat Satellite Lifecycle Environment.
---

# satellite_lifecycle_environment (Data Source)

Data source to access information about a Red Hat Satellite Lifecycle Environment.

## Example Usage

```terraform
resource "satellite_lifecycle_environment" "library" {
  name            = "Library"
  organization_id = 5
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) The name of the Lifecycle Environment.
- `organization_id` (Number) The ID of the organization that contains the Lifecycle Environment.
- `search` (String) A search filter for the Lifecycle Environment search. The search must only return 1 Lifecycle Environment.

### Read-Only

- `counts` (Map of Number) A map of various counts of attributes of the Lifecycle Environment.
- `created_at` (String) Timestamp of when the Lifecycle Environment was created.
- `description` (String) A description of the Lifecycle Environment.
- `id` (String) The ID of this resource.
- `label` (String) A label for the Lifecycle Environment.
- `library` (Boolean) Is the Lifecycle Environment a base Library?
- `organization` (Map of String) The organization that contains the Lifecycle Environment.
- `permissions` (Map of Boolean) The current Satellite user's permissions for the Lifecycle Environment.
- `prior` (Map of String) The Lifecycle Environment directly before this one.
- `registry_name_pattern` (String) TODO
- `registry_unauthenticated_pull` (Boolean) TODO
- `successor` (Map of String) The Lifecycle Environment directly after this one.
- `updated_at` (String) Timestamp of when the Lifecycle Environment was last updated.
