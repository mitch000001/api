-- Verify 02-create-positions-table

BEGIN;

SELECT id, total_amount_cents, currency
  FROM public.positions
  WHERE FALSE;

ROLLBACK;
