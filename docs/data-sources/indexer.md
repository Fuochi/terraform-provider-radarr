---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_indexer Data Source - terraform-provider-radarr"
subcategory: "Indexers"
description: |-
  Single Indexer ../resources/indexer.
---

# radarr_indexer (Data Source)

<!-- subcategory:Indexers -->Single [Indexer](../resources/indexer).

## Example Usage

```terraform
data "radarr_indexer" "test" {
  name = "Example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Indexer name.

### Read-Only

- `additional_parameters` (String) Additional parameters.
- `allow_zero_size` (Boolean) Allow zero size files.
- `api_key` (String) API key.
- `api_path` (String) API path.
- `api_user` (String) API User.
- `base_url` (String) Base URL.
- `captcha_token` (String) Captcha token.
- `categories` (Set of Number) Series list.
- `codecs` (Set of Number) Codecs.
- `config_contract` (String) Indexer configuration template.
- `cookie` (String) Cookie.
- `delay` (Number) Delay before grabbing.
- `download_client_id` (Number) Download client ID.
- `enable_automatic_search` (Boolean) Enable automatic search flag.
- `enable_interactive_search` (Boolean) Enable interactive search flag.
- `enable_rss` (Boolean) Enable RSS flag.
- `id` (Number) Indexer ID.
- `implementation` (String) Indexer implementation name.
- `mediums` (Set of Number) Mediumd.
- `minimum_seeders` (Number) Minimum seeders.
- `multi_languages` (Set of Number) Language list.
- `passkey` (String) Passkey.
- `priority` (Number) Priority.
- `protocol` (String) Protocol. Valid values are 'usenet' and 'torrent'.
- `ranked_only` (Boolean) Allow ranked only.
- `remove_year` (Boolean) Remove year.
- `required_flags` (Set of Number) Computed flags.
- `seed_ratio` (Number) Seed ratio.
- `seed_time` (Number) Seed time.
- `tags` (Set of Number) List of associated tags.
- `username` (String) Username.

