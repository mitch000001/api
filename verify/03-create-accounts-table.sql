-- Verify 03-create-accounts-table

BEGIN;

SELECT code, label
  FROM public.accounts
  WHERE FALSE;

ROLLBACK;
