---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_custom_format_condition_size Data Source - terraform-provider-radarr"
subcategory: "Profiles"
description: |-
  Custom Format Condition Size data source.
  For more information refer to Custom Format Conditions https://wiki.servarr.com/radarr/settings#conditions.
---

# radarr_custom_format_condition_size (Data Source)

<!-- subcategory:Profiles --> Custom Format Condition Size data source.
For more information refer to [Custom Format Conditions](https://wiki.servarr.com/radarr/settings#conditions).

## Example Usage

```terraform
data "radarr_custom_format_condition_size" "example" {
  name     = "Example"
  negate   = false
  required = false
  min      = 5
  max      = 50
}

resource "radarr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.radarr_custom_format_condition_size.example]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `max` (Number) Max size in GB.
- `min` (Number) Min size in GB.
- `name` (String) Specification name.
- `negate` (Boolean) Negate flag.
- `required` (Boolean) Computed flag.

### Read-Only

- `id` (Number) Custom format condition size ID.
- `implementation` (String) Implementation.

