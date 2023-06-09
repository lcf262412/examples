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
-- Name: teste; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.teste (
    codigo integer NOT NULL,
    id integer NOT NULL
);


ALTER TABLE public.teste OWNER TO postgres;

--
-- Name: teste_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.teste ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.teste_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);