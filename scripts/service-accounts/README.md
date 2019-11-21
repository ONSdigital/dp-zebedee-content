# Service Account Generation

Set up script for generating Zebedee service accounts.

### Service Accounts.
Service accounts are how our API services authenticate inbound and "sign" outbound requests. Service account files are 
simple `.json` files containing an ID field - the service name (confusingly). The file name is
 a UUID which is the **service token value** used to sign a request/identify a caller.

#### dp-frontend-dataset-controller Example

`1L1YlW7aA2hMMGetIbRv3IE3jIgdqYaXkeF8NTXYyZUh3XyvbHh5tUeYnSSCw0x9.json`
```json
{
  "id": "dp-frontend-dataset-controller"
}
```
If dp-frontend-dataset-controller makes a request to an API requiring service authentication it will set the Auth header with the above token to identify itself:
```
Authorization: Bearer 1L1YlW7aA2hMMGetIbRv3IE3jIgdqYaXkeF8NTXYyZUh3XyvbHh5tUeYnSSCw0x9
```
The API receiving the request will check the request contains an authentication header and make a request to the 
__Identity__ API to verify the caller is known. If the identity check request is successful the request is allowed to 
continue to its destination otherwise it is rejected immediately with the appropriate HTTP status code.

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

4. Clone the script - the script uses Go Modules so should be placed outside of the **$GOPATH**. 

    **Note:** When cloning the repo make sure you use `https` instead of `ssh`.
```
git clone https://github.com/ONSdigital/dp-zebedee-content.git
```

5. Move into the scripts dir and build the generator binary
```
cd dp-zebedee-content/scripts/service-accounts
go build -o generator
```

6. Run the script: Assuming everything is good the script will generate a service account for each service listed in 

        dp-zebedee-content/scripts/service-accounts/generateServiceAccounts.go
 
    `-dir` is the path of the service account directory where the generated service account files will be written. If 
    you have created a volume mapping as defined in step 2 this will be `-dir="/zebedee/services"`.
    
     **The script is non destructive**: If a service account already exists it will not be overwritten. 
```
export HUMAN_LOG=true

./generator -dir="/zebedee/services"
```
The script will output a list of tokens created and the service they belong to
```
1L1YlW7aA2hMMGetIbRv3IE3jIgdqYaXkeF8NTXYyZUh3XyvbHh5tUeYnSSCw0x9    florence
```

7. Once you've generated the required service accounts please stop and remove the container.
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

1. Move into the scripts dir
```
cd dp-zebedee-content/scripts/service-accounts
```

2. Build the generator binary
```
go build -o generator
```

3. Set human readable logging:

`export HUMAN_LOG=true`

4. The script has several flags some of which are optional, see table below:

| flag      | mandatory | example                         | description                                                             |
|-----------|-----------|---------------------------------|-------------------------------------------------------------------------|
| dir       | yes       | <zebedee_root>/zebedee/services | The path to the location of your stored content for zebedee             |
| replace   | no        | true                            | A boolean flag to clear out service auth tokens before regenerating     |
| set-path  | no        | go/src/github.com/ONSdigital    | The path to where your digital publishing services exist                |
| update-mk | no        | true                            | A boolean flag to update makefiles to contain unique SERVICE_AUTH_TOKEN |

a) If you are planning to generate unique `SERVICE_AUTH_TOKEN`S for each service but do not want
to update you Makefile for each service, run the following command:

`./generator -dir=$zebedee_root>/zebedee/services`

b) If you want to add any new `SERVICE_AUTH_TOKEN`'s to service Makefiles, run the following command:

`./generator -dir=$zebedee_root>/zebedee/services -set-path="~/go/src/github.com/ONSdigital" -update-mk=true`

For each new service auth token generated it will add the following line to that services Makefile:
`export SERVICE_AUTH_TOKEN=<generated service auth token>` 

c) If you want to replace all service auth tokens add the replace flag to either a or b options, like so:
```
a) ./generator -dir=$zebedee_root>/zebedee/services -replace=true

b) ./generator -dir=$zebedee_root>/zebedee/services -set-path="~/go/src/github.com/ONSdigital" -update-mk=true -replace=true
```




