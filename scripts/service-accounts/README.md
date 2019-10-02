# Service Account Generation

Set up script for generating Zebedee service accounts.

## Setting up from scratch

1. SSH into the environment

2. Run a golang container with a volume that maps the zebedee directory on the publishing box:
```
sudo docker run -i -t --name service-accounts \
   --userns=host \
   -v <CONTENT_DIR>:/zebedee:rw \
   golang /bin/bash
```

3. Update the container and install vim (always useful)
```
apt-get update && apt-get install vim
```

4. Clone the script - script uses Go Modules so should be placed outside of the $GOPATH
```
git clone git@github.com:ONSdigital/dp-zebedee-content.git
```

5. Move into the scripts dir
```
cd dp-zebedee-content/scripts/service-accounts
```

6. Build the Go binary
```
go build -o generator
```

7. Run the script: Assuming everything is good the script will generate each of the required service accounts. The 
script is non destructive so if a service account already exists it will not be overwritten.

```
export HUMAN_LOG=true

./generator -dir="/service"
```

### Service Account files.
Service account files is a simple `.json` file containing an ID field - the service name. The file name is a UUID which 
is the **service token value**.

#### Example Florence service account

`1L1YlW7aA2hMMGetIbRv3IE3jIgdqYaXkeF8NTXYyZUh3XyvbHh5tUeYnSSCw0x9.json`
```json
{
  "id": "florence"
}
```

If Florence makes a request to an API requiring service authentication it will send an auth request header with value 
`1L1YlW7aA2hMMGetIbRv3IE3jIgdqYaXkeF8NTXYyZUh3XyvbHh5tUeYnSSCw0x9` to identify its self. 
