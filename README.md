# dp-zebedee-content

Command line helper tool for developer Zebedee-CMS set up:
 - Generates the required file system structure.
 - Populates the CMS with basic content
 - Generates a service account (required if running with the _CMD_ feature enabled).
 - Generates a `run-cms.sh` for running Zebedee locally with typical developer configurations.

### Prerequisites
- Go 1.10.2
- [Govendor][1] 

### Getting started
```
go get github.com/ONSdigital/dp-zebedee-content
```

There are 2 commands provided by this repository: content and service-accounts

To build content:
```
cd dp-zebedee-content/cmd/content
go build -o ../../zebContent
```

To build service-accounts:
```
cd dp-zebedee-content/cmd/service-accounts
go build -o ../../service-accounts
```

This will result in executable, clearly named files in the root of this repository.

### Run it
```
./zebContent -content_dir=[CONTENT_DIR] -project_dir=[PROJECT_DIR] -enable_cmd[true/false]
```

### Flags
| Flag         | Description                                                                   | Example                                                |
| ------------ |------------------------------------------------------------------------------ | ------------------------------------------------------ |
| h / help     | Display the help menu.                                                        |                                                        |
| content_dir  | The directory to build the Zebedee directory structure and unpack the content | `/Users/RickSanchez/Desktop/zebedee-content/generated` |
| project_dir  | The root directory of your Zebedee java project.                              | `/Users/RickSanchez/IdeaProjects/zebedee`              |
| enable_cmd   | Enable or disable the _CMD_ feature.                                          |                                                        |


Once you have run generator (assuming it has completed successfully) you should now have the required directories, content and configurations to run Zebedee locally.

If you are running Zebedee with the _CMD_ feature enabled you can find the _CMD_ configuration values, auth tokens & service accounts ID's under the CMD section of the `./run-cmd.sh`.

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
