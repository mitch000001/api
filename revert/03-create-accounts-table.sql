-- Revert 03-create-accounts-table

BEGIN;

DROP TABLE public.accounts;

COMMIT;
