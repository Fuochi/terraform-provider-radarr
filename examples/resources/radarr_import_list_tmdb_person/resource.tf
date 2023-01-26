resource "radarr_import_list_tmdb_person" "example" {
  enabled              = true
  enable_auto          = false
  search_on_add        = false
  root_folder_path     = "/config"
  monitor              = "none"
  minimum_availability = "tba"
  quality_profile_id   = 1
  name                 = "Example"
  person_id            = "11842"
  cast                 = true
  cast_director        = true
  cast_producer        = true
  cast_sound           = true
  cast_writing         = true
}