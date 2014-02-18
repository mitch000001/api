-- Deploy 03-create-accounts-table

BEGIN;

SET client_min_messages = 'warning';

CREATE TABLE public.accounts (
  code       character varying(5)  PRIMARY KEY,
  label      character varying(16) NOT NULL
);

COMMIT;
