# Terraform Provider for LogDNA

[![Coverage Status](https://coveralls.io/repos/github/logdna/terraform-provider-logdna/badge.svg?branch=main)](https://coveralls.io/github/logdna/terraform-provider-logdna?branch=main)
[![Public Beta](https://img.shields.io/badge/-Public%20Beta-404346?style=flat)](#)

[LogDNA](https://logdna.com) is a centralized log management platform. The LogDNA Terraform Provider allows organizations to manage certain LogDNA resources (alerts, views, etc) programmatically via Terraform.

The [official docs for the LogDNA terraform provider](https://registry.terraform.io/providers/logdna/logdna/latest/docs) can be found in the Terraform registry.

Using the `logdna_view` resource, a user can create a View with a provided `name`, `query`, `hosts`, `categories`, `tags`, `email_channel`, `pagerduty_channel`, and `webhook_channel`, delete a View with a given `viewid` or update a View using the `viewid` and `name`.

Using the `logdna_alert` resource, a user can create a Preset Alert with a provided `name`, `email_channel`, `pagerduty_channel` and `webhook_channel`, delete a Preset Alert with a given `presetid` or update a Preset Alert using the `presetid` and `name`.

Run `terraform init`, `terraform plan`, and `terraform apply`, refresh your browser and then navigate to the UI to see your updates.

In addition to the examples provided below, sample .tf files can be found [here](https://github.com/logdna/terraform-provider-logdna/tree/main/examples).

## Example Terraform Configuration for Preset Alerts
```
provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_alert" "my_alert" {
  name  = "Email and PagerDuty Preset Alert"
  email_channel {    
    emails         = ["test@logdna.com"]                 
    operator       = "absence"
    timezone       = "Pacific/Samoa"
    triggerlimit   = 15                  
  }

  pagerduty_channel {
    immediate       = "false"
    key             = "Your PagerDuty API key goes here"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
}
```

## Example Terraform Configuration for Views
```
provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_view" "my_view" {
  apps     = ["app1", "app2"]
  categories = ["Demo1", "Demo2"]
  hosts    = ["host1", "host2"]
  levels   = ["fatal", "critical"]
  name     = "Email PagerDuty and Webhook View-specific Alerts"
  query    = "test"
  tags     = ["tag1", "tag2"]

  email_channel {
    emails          = ["test@logdna.com"]
    immediate       = "false"
    operator        = "absence"
    terminal        = "true"
    timezone        = "Pacific/Samoa"
    triggerinterval = "15m"
    triggerlimit    = 15
  }

  pagerduty_channel {
    immediate       = "false"
    key             = "Your PagerDuty API key goes here"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }

  webhook_channel {
    bodytemplate = jsonencode({
      hello = "test1"
      test  = "test2"
    })
    headers = {
      hello = "test3"
      test  = "test2"
    }
    immediate       = "false"
    method          = "post"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
    url             = "https://yourwebhook/endpoint"
  }
}
```

## Development

### Prerequisites

In order to test the provider you will need to have a `SERVICE_KEY` environment variable
exported in your shell. Your service key can be generated or retrieved from your LogDNA
account at **Settings > Organization > API Keys**.

To run the archiving tests, you will need to have `S3_BUCKET`, `GCS_BUCKET`, `GCS_PROJECTID` environment
variables exported in your shell. These should be valid settings to create S3 and GCS archiving
configurations.

### Local Test, Build, & Install

During development, the full test suite can be run with:

```sh
make test-local
```

The provider can be built and installed locally in `$HOME` by running:

```sh
make install-local
```

After running `make install-local`, you will be able to reference the snapshot build
version in any local Terraform configuration. For example:

```hcl
terraform {
  required_providers {
    logdna = {
      source = "logdna.com/logdna/logdna"
      version = "1.2.0-pre-SNAPSHOT-b9faaaa"
    }
  }
}
```

### Docker

The included tooling can be used to test and build the provider inside a Docker build
environment, without installing any dependencies locally. 

You will need an ascii-armored GPG key in the root of the project at `./gpgkey.asc` for
signing test builds. If you do not already have a personal GPG key, you can generate one
by [following this guide](https://docs.github.com/en/github/authenticating-to-github/managing-commit-signature-verification/generating-a-new-gpg-key).

Export your key to the root of this repository:

```sh
gpg --armor --export <ID> > ./gpgkey.asc
```

**NOTE:** This is only for local testing via Docker. The release process should
only be run from CI which will sign binaries with the proper production key.

The following build targets are useful for running locally within Docker:

```sh
make test         # run tests
make build        # build the provider for your host OS/ARCH
make test-release # build for all supported targets
```

### Tagging and Release

This project uses [`svu`](https://github.com/caarlos0/svu) for parsing
[conventional commit messages](https://github.com/caarlos0/svu#commit-messages-vs-what-they-do)
and determining the next version/tag. Once code is merged into `main`, a release will be
built and created as a draft in GitHub. The draft release must be published manually in
GitHub for it to be pulled in by the Terraform Registry.
