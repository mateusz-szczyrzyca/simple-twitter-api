## Requirements

This entire application stack was successfully tested and deployed 
under x86_64 Linux OS. Do not try run this, let say under macOS.

To get this working localy you need:
- working go environment (check via `go version` command)
- git
- GNU make

Default http port for application is `58123`


## Application deployment (local)

This app can be deployed locally very easily if you have git and GNU make

`git clone git@github.com:mateusz-szczyrzyca/simple-twitter-api.git`

The repository contains simple Makefile with the following actions:

`make environment` - downloads Go version and simply prepares Go environment for the 
                     rest actions used in Makefile

`make db`    - download cockroachdb binary, unpack it, starts and creates default 
               schema with some test data.
               This is designed only for local configuration only, do not try this
               on prod or internet-exposed environment as Cockroach has no certs 
               and password in such configuration

`make run`   - start application on default settings dsn (look in Makefile if you 
               want change ports or host). Because it won't go into background, use 
               this command with `&` char or inside `screen/tmux` session (logs are 
               shown on stdout)
            

`make tests` - it executed API tests on default database schema (3 users, 20 messages)
               from api_tests/ directory. Before you run this, make sure your application
               is running (hence database with schema and tests data). 
               API tests are NOT unit tests.
               If these tests work correctly - congratulations, you have successfully 
               bootstrapped application locally

               *WARNING*: Because there is no message deleting option, do not run this
               option more than once as data in database will be changed. To verify tests
               AGAIN you have to clean database first by `make clean`


`make clean` - kill database, app and clean db directory

Please have a look in Makefile to change default settings (for DSN etc)

Endpoint address, database dsn string and dbhost can be replaced by your own options:

`make run dbdsn="postgresql://root@localhost:26257/twitter?sslmode=disable" dbhost="localhost" endpoint="localhost:58123"`


## Ansible

In ansible/ directory there are ansible playbooks which can be used to deploy cockroachdb
and app in the cloud. These playbooks were tested on DO cloud and requires `doctl` DO cloud to 
generate static inventory file, based on tags, via `create_ansible_inventory.py` script.

It's able to execute deploy many cockroachdb/app nodes within SAME cluster but in INSECURE 
mode. However, all nodes (hosts) has to be added to the cloud before this action with 
private network and tags. 

You probably won't be able to test this ansible playbook by yourself, however be aware 
that it was used to deploy cockroachdb and app clusters which is used for presentation 
purposes

## JSON samples

In `api_test/json` there are some json samples that can be used with API clien like httpie 
to perform some requests by hand.

### Examples (using httpie)

`http POST localhost:58123/users/login @good_login.json`

`http POST localhost:58123/users/logout @badlogout.json`

`http GET localhost:58123/messages @get_messages_simplest_.json`

`http POST localhost:58123/messages @add_message.json`
