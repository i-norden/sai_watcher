--
-- PostgreSQL database dump
--

-- Dumped from database version 10.1
-- Dumped by pg_dump version 10.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: maker; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA maker;


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = maker, pg_catalog;

--
-- Name: act; Type: TYPE; Schema: maker; Owner: -
--

CREATE TYPE act AS ENUM (
    'open',
    'join',
    'exit',
    'lock',
    'free',
    'draw',
    'wipe',
    'shut',
    'bite'
);


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: cup_action; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE cup_action (
    log_id integer,
    id integer,
    tx character varying(66) NOT NULL,
    act act NOT NULL,
    arg character varying(66),
    lad character varying(66) NOT NULL,
    ink numeric DEFAULT 0 NOT NULL,
    art numeric DEFAULT 0 NOT NULL,
    ire numeric DEFAULT 0 NOT NULL,
    block integer NOT NULL,
    deleted boolean DEFAULT false
);


--
-- Name: peps_everyblock; Type: TABLE; Schema: maker; Owner: -
--

CREATE TABLE peps_everyblock (
    id integer NOT NULL,
    block_number integer NOT NULL,
    pep numeric,
    pip numeric,
    per numeric
);


--
-- Name: peps_everyblock_id_seq; Type: SEQUENCE; Schema: maker; Owner: -
--

CREATE SEQUENCE peps_everyblock_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: peps_everyblock_id_seq; Type: SEQUENCE OWNED BY; Schema: maker; Owner: -
--

ALTER SEQUENCE peps_everyblock_id_seq OWNED BY peps_everyblock.id;


SET search_path = public, pg_catalog;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


SET search_path = maker, pg_catalog;

--
-- Name: peps_everyblock id; Type: DEFAULT; Schema: maker; Owner: -
--

ALTER TABLE ONLY peps_everyblock ALTER COLUMN id SET DEFAULT nextval('peps_everyblock_id_seq'::regclass);


SET search_path = public, pg_catalog;

--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- PostgreSQL database dump complete
--

