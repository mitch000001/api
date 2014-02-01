-- Deploy 01-create-fiscalPeriods-table

BEGIN;

SET client_min_messages = 'warning';

CREATE TABLE public.fiscal_periods (
  id          SERIAL         PRIMARY KEY,
  year        INT            NOT NULL,
  created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

COMMIT;
