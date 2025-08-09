CREATE TABLE IF NOT EXISTS rtoken
(
    "UUID" bytea NOT NULL,
    "userUUID" bytea NOT NULL,
    nickname character varying(100) COLLATE pg_catalog."default" NOT NULL,
    "createdAt" integer NOT NULL,
    "expiresAt" integer NOT NULL,
    CONSTRAINT rtoken_pkey PRIMARY KEY ("UUID")
);
