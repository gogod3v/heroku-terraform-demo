variable "heroku_email" {}
variable "heroku_api_key" {}

provider "heroku" {
  email   = var.heroku_email
  api_key = var.heroku_api_key
  version = "~> 2.2"
}

resource "heroku_app" "gogod3v-testapp-1" {
  name   = "gogod3v-testapp-1"
  region = "us"

  config_vars = {
    FOOBAR = "baz"
    TEST2  = "test 2"
    TEST3  = "t3"
    TEST4  = "t4"
  }

  buildpacks = [
    "heroku/go"
  ]
}
