-- Verify 01-create-fiscalPeriods-table

BEGIN;

SELECT id, year, created_at, updated_at
  FROM public.fiscal_periods
  WHERE FALSE;

ROLLBACK;
