## Database

CockroachDB is used: http://cockroachlabs.com/

CockroachDB is chosen as it's very scallable and very easy to deploy 
(it's written in Go)

It's similar to SQL databases (even pg driver is fully compatibile) hence it's 
sufficient for this purpose. Cassandra was considered as well but it's much 
more demanding in deployment process (Java)

## Local deploy

Just clone repository and write `make db` to simply quick install and bootstrap
this database in single node in *INSECURE* mode with test data.

## Cloud deploy

Deployment was tested and performed on Digital Ocean Cloud.

First step is to add (no matter what) nodes for cockroachdb, they should be 
tagged. So inn my case I have 3 nodes and all of them are tagged "database". 
The deployment was tested on Ubuntu Server, however it should not be problem to 
deploy it on Centos/RH (unless SELinux blocks something)

Moreover, nodes should resident in same datacenter because of their INSECURE 
configuration I used private ip network which is offered by Digital Ocean.

Hence, each my node has two IP addresses: private and public - private 
network is used to establish INSECURE cluster configuration, public for their 
management and UI.

In ansible/ directory there is `create_ansible_inventory.py` script which creates 
static ansible inventory based on "doctl" tool from DigitalOcean and tags.
This script creates groups and evens divides IP addresses for private and public. 
Because of that we are able to perform configuration CockroachDB cluster as 
ansible playbook knows when to use private IP address.
Such file is created for current configuration and placed in the repo.

In such manner script creates file with tagges hosts which later can be used by 
ansible with `-i file` option.