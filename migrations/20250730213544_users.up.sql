CREATE TABLE IF NOT EXISTS users
(
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    lastname character varying(100) COLLATE pg_catalog."default" NOT NULL,
    nickname character varying(100) COLLATE pg_catalog."default" NOT NULL,
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    "userType" character(10) COLLATE pg_catalog."default" NOT NULL,
    "isActive" boolean NOT NULL,
    phone bigint NOT NULL,
    "UUID" bytea NOT NULL,
    "passHash" bytea,
    "createdAt" timestamp without time zone NOT NULL,
    "updateAt" timestamp without time zone,
    CONSTRAINT users_pkey PRIMARY KEY ("UUID"),
    CONSTRAINT uniq_email UNIQUE (email),
    CONSTRAINT uniq_nickname UNIQUE (nickname),
    CONSTRAINT uniq_phone UNIQUE (phone)
);
