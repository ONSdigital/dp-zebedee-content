# dp-zebedee-content

Command line tool generating default content required to run Zebedee CMS.

### Prerequisites
- Go 1.12 +

### Getting started
dp-zebedee-content is a Go Module so needs to be cloned to a directory **outside of your $GOPATH**
```
git clone git@github.com:ONSdigital/dp-zebedee-content.git
```

### Install
```
make install
```

### Run
```bash
./cli generate -c=~/path_to_content_dir -z=~/path_to_zebedee_project_dir
```

See [Flags](#Flags) for further details. 

The `generate` command will:
 - Generate the directory structure required by Zebedee-CMS.
 - Populate the CMS with default content.
 - Generates default user, teams, permissions and service token content.
 - Generates a `run-cms.sh` to running Zebedee locally with typical developer configurations. 

**Note** It's safe to run the `generate` command multiple times. Doing so will overwrite any previously generated 
content and reset the CMS content, users, teams etc. to the default state.  

### Flags
| Flag         | Description                                                                              | Example                                                             |
| ------------ |------------------------------------------------------------------------------------------| ------------------------------------------------------------------- |
| h / help     | Display the help menu.                                                                   |                                                                     |
| c / content  | The output directory the generated content will be written - this can anywhere you like. | `~/Desktop/zebedee-content/generated` (`~` prefix will be expanded) |
| z / zebedee  | The directory of your Zebedee Java project.                                              | `~/IdeaProjects/zebedee`                                            |

Once you have run generator (assuming it has completed successfully) you should now have the required directories, content and configurations to run Zebedee locally.

### CMD configuration
:warning: The `run-cms.sh` script applies default values for the following config properties. When running the script **any existing configuration you have set will take precedence over these defaults**.

```bash
zebedee_root
PORT
ENABLE_DATASET_IMPORT
ENABLE_PERMISSIONS_AUTH
DATASET_API_URL
DATASET_API_AUTH_TOKEN
SERVICE_AUTH_TOKEN
```

Once you have completed the Zebedee set up guide run `./run-cmd.sh` (found in the root of your Zebebee project directory).
 
_NOTE_: You may be required to make it an executable before you can run it.
```
sudo chmod +x run-cms.sh
```

### Help/Issues
If you experience any problems with this tool please speak to a member of the dev team. If you believe there is a defect or issue with the code you can either:
- Raise a [github issue][2].
- Open pull request.

Please be sure to provide a description of the problem and steps to recreate. 

Kind regards
The Dev team  

[1]: https://github.com/kardianos/govendor
[2]: https://github.com/ONSdigital/dp-zebedee-content/issues
