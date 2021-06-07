resource "satellite_location" "Tatooine" {
  name        = "Tatooine"
  description = "Well, if there's a bright center of the universe, you're on the planet that it's farthest from."
}

resource "satellite_location" "Mos_Eisley" {
  name        = "Mos Eisley"
  description = "You will never find a more wretched hive of scum and villainy."
  parent_id   = satellite_location.Tatooine.id
}
