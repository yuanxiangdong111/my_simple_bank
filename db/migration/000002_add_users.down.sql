ALTER TABLE IF Exists "accounts" DROP CONSTRAINT "owner_currency_key";

ALTER TABLE IF Exists "accounts" DROP CONSTRAINT "accounts_owner_fkey";

DROP TABLE IF EXISTS "users";