# dp-zebedee-content

Command line tool generating default content required to run Zebedee CMS.

![Alt text](preview.png?raw=true "Optional Title")

## Prerequisites

- Go version >= `1.15`
- Access to the AWS dev account
- Your aws config/credentials files will need entries for the `[default]` profile:

Example: `.aws/config`

```ini
  [default]
  region=eu-west-1
  ```

Example: `.aws/credentials`

```ini
...
[default]
aws_access_key_id=...
aws_secret_access_key=...

[development]
aws_access_key_id=....
aws_secret_access_key=...
region=eu-west-1
...
```

## Getting started

dp-zebedee-content is a Go Module so needs to be cloned to a directory **outside of your $GOPATH**

```shell
git clone git@github.com:ONSdigital/dp-zebedee-content.git
```

### Install

```shell
make install
```

### Run

```shell
dp-zebedee-content generate -c=~/path_where_you_want_the_content_to_be_generated
```

See [Flags](#Flags) for further details.

The `generate` command will:

- Generate the directory structure required by Zebedee-CMS.
- Populate the CMS with default content.
- Generates default user, teams, permissions and service token content.

**Note** It's safe to run the `generate` command multiple times. Doing so will overwrite any previously generated
content and reset the CMS content, users, teams etc. to the default state.

### Flags

| Flag        | Description                                                                                                                                                 | Example                                                              |
| ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| h / help    | Display the help menu.                                                                                                                                      |                                                                      |
| c / content | The output directory the generated content will be written - this can be anywhere you like. If unset, the `zebedee_root` env var will be attempted instead. | `~/Desktop/zebedee-content/generated` (`~` prefix will be expanded). |

Once you have run generator add the output `zebedee_root` and `SERVICE_AUTH_TOKEN` values to your ENV vars.
You should now have the required directories, content and configurations to run a local copy of Zebedee CMS.

## Help/Issues

If you experience any problems with this tool please speak to a member of the dev team. If you believe there is a defect or issue with the code you can either:

- Raise a [github issue][2].
- Open pull request.

Please be sure to provide a description of the problem and steps to recreate.

Kind regards

The Dev team

[1]: https://github.com/kardianos/govendor
[2]: https://github.com/ONSdigital/dp-zebedee-content/issues
