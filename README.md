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
make check # install missing deps
DATABASE=umsatz_test go test ./...
```

## Executing

``` bash
createdb umsatz
sqitch deploy

DATABASE=umsatz go run umsatz.go
```

go clean -i -r

## Fake Data

INSERT INTO "fiscal_periods" (year) VALUES (2014);

INSERT INTO "positions" (account_code_from, account_code_to, type, invoice_date, booking_date, invoice_number, total_amount_cents, currency, tax, fiscal_period_id, description, attachment_path) VALUES ('5900', '1100', 'income', '2013-04-04T00:00:00Z', '2013-04-07T00:00:00Z', '20130401', 2001, 'EUR', 1900, 1, '', '');

[1]:https://github.com/theory/sqitch