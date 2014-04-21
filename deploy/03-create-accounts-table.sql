-- Deploy 03-create-accounts-table

BEGIN;

SET client_min_messages = 'warning';

CREATE TABLE public.accounts (
  id         SERIAL                PRIMARY KEY,
  code       character varying(5),
  label      character varying(16) NOT NULL
);

CREATE UNIQUE INDEX account_code ON public.accounts (code);

COMMIT;
