CREATE KEYSPACE IF NOT EXISTS "public"
WITH replication = {'class':'SimpleStrategy','replication_factor':'1'};

CREATE TABLE public.message
(
    id           SERIAL PRIMARY KEY,
    email        VARCHAR(50) NOT NULL UNIQUE,
    title        VARCHAR(50) NOT NULL,
    content      TEXT        NOT NULL,
    magic_number INT         NOT NULL
);