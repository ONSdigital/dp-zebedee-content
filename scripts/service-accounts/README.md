# Service Account Generation

Set up script for generating Zebedee service accounts.

### Service Accounts.
Service accounts are how our API services authenticate inbound and "sign" outbound requests and identify the caller. 
Service account files are simple `.json` files containing an ID field - the service name (confusingly). The file name is
 a UUID which is the **service token value**.

#### Example Florence service account

`1L1YlW7aA2hMMGetIbRv3IE3jIgdqYaXkeF8NTXYyZUh3XyvbHh5tUeYnSSCw0x9.json`
```json
{
  "id": "florence"
}
```
If Florence makes a request to an API requiring service authentication it will send an auth request header with value 
`1L1YlW7aA2hMMGetIbRv3IE3jIgdqYaXkeF8NTXYyZUh3XyvbHh5tUeYnSSCw0x9` to identify itself. 

## Generating Service Accounts

### On an environment 

1. SSH onto the Zebedee publishing-mount.

2. Run a golang docker container with a volume that maps the zebedee directory on the publishing box:
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

4. Clone the script - the script uses Go Modules so should be placed outside of the $GOPATH
**Note:** When cloning the repo make sure you use `https` instead of `ssh`.
```
git clone https://github.com/ONSdigital/dp-zebedee-content.git
```

5. Move into the scripts dir
```
cd dp-zebedee-content/scripts/service-accounts
```

6. Build the generator binary
```
go build -o generator
```

7. Run the script: Assuming everything is good the script will generate each of the required service accounts. The 
script is non destructive so if a service account already exists it will not be overwritten. 
 
```
export HUMAN_LOG=true

./generator -dir="/zebedee/services"
```
The script will output a list of tokens created and the service they belong to
```
1L1YlW7aA2hMMGetIbRv3IE3jIgdqYaXkeF8NTXYyZUh3XyvbHh5tUeYnSSCw0x9    florence
```

8. Once you've finished stop and remove the container
```bash
docker stop <CONTAINER_ID>

docker rm <CONTAINER_ID>
```




