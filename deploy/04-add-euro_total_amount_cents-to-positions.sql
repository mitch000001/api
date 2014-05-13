-- Deploy 04-add-euro_total_amount_cents-to-positions
-- requires: 02-create-positions-table

BEGIN;

SET client_min_messages = 'warning';

ALTER TABLE public.positions
  ADD COLUMN total_amount_cents_in_eur int NOT NULL DEFAULT 0;

COMMIT;
