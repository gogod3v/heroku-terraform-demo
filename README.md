# heroku-terraform-demo

Heroku + Terraform demo

## Demo

```shell
git clone git@github.com:gogod3v/heroku-terraform-demo.git
cd heroku-terraform-demo

terraform init
# register & manage Heroku API keys at https://dashboard.heroku.com/account/applications
terraform plan -var-file terraform.tfvars -out=/tmp/plan
# check https://dashboard.heroku.com/apps - no apps yet
terraform apply /tmp/plan
```

Example tfvars file:

```text
heroku_email   = "EMAIL"
heroku_api_key = "HEROKU-API-KEY"
```
