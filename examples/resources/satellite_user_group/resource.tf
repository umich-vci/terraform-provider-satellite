resource "satellite_user_group" "my_group" {
  name     = "My Group"
  role_ids = [1]
}
