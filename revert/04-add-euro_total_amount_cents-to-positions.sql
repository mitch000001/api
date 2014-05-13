-- Revert 04-add-euro_total_amount_cents-to-positions

BEGIN;

ALTER TABLE public.positions
  DROP COLUMN total_amount_cents_in_eur;

COMMIT;
