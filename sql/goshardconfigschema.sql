--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: database_mappings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.database_mappings (
    id integer NOT NULL,
    shardid smallint,
    sharduid character varying(50),
    dsn character varying(255),
    user_id bigint
);


ALTER TABLE public.database_mappings OWNER TO postgres;

--
-- Name: database_mappings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.database_mappings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.database_mappings_id_seq OWNER TO postgres;

--
-- Name: database_mappings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.database_mappings_id_seq OWNED BY public.database_mappings.id;


--
-- Name: user_schemas; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_schemas (
    id integer NOT NULL,
    user_id bigint,
    schema text
);


ALTER TABLE public.user_schemas OWNER TO postgres;

--
-- Name: user_schemas_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_schemas_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_schemas_id_seq OWNER TO postgres;

--
-- Name: user_schemas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_schemas_id_seq OWNED BY public.user_schemas.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    token character varying(50) NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: database_mappings id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.database_mappings ALTER COLUMN id SET DEFAULT nextval('public.database_mappings_id_seq'::regclass);


--
-- Name: user_schemas id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_schemas ALTER COLUMN id SET DEFAULT nextval('public.user_schemas_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: database_mappings database_mappings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.database_mappings
    ADD CONSTRAINT database_mappings_pkey PRIMARY KEY (id);


--
-- Name: user_schemas user_schemas_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_schemas
    ADD CONSTRAINT user_schemas_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

