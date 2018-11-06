--
-- PostgreSQL database dump
--

-- Dumped from database version 11.0 (Debian 11.0-1.pgdg90+2)
-- Dumped by pg_dump version 11.0

-- Started on 2018-11-06 13:08:08 WIB

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE srp_demo;
--
-- TOC entry 2877 (class 1262 OID 16384)
-- Name: srp_demo; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE srp_demo WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';


ALTER DATABASE srp_demo OWNER TO postgres;

\connect srp_demo

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 598 (class 1247 OID 16402)
-- Name: session_state; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.session_state AS ENUM (
    'CHALLENGE_SENT',
    'COMPLETED'
);


ALTER TYPE public.session_state OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 196 (class 1259 OID 16385)
-- Name: identities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.identities (
    username character varying(32) NOT NULL,
    salt bytea NOT NULL,
    verifier bytea NOT NULL
);


ALTER TABLE public.identities OWNER TO postgres;

--
-- TOC entry 197 (class 1259 OID 16393)
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    id uuid NOT NULL,
    master_key bytea,
    v bytea,
    state public.session_state,
    "A" bytea,
    b bytea
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- TOC entry 2748 (class 2606 OID 16392)
-- Name: identities identities_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.identities
    ADD CONSTRAINT identities_pkey PRIMARY KEY (username);


--
-- TOC entry 2750 (class 2606 OID 16400)
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


-- Completed on 2018-11-06 13:08:09 WIB

--
-- PostgreSQL database dump complete
--

