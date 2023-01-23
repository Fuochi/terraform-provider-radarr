package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMovieResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccMovieResourceConfig("test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_movie.test", "path", "/config/test"),
					resource.TestCheckResourceAttrSet("radarr_movie.test", "id"),
					resource.TestCheckResourceAttr("radarr_movie.test", "original_title", "The Matrix"),
					resource.TestCheckResourceAttr("radarr_movie.test", "status", "released"),
					resource.TestCheckResourceAttr("radarr_movie.test", "monitored", "false"),
					resource.TestCheckResourceAttr("radarr_movie.test", "year", "1999"),
					resource.TestCheckResourceAttr("radarr_movie.test", "minimum_availability", "inCinemas"),
					resource.TestCheckResourceAttr("radarr_movie.test", "imdb_id", "tt0133093"),
					resource.TestCheckResourceAttr("radarr_movie.test", "is_available", "true"),
					resource.TestCheckResourceAttr("radarr_movie.test", "original_language.id", "1"),
					resource.TestCheckResourceAttr("radarr_movie.test", "original_language.name", "English"),
					resource.TestCheckResourceAttr("radarr_movie.test", "genres.0", "Action"),
				),
			},
			// Update and Read testing
			{
				Config: testAccMovieResourceConfig("test123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_movie.test", "path", "/config/test123"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_movie.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMovieResourceConfig(path string) string {
	return fmt.Sprintf(`
		resource "radarr_movie" "test" {
			monitored = false
			title = "The Matrix"
			path = "/config/%s"
			quality_profile_id = 1
			tmdb_id = 603

			minimum_availability = "inCinemas"
		}
	`, path)
}
