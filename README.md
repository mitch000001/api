# umsatz-api

# revisioneer

WIP api/ backend for a freelance accounting app

## Tests

To run the testsuite you need to have a PostgreSQL server running & deployed.
Umsatz uses [sqitch][1] for schema management. Thus you need to run

``` bash
createdb umsatz_test
sqitch -d umsatz_test deploy
```

Then you can use `gom` to run the testsuite:

```
REV_DSN="user=$(whoami) dbname=umsatz_test sslmode=disable" gom test
```

## Executing

``` bash
createdb umsatz
sqitch deploy

gom build
REV_DSN="user=$(whoami) dbname=umsatz sslmode=disable" ./umsatz
```

[1]:https://github.com/theory/sqitch