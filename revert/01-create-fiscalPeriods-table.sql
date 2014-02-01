-- Revert 01-create-fiscalPeriods-table

BEGIN;

DROP TABLE public.fiscal_periods;

COMMIT;
