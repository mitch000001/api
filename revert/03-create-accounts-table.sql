-- Revert 03-create-accounts-table

BEGIN;

DROP INDEX account_code;
DROP TABLE public.accounts;

COMMIT;
