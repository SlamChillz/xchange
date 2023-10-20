--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4 (Ubuntu 15.4-2.pgdg20.04+1)
-- Dumped by pg_dump version 15.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', 'public', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;


--
-- Name: customer; Type: TABLE; Schema:  Owner: goswap
--

CREATE TABLE customer (
    id SERIAL PRIMARY KEY,
    last_login TIMESTAMP WITH TIME ZONE,
    photo CHARACTER VARYING(100),
    first_name CHARACTER VARYING(50) NOT NULL,
    last_name CHARACTER VARYING(50) NOT NULL,
    email CHARACTER VARYING(254) UNIQUE NOT NULL,
    password CHARACTER VARYING(300) NOT NULL,
    phone CHARACTER VARYING(30) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE NOT NULL,
    is_staff BOOLEAN DEFAULT FALSE NOT NULL,
    is_supercustomer BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: customerkyc; Type: TABLE; Schema:  Owner: goswap
--

CREATE TABLE customerkyc (
    id SERIAL PRIMARY KEY,
    country_of_residence CHARACTER VARYING(50) NOT NULL,
    first_name CHARACTER VARYING(100) NOT NULL,
    last_name CHARACTER VARYING(100) NOT NULL,
    residential_address CHARACTER VARYING(100) NOT NULL,
    govt_id_photo CHARACTER VARYING(100) NOT NULL,
    govt_id_number CHARACTER VARYING(150) NOT NULL,
    govt_id_type CHARACTER VARYING(150) NOT NULL,
    expiry_date DATE NOT NULL,
    selfie_photo CHARACTER VARYING(100) NOT NULL,
    customer_id SERIAL NOT NULL,
    kyc_verified CHARACTER VARYING(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_customer_di_kyc FOREIGN KEY (customer_id) REFERENCES customer(id) DEFERRABLE INITIALLY DEFERRED
);


--
-- Name: admindata; Type: TABLE; Schema:  Owner: goswap
--

CREATE TABLE admindata (
    id SERIAL PRIMARY KEY,
    bitpowr_account_id CHARACTER VARYING(100) NOT NULL,
    btc_address CHARACTER VARYING(200) DEFAULT NULL,
    usdt_address CHARACTER VARYING(200) DEFAULT NULL,
    usdt_tron_address CHARACTER VARYING(200) DEFAULT NULL,
    admin_email CHARACTER VARYING(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: customerasset; Type: TABLE; Schema:  Owner: goswap
--

CREATE TABLE customerasset (
    id SERIAL PRIMARY KEY,
    btc_address CHARACTER VARYING(200),
    btc_network CHARACTER VARYING(50),
    usdt_tron_address CHARACTER VARYING(200),
    usdc_tron_network CHARACTER VARYING(50),
    usdt_address CHARACTER VARYING(200),
    usdt_network CHARACTER VARYING(50),
    customer_id SERIAL NOT NULL,
    btc_address_uid CHARACTER VARYING(200),
    usdt_tron_address_uid CHARACTER VARYING(200),
    usdt_address_uid CHARACTER VARYING(200),
    usdc_bsc_network CHARACTER VARYING(50),
    usdt_bsc_address CHARACTER VARYING(200),
    usdt_bsc_address_uid CHARACTER VARYING(200),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_customer_id_coin_addresses FOREIGN KEY (customer_id) REFERENCES customer(id)
);


--
-- Name: coinswap; Type: TABLE; Schema:  Owner: goswap
--

CREATE TABLE coinswap (
    id SERIAL PRIMARY KEY,
    coin_name CHARACTER VARYING(100) NOT NULL,
    coin_amount_to_swap numeric(50,8) NOT NULL,
    network CHARACTER VARYING(100) NOT NULL,
    phone_number CHARACTER VARYING(100) NOT NULL,
    coin_address CHARACTER VARYING(200) NOT NULL,
    transaction_ref CHARACTER VARYING(100) UNIQUE NOT NULL,
    transaction_status CHARACTER VARYING(100) NOT NULL,
    current_usdt_ngn_rate CHARACTER VARYING(100) NOT NULL,
    customer_id SERIAL NOT NULL,
    ngn_equivalent numeric(50,8) NOT NULL,
    payout_status CHARACTER VARYING(100) NOT NULL,
    bank_acc_name CHARACTER VARYING(200) NOT NULL,
    bank_acc_number CHARACTER VARYING(50) NOT NULL,
    bitpowr_ref CHARACTER VARYING(200),
    trans_address CHARACTER VARYING(200),
    trans_amount CHARACTER VARYING(200),
    trans_chain CHARACTER VARYING(200),
    trans_hash CHARACTER VARYING(200),
    bank_code CHARACTER VARYING(20) NOT NULL,
    admin_trans_amount CHARACTER VARYING(200),
    admin_trans_fee CHARACTER VARYING(200),
    admin_trans_ref CHARACTER VARYING(200),
    admin_trans_uid CHARACTER VARYING(200),
    trans_amount_ngn CHARACTER VARYING(200),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_customer_id_coinswap FOREIGN KEY (customer_id) REFERENCES customer(id) DEFERRABLE INITIALLY DEFERRED
);


--
-- Name: ratecalculator; Type: TABLE; Schema:  Owner: goswap
--

CREATE TABLE ratecalculator (
    id SERIAL PRIMARY KEY,
    coin_name CHARACTER VARYING(100) NOT NULL,
    coin_amount_to_calc numeric(100,8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: usdtngnrate; Type: TABLE; Schema:  Owner: goswap
--

CREATE TABLE usdtngnrate (
    id SERIAL PRIMARY KEY NOT NULL,
    usdt_ngn_rate numeric(50,8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- PostgreSQL database dump complete
--
