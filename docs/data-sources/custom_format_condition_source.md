---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_custom_format_condition_source Data Source - terraform-provider-radarr"
subcategory: "Profiles"
description: |-
  Custom Format Condition Source data source.
  For more information refer to Custom Format Conditions https://wiki.servarr.com/radarr/settings#conditions.
---

# radarr_custom_format_condition_source (Data Source)

<!-- subcategory:Profiles --> Custom Format Condition Source data source.
For more information refer to [Custom Format Conditions](https://wiki.servarr.com/radarr/settings#conditions).

## Example Usage

```terraform
data "radarr_custom_format_condition_source" "example" {
  name     = "Example"
  negate   = false
  required = false
  value    = "7"
}

resource "radarr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.radarr_custom_format_condition_source.example]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Specification name.
- `negate` (Boolean) Negate flag.
- `required` (Boolean) Computed flag.
- `value` (String) Source ID. `0` unknown, `1` cam, `2` telesync, `3` telecine, `4` workprint, `5` dvd, `6` tv, `7` webdl, `8` webrip, `9` bluray.

### Read-Only

- `id` (Number) Custom format condition source ID.
- `implementation` (String) Implementation.

