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
go build -o zebContent
```

### Run it
```
./zebContent -content_dir=[CONTENT_DIR] -project_dir=[PROJECT_DIR] -enable_cmd[true/false]
```

| Flag          | Description                                                                                    | Example                                                  |
| ------------- |----------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| -h / -help    | Display the help menu.                                                                         |                                                          |
| -content_dir  | The directory in which to build the Zebedee directory structure and unpack the default content | `/Users/RickSanchez/Desktop/zebedee-content/generated`   |
| -project_dir  | The root directory of your Zebedee java project                                                | `/Users/RickSanchez/IdeaProjects/zebedee`                |
| -enable_cmd   | If `enabled` the generated `run-cms.sh` script will have the CMD feature enabled               |                                                          |

Once you have run generator (assuming it has completed successfully) you should have the required directories, content and config to run Zebedee locally.
To do so make sure you have completed the Zebedee set up guide and simply  run `./run-cmd.sh` (found in the root of your Zebebee project directory).
 
_NOTE_: You may be required to make it an executable before you can run it.
```
sudo chmod +x run-cms.sh
```

[1]: https://github.com/kardianos/govendor