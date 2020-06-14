## General

### This is not security-focused app - do not use on prod!
### Authorization token is very simple
### Authorization method is very simple and based on HTTP

## Database 

### It's starting (single and with as a cluster) in insecure mode, on cloud it's communicating via private network
### Password are stored as sha512 hash but without salt!
### There is no rules against existence more users with same username and password which can lead to unpredictable behaviour
### There is DOS possibility with LOGIN and LOGOUT actions and there is no additional checking if an user is already logged or performed some actions recently. There are no limits. This is very basic
### LOGOUT action only checks existence of token in table users - may be not scallable with larger user database. Also if will be two users with same token then both of them will be reset at the same time.

## HTTP API

### Uses HTTP protocol, not HTTPS - highly insecure for communication



