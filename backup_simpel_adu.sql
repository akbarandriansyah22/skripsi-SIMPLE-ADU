--
-- PostgreSQL database dump
--

\restrict aZlm3GO70or3vfeWwz5eJNJIK3Hd9B8RdXuhcv0eTHKTT7XeRBe4uzu8W2715mb

-- Dumped from database version 16.11 (Debian 16.11-1.pgdg13+1)
-- Dumped by pg_dump version 16.11 (Debian 16.11-1.pgdg13+1)

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
-- Name: disposisi; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.disposisi (
    id bigint NOT NULL,
    pengaduan_id bigint NOT NULL,
    pimpinan_id bigint NOT NULL,
    catatan text,
    created_at timestamp with time zone
);


ALTER TABLE public.disposisi OWNER TO postgres;

--
-- Name: disposisi_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.disposisi_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.disposisi_id_seq OWNER TO postgres;

--
-- Name: disposisi_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.disposisi_id_seq OWNED BY public.disposisi.id;


--
-- Name: hasil_ai; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hasil_ai (
    id bigint NOT NULL,
    pengaduan_id bigint NOT NULL,
    sentimen character varying(20) NOT NULL,
    skor_sentimen bigint NOT NULL,
    urgensi character varying(20) NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.hasil_ai OWNER TO postgres;

--
-- Name: hasil_ai_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.hasil_ai_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.hasil_ai_id_seq OWNER TO postgres;

--
-- Name: hasil_ai_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.hasil_ai_id_seq OWNED BY public.hasil_ai.id;


--
-- Name: kategori_pengaduan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.kategori_pengaduan (
    id bigint NOT NULL,
    nama character varying(100) NOT NULL,
    deskripsi text
);


ALTER TABLE public.kategori_pengaduan OWNER TO postgres;

--
-- Name: kategori_pengaduan_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.kategori_pengaduan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.kategori_pengaduan_id_seq OWNER TO postgres;

--
-- Name: kategori_pengaduan_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.kategori_pengaduan_id_seq OWNED BY public.kategori_pengaduan.id;


--
-- Name: mahasiswa; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mahasiswa (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    nim character varying(20) NOT NULL,
    program_studi character varying(100) NOT NULL,
    angkatan bigint NOT NULL,
    no_hp character varying(20)
);


ALTER TABLE public.mahasiswa OWNER TO postgres;

--
-- Name: mahasiswa_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.mahasiswa_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.mahasiswa_id_seq OWNER TO postgres;

--
-- Name: mahasiswa_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.mahasiswa_id_seq OWNED BY public.mahasiswa.id;


--
-- Name: notifikasi; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notifikasi (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    judul character varying(255) NOT NULL,
    isi text NOT NULL,
    is_read boolean DEFAULT false,
    created_at timestamp with time zone
);


ALTER TABLE public.notifikasi OWNER TO postgres;

--
-- Name: notifikasi_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notifikasi_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.notifikasi_id_seq OWNER TO postgres;

--
-- Name: notifikasi_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.notifikasi_id_seq OWNED BY public.notifikasi.id;


--
-- Name: pengaduan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pengaduan (
    id bigint NOT NULL,
    kode_tiket character varying(30) NOT NULL,
    user_id bigint NOT NULL,
    kategori_id bigint NOT NULL,
    unit_id bigint,
    judul character varying(255) NOT NULL,
    deskripsi text NOT NULL,
    lampiran text,
    status character varying(30) DEFAULT 'Menunggu Verifikasi'::character varying,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    tanggal_selesai timestamp with time zone
);


ALTER TABLE public.pengaduan OWNER TO postgres;

--
-- Name: pengaduan_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pengaduan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pengaduan_id_seq OWNER TO postgres;

--
-- Name: pengaduan_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pengaduan_id_seq OWNED BY public.pengaduan.id;


--
-- Name: respon_pengaduan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.respon_pengaduan (
    id bigint NOT NULL,
    pengaduan_id bigint NOT NULL,
    user_id bigint NOT NULL,
    pesan text NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.respon_pengaduan OWNER TO postgres;

--
-- Name: respon_pengaduan_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.respon_pengaduan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.respon_pengaduan_id_seq OWNER TO postgres;

--
-- Name: respon_pengaduan_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.respon_pengaduan_id_seq OWNED BY public.respon_pengaduan.id;


--
-- Name: unit; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.unit (
    id bigint NOT NULL,
    nama_unit character varying(150) NOT NULL,
    email character varying(150)
);


ALTER TABLE public.unit OWNER TO postgres;

--
-- Name: unit_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.unit_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.unit_id_seq OWNER TO postgres;

--
-- Name: unit_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.unit_id_seq OWNED BY public.unit.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    nama_lengkap character varying(150) NOT NULL,
    email character varying(150) NOT NULL,
    password_hash text NOT NULL,
    role character varying(20) NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
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
-- Name: disposisi id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.disposisi ALTER COLUMN id SET DEFAULT nextval('public.disposisi_id_seq'::regclass);


--
-- Name: hasil_ai id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hasil_ai ALTER COLUMN id SET DEFAULT nextval('public.hasil_ai_id_seq'::regclass);


--
-- Name: kategori_pengaduan id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kategori_pengaduan ALTER COLUMN id SET DEFAULT nextval('public.kategori_pengaduan_id_seq'::regclass);


--
-- Name: mahasiswa id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mahasiswa ALTER COLUMN id SET DEFAULT nextval('public.mahasiswa_id_seq'::regclass);


--
-- Name: notifikasi id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifikasi ALTER COLUMN id SET DEFAULT nextval('public.notifikasi_id_seq'::regclass);


--
-- Name: pengaduan id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengaduan ALTER COLUMN id SET DEFAULT nextval('public.pengaduan_id_seq'::regclass);


--
-- Name: respon_pengaduan id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.respon_pengaduan ALTER COLUMN id SET DEFAULT nextval('public.respon_pengaduan_id_seq'::regclass);


--
-- Name: unit id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unit ALTER COLUMN id SET DEFAULT nextval('public.unit_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: disposisi; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.disposisi (id, pengaduan_id, pimpinan_id, catatan, created_at) FROM stdin;
\.


--
-- Data for Name: hasil_ai; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.hasil_ai (id, pengaduan_id, sentimen, skor_sentimen, urgensi, created_at) FROM stdin;
\.


--
-- Data for Name: kategori_pengaduan; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.kategori_pengaduan (id, nama, deskripsi) FROM stdin;
1	Akademik	Pengaduan terkait dosen, jadwal, kurikulum, dan proses perkuliahan
2	Fasilitas	Pengaduan terkait kelas, laboratorium, toilet, dan sarana kampus
3	Kemahasiswaan	Pengaduan terkait organisasi, beasiswa, dan kegiatan mahasiswa
\.


--
-- Data for Name: mahasiswa; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.mahasiswa (id, user_id, nim, program_studi, angkatan, no_hp) FROM stdin;
1	3	2026000001	Teknik Informatika	2026	081234567890
\.


--
-- Data for Name: notifikasi; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.notifikasi (id, user_id, judul, isi, is_read, created_at) FROM stdin;
\.


--
-- Data for Name: pengaduan; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pengaduan (id, kode_tiket, user_id, kategori_id, unit_id, judul, deskripsi, lampiran, status, created_at, updated_at, tanggal_selesai) FROM stdin;
\.


--
-- Data for Name: respon_pengaduan; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.respon_pengaduan (id, pengaduan_id, user_id, pesan, created_at) FROM stdin;
\.


--
-- Data for Name: unit; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.unit (id, nama_unit, email) FROM stdin;
1	Administrasi Fakultas	admin.fakultas@simpel-adu.test
2	Program Studi Teknik Informatika	informatika@simpel-adu.test
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, nama_lengkap, email, password_hash, role, is_active, created_at, updated_at) FROM stdin;
1	Admin Fakultas	admin@simpel-adu.test	$2a$10$YbUnyWuzB4TU1TfxQwPEiOusAS4t2Cc7MXUQ4Y/LXlO7S/nw5QbGa	admin_fakultas	t	2026-07-11 20:28:02.498432+00	2026-07-11 20:28:02.498432+00
2	Pimpinan Fakultas	pimpinan@simpel-adu.test	$2a$10$XAMwjD1xmDfMttRxoGFfueQV9RLX2aM99wipgZIB65y5/oVxa1O9e	pimpinan_fakultas	t	2026-07-11 20:28:02.498432+00	2026-07-11 20:28:02.498432+00
3	Mahasiswa Testing	mahasiswa@simpel-adu.test	$2a$10$igJzQgKAtftSx7uLF1dX3.IdVvdcAN7EIrmet30yoHOAk2aaVnEqa	mahasiswa	t	2026-07-11 20:28:02.498432+00	2026-07-11 20:28:02.498432+00
\.


--
-- Name: disposisi_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.disposisi_id_seq', 1, false);


--
-- Name: hasil_ai_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.hasil_ai_id_seq', 1, false);


--
-- Name: kategori_pengaduan_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.kategori_pengaduan_id_seq', 3, true);


--
-- Name: mahasiswa_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.mahasiswa_id_seq', 1, true);


--
-- Name: notifikasi_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.notifikasi_id_seq', 1, false);


--
-- Name: pengaduan_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.pengaduan_id_seq', 1, false);


--
-- Name: respon_pengaduan_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.respon_pengaduan_id_seq', 1, false);


--
-- Name: unit_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.unit_id_seq', 2, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- Name: disposisi disposisi_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.disposisi
    ADD CONSTRAINT disposisi_pkey PRIMARY KEY (id);


--
-- Name: hasil_ai hasil_ai_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hasil_ai
    ADD CONSTRAINT hasil_ai_pkey PRIMARY KEY (id);


--
-- Name: kategori_pengaduan kategori_pengaduan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kategori_pengaduan
    ADD CONSTRAINT kategori_pengaduan_pkey PRIMARY KEY (id);


--
-- Name: mahasiswa mahasiswa_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mahasiswa
    ADD CONSTRAINT mahasiswa_pkey PRIMARY KEY (id);


--
-- Name: notifikasi notifikasi_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifikasi
    ADD CONSTRAINT notifikasi_pkey PRIMARY KEY (id);


--
-- Name: pengaduan pengaduan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengaduan
    ADD CONSTRAINT pengaduan_pkey PRIMARY KEY (id);


--
-- Name: respon_pengaduan respon_pengaduan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.respon_pengaduan
    ADD CONSTRAINT respon_pengaduan_pkey PRIMARY KEY (id);


--
-- Name: disposisi uni_disposisi_pengaduan_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.disposisi
    ADD CONSTRAINT uni_disposisi_pengaduan_id UNIQUE (pengaduan_id);


--
-- Name: hasil_ai uni_hasil_ai_pengaduan_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hasil_ai
    ADD CONSTRAINT uni_hasil_ai_pengaduan_id UNIQUE (pengaduan_id);


--
-- Name: kategori_pengaduan uni_kategori_pengaduan_nama; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kategori_pengaduan
    ADD CONSTRAINT uni_kategori_pengaduan_nama UNIQUE (nama);


--
-- Name: mahasiswa uni_mahasiswa_nim; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mahasiswa
    ADD CONSTRAINT uni_mahasiswa_nim UNIQUE (nim);


--
-- Name: mahasiswa uni_mahasiswa_user_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mahasiswa
    ADD CONSTRAINT uni_mahasiswa_user_id UNIQUE (user_id);


--
-- Name: pengaduan uni_pengaduan_kode_tiket; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengaduan
    ADD CONSTRAINT uni_pengaduan_kode_tiket UNIQUE (kode_tiket);


--
-- Name: unit uni_unit_nama_unit; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unit
    ADD CONSTRAINT uni_unit_nama_unit UNIQUE (nama_unit);


--
-- Name: users uni_users_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_email UNIQUE (email);


--
-- Name: unit unit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.unit
    ADD CONSTRAINT unit_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: disposisi fk_disposisi_pimpinan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.disposisi
    ADD CONSTRAINT fk_disposisi_pimpinan FOREIGN KEY (pimpinan_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: pengaduan fk_kategori_pengaduan_pengaduan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengaduan
    ADD CONSTRAINT fk_kategori_pengaduan_pengaduan FOREIGN KEY (kategori_id) REFERENCES public.kategori_pengaduan(id);


--
-- Name: mahasiswa fk_mahasiswa_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mahasiswa
    ADD CONSTRAINT fk_mahasiswa_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: disposisi fk_pengaduan_disposisi; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.disposisi
    ADD CONSTRAINT fk_pengaduan_disposisi FOREIGN KEY (pengaduan_id) REFERENCES public.pengaduan(id);


--
-- Name: hasil_ai fk_pengaduan_hasil_ai; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hasil_ai
    ADD CONSTRAINT fk_pengaduan_hasil_ai FOREIGN KEY (pengaduan_id) REFERENCES public.pengaduan(id);


--
-- Name: respon_pengaduan fk_pengaduan_respon_pengaduan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.respon_pengaduan
    ADD CONSTRAINT fk_pengaduan_respon_pengaduan FOREIGN KEY (pengaduan_id) REFERENCES public.pengaduan(id);


--
-- Name: respon_pengaduan fk_respon_pengaduan_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.respon_pengaduan
    ADD CONSTRAINT fk_respon_pengaduan_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: pengaduan fk_unit_pengaduan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengaduan
    ADD CONSTRAINT fk_unit_pengaduan FOREIGN KEY (unit_id) REFERENCES public.unit(id);


--
-- Name: notifikasi fk_users_notifikasi; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifikasi
    ADD CONSTRAINT fk_users_notifikasi FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: pengaduan fk_users_pengaduan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengaduan
    ADD CONSTRAINT fk_users_pengaduan FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict aZlm3GO70or3vfeWwz5eJNJIK3Hd9B8RdXuhcv0eTHKTT7XeRBe4uzu8W2715mb

