-- Revert 02-create-positions-table

BEGIN;

DROP TABLE public.positions;

DROP TYPE position_type;
DROP TYPE position_currency;

COMMIT;
