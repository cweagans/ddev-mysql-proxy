# ddev-mysql-proxy

Goal: one MySQL endpoint that knows about + can proxy to all of your currently running ddev databases.

When you `show databases`, you should see a list of ddev project names. When you `use whateverproject`, further queries should be proxied to the `ddev-whateverproject-db` database server using the `db` database that ddev ships with OOTB.

This project is still very early and mostly doesn't work.

## what works

* compilation
* connecting to the proxy with a mysql cli client
* discovery of all ddev database containers

## what doesn't work

* pretty much everything else

## building

`go mod vendor`
`go build`
