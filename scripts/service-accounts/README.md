# Service Account Generation

Set up script for generating Zebedee service accounts.

### Service Accounts.
Service accounts are how our API services authenticate inbound and "sign" outbound requests. Service account files are 
simple `.json` files containing an ID field - the service name (confusingly). The file name is
 a UUID which is the **service token value** used to sign a request/identify the caller.

#### Example

`1L1YlW7aA2hMMGetIbRv3IE3jIgdqYaXkeF8NTXYyZUh3XyvbHh5tUeYnSSCw0x9.json`
```json
{
  "id": "florence"
}
```
If Florence makes a request to an API requiring service authentication it will set the Auth header with the above token to identify itself:
```
Authorization: Bearer 1L1YlW7aA2hMMGetIbRv3IE3jIgdqYaXkeF8NTXYyZUh3XyvbHh5tUeYnSSCw0x9
```
The API receiving the request will check the request contains an authentication header and make a request to the 
__Identity__ API to verify the caller is known. If the identity check request is successful the request continues to its 
destination otherwise it is rejected immediately with the appropriate HTTP status code.

## Generating Service Accounts

### On an environment 

1. SSH onto the Zebedee publishing-mount.

2. Run a golang docker container with a volume mapping the root zebedee directory on the publishing mount box to `/zebedee`:
```
sudo docker run -i -t --name service-accounts \
   --userns=host \
   -v <ROOT_ZEBEDEE_DIR>:/zebedee:rw \
   golang /bin/bash
```

3. Update the container and install vim (always useful)
```
apt-get update && apt-get install vim
```

4. Clone the script - the script uses Go Modules so should be placed outside of the **$GOPATH**
**Note:** When cloning the repo make sure you use `https` instead of `ssh`.
```
git clone https://github.com/ONSdigital/dp-zebedee-content.git
```

5. Move into the scripts dir and build the generator binary
```
cd dp-zebedee-content/scripts/service-accounts
go build -o generator
```

6. Run the script: Assuming everything is good the script will generate a service account for each service listed under 
`dp-zebedee-content/scripts/service-accounts/generateServiceAccounts.go` var `services`. The script is non destructive 
so if a service account already exists it will not be overwritten. 
 
`-dir` is the path of the service account directory where generated service account files will be written. If you have 
created a volume mapping as defined in step 2 this will be `-dir="/zebedee/services"` 
```
export HUMAN_LOG=true

./generator -dir="/zebedee/services"
```
The script will output a list of tokens created and the service they belong to
```
1L1YlW7aA2hMMGetIbRv3IE3jIgdqYaXkeF8NTXYyZUh3XyvbHh5tUeYnSSCw0x9    florence
```

7. Once you've finished stop and remove the container
```bash
docker stop <CONTAINER_ID>

docker rm <CONTAINER_ID>
```


### Generating service accounts for local set 
To generate service accounts for your local set up:

1 . Clone the script - the script uses Go Modules so should be placed outside of the $GOPATH
```
git clone https://github.com/ONSdigital/dp-zebedee-content.git
```

2. Move into the scripts dir
```
cd dp-zebedee-content/scripts/service-accounts
```

2. Build the generator binary
```
go build -o generator
```

7. Run the script: where `dir` is your `$zebedee_root` + `/zebedee/services` path 
 
```
export HUMAN_LOG=true

./generator -dir="<YOUR_ZEBEDEE_SERVICES_PATH>"
```




