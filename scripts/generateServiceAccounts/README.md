
# generateServiceAccounts

Intended to be ran on the box, generates service account .json files.


Usage

- `go run generateServiceAccounts.go` to generate all accounts as listed by default (see script)
- run with the flag `-svc <A SERVICE>` to generate a specific service account
- run with the flags `-svc <A SERVICE> -id <AN ID>` to generate a specific service account with a specific ID
In all scenarios the generated ID(s) and service name(s) will be printed to terminal.

Running on the box:

- SSH on to the `publishing_mount` of your chosen environment
- `cd /var/florence/zebedee/services`
- get this script via `curl -o generateServiceAccounts.go https://raw.githubusercontent.com/ONSdigital/dp-zebedee-content/master/scripts/generateServiceAccounts/generateServiceAccounts.go`
- start the go docker image via `sudo docker run -i -t --name <NAME> --userns=host -v /var/florence/zebedee:/zebe-test:rw golang /bin/bash`
- within the container `/zebe-test/services`
- Use script as per usage section

Note - you may need to change <NAME> to a unique value.