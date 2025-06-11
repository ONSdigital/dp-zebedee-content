# dp-zebedee-content

Test content for local development of the dissemination "legacy core" stack. Specifically, this repo generates the directory structure and test content required by zebedee (aka zebedee cms), zebedee reader and the-train.

## Getting started

It is strongly recommended that you use a dp-compose stack for your local environment rather than creating one manually. These stacks will complete the following automatically through the stack's provisioning and volume mounts.

If you choose to manually set up an environment then:

1. Clone this repo:

   ```shell
   git clone git@github.com:ONSdigital/dp-zebedee-content
   ```

2. Initialise the content

   ```shell
   cd dp-zebedee-content
   make init
   ```

3. Set environment variables to configure the appropriate directories for the apps:

   ```shell
   export DP_ZEBEDEE_CONTENT=<PATH TO THIS REPO>

   # zebedee cms
   export zebedee_root=${DP_ZEBEDEE_CONTENT}/generated/publishing

   # zebedee reader
   export content_dir=${DP_ZEBEDEE_CONTENT}/generated/web/site

   # the-train
   export TRANSACTION_STORE=${DP_ZEBEDEE_CONTENT}/generated/web/transactions
   export WEBSITE=${DP_ZEBEDEE_CONTENT}/generated/web/site
   ```

### Resetting the content

You can delete and recreate your directories and content by re-running:

```shell
make init
```

## Adding content

Additional test content can be added by:

1. On a feature branch, copy all files from the directory corresponding to a page on the site to the corresponding location under `src/content`

2. Re-initialise to test the new content

   ```shell
   make init
   ```

3. If the content works as expected then open a PR for the change

## Licence

Copyright (c) Crown Copyright ([Office for National Statistics](https://www.ons.gov.uk))

Test content provided by this repo under contains public sector information sourced from [ONS.GOV.UK](https://www.ons.gov.uk/) and licensed under the [Open Government Licence v3.0](https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/).

Unless stated otherwise, this repo is released under MIT license, see [LICENCE](LICENCE) for details.
