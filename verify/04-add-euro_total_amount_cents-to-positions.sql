-- Verify 04-add-euro_total_amount_cents-to-positions

BEGIN;

SELECT total_amount_cents_in_eur
  FROM public.positions
  WHERE false;

ROLLBACK;
