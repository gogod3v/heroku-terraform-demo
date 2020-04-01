variable "heroku_email" {}
variable "heroku_api_key" {}

# Heroku auth & terraform provider version infos
provider "heroku" {
  email   = var.heroku_email
  api_key = var.heroku_api_key
  version = "~> 2.2"
}

# A random ID for this env
# Referencing this will always give you the same value
# until this random value is `terraform destroy`-d or otherwise
# removed from the Terraform state file, which triggers it to be re-generated.
#
# See e.g. the two created heroku apps' name,
# both references this random value so both gets the same value.
resource "random_id" "env_name_id" {
  byte_length = 4
}

# The heroku apps
resource "heroku_app" "redis_sample_app" {
  name   = "demo-${random_id.env_name_id.hex}-redis-sample"
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

## Heroku Addons

### Redis
resource "heroku_addon" "redis_sample_redis_db" {
  app  = heroku_app.redis_sample_app.name
  plan = "heroku-redis:hobby-dev"
}

### Logging
resource "heroku_addon" "redis_sample_logging" {
  app  = heroku_app.redis_sample_app.name
  plan = "logentries:le_tryit"
}


# Source code ( https://www.terraform.io/docs/providers/heroku/r/build.html )
# OPTIONAL - if you prefer to deploy separately to the heroku app, instead of doing
# that as part of the Terraform config, simply leave out this section.
# 
# That said, to be able to quickly create new ad-hoc environments
# the best is if you also include the source code build part.
# 
# The most reliable/reproducible method is using Docker images, where you first build the image,
# then you just reference that image here.
resource "heroku_build" "redis_sample_app" {
  app = heroku_app.redis_sample_app.id

  source = {
    # A local directory, changing its contents will
    # force a new build during `terraform apply`
    path = "./apps/redis-sample"
  }
}

# Heroku Formation: how many dynos to run, and what kind.
resource "heroku_formation" "redis_sample_app" {
  app      = heroku_app.redis_sample_app.id
  type     = "web"
  quantity = 1
  size     = "Free" # "Standard-1x"
  # depends_on : wait until that action is finished - in this case wait until the source build is ready and deployed
  # as you can't scale the web dynos before there's something to scale.
  depends_on = [heroku_build.redis_sample_app]
}


## Second Heroku App

resource "heroku_app" "redis_sample_app2" {
  name   = "demo-${random_id.env_name_id.hex}-redis-sample2"
  region = "us"

  config_vars = {
    # reference the first app's URL
    APP1_URL = heroku_app.redis_sample_app.web_url
  }

  buildpacks = [
    "heroku/go"
  ]
}
