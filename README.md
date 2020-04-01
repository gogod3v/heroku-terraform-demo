# heroku-terraform-demo

Heroku + Terraform demo

## Demo

```shell
git clone git@github.com:gogod3v/heroku-terraform-demo.git
cd heroku-terraform-demo

# install the required terraform version, e.g. via tfenv:
tfenv install 0.12.18

terraform init
# register & manage Heroku API keys at https://dashboard.heroku.com/account/applications
terraform plan -var-file terraform.tfvars -out=/tmp/plan
# check https://dashboard.heroku.com/apps - no apps yet


terraform apply /tmp/plan
# check https://dashboard.heroku.com/apps again - app was created

# go through the heroku.tf config and explain what's what

# run terraform plan again - no changes
terraform plan -var-file terraform.tfvars -out=/tmp/plan

# change something in heroku.tf, e.g. an env var of one of the heroku apps, then terraform plan again
terraform plan -var-file terraform.tfvars -out=/tmp/plan
# there's a change! Apply it!
terraform apply /tmp/plan
# Check that it's changed on Heroku.

# cleanup:
terraform destroy
```

Example tfvars file:

```text
heroku_email   = "EMAIL"
heroku_api_key = "HEROKU-API-KEY"
```

## Sources

- https://devcenter.heroku.com/articles/using-terraform-with-heroku
- https://www.terraform.io/docs/providers/heroku/index.html
    - https://www.terraform.io/docs/providers/heroku/r/build.html
