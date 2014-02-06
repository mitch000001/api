-- Deploy 02-create-positions-table
-- requires: 01-create-fiscalPeriods-table

BEGIN;

CREATE TYPE position_type AS ENUM ('income', 'expense');
CREATE TYPE position_currency AS ENUM ('EUR', 'USD', 'GBP');

CREATE TABLE public.positions (
  id                SERIAL               PRIMARY KEY,
  category           character varying(16)          NOT NULL DEFAULT '',
  account            character varying(5)           NOT NULL,
  type              position_type        NOT NULL,
  invoice_date      TIMESTAMPTZ          NOT NULL,
  invoice_number     character varying(32)          NOT NULL,
  total_amount      int                  NOT NULL DEFAULT 0,
  currency          position_currency    NOT NULL,
  tax               int                  NOT NULL,
  fiscal_period_id  int                  NOT NULL,
  description       text,
  created_at        TIMESTAMPTZ          NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMPTZ          NOT NULL DEFAULT NOW(),

  CONSTRAINT fiscalPeriodfk FOREIGN KEY (fiscal_period_id) REFERENCES fiscal_periods (id) MATCH FULL
);


COMMIT;
