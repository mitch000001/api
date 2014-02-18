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

Then you can use `gpm` to install dependencies and then run the test suite

```
gpm
DATABASE=umsatz_test go test ./...
```

## Executing

``` bash
createdb umsatz
sqitch deploy

DATABASE=umsatz go run umsatz.go
```

## Fake Data

INSERT INTO "positions" (category, account, type, invoice_date, invoice_number, total_amount, currency, tax, fiscal_period_id, description) VALUES ('foo', '5900', 'income', '2013-04-04T00:00:00Z', '20130401', 2001, 'EUR', 1900, 2, '')

[1]:https://github.com/theory/sqitch