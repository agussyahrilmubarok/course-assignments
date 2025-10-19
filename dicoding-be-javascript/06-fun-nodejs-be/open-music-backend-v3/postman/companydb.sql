--
-- PostgreSQL database dump
--

-- Dumped from database version 15.2 (Debian 15.2-1.pgdg110+1)
-- Dumped by pg_dump version 15.2 (Debian 15.2-1.pgdg110+1)

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
-- Name: albums; Type: TABLE; Schema: public; Owner: developer
--

CREATE TABLE public.albums (
    id character varying(50) NOT NULL,
    name text NOT NULL,
    year integer NOT NULL,
    "coverUrl" text
);


ALTER TABLE public.albums OWNER TO developer;

--
-- Name: authentications; Type: TABLE; Schema: public; Owner: developer
--

CREATE TABLE public.authentications (
    token text NOT NULL
);


ALTER TABLE public.authentications OWNER TO developer;

--
-- Name: collaborations; Type: TABLE; Schema: public; Owner: developer
--

CREATE TABLE public.collaborations (
    id character varying(50) NOT NULL,
    playlist_id character varying(50) NOT NULL,
    user_id character varying(50) NOT NULL
);


ALTER TABLE public.collaborations OWNER TO developer;

--
-- Name: pgmigrations; Type: TABLE; Schema: public; Owner: developer
--

CREATE TABLE public.pgmigrations (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    run_on timestamp without time zone NOT NULL
);


ALTER TABLE public.pgmigrations OWNER TO developer;

--
-- Name: pgmigrations_id_seq; Type: SEQUENCE; Schema: public; Owner: developer
--

CREATE SEQUENCE public.pgmigrations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.pgmigrations_id_seq OWNER TO developer;

--
-- Name: pgmigrations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: developer
--

ALTER SEQUENCE public.pgmigrations_id_seq OWNED BY public.pgmigrations.id;


--
-- Name: playlist_song_activities; Type: TABLE; Schema: public; Owner: developer
--

CREATE TABLE public.playlist_song_activities (
    id character varying(50) NOT NULL,
    playlist_id character varying(50) NOT NULL,
    song_id character varying(50) NOT NULL,
    user_id character varying(50) NOT NULL,
    action character varying(50) NOT NULL,
    "time" timestamp without time zone NOT NULL
);


ALTER TABLE public.playlist_song_activities OWNER TO developer;

--
-- Name: playlist_songs; Type: TABLE; Schema: public; Owner: developer
--

CREATE TABLE public.playlist_songs (
    id character varying(50) NOT NULL,
    playlist_id character varying(50) NOT NULL,
    song_id character varying(50) NOT NULL
);


ALTER TABLE public.playlist_songs OWNER TO developer;

--
-- Name: playlists; Type: TABLE; Schema: public; Owner: developer
--

CREATE TABLE public.playlists (
    id character varying(50) NOT NULL,
    name text NOT NULL,
    owner character varying(50) NOT NULL
);


ALTER TABLE public.playlists OWNER TO developer;

--
-- Name: songs; Type: TABLE; Schema: public; Owner: developer
--

CREATE TABLE public.songs (
    id character varying(50) NOT NULL,
    title text NOT NULL,
    year integer NOT NULL,
    genre text NOT NULL,
    performer text NOT NULL,
    duration integer,
    album_id text
);


ALTER TABLE public.songs OWNER TO developer;

--
-- Name: user_album_likes; Type: TABLE; Schema: public; Owner: developer
--

CREATE TABLE public.user_album_likes (
    id character varying(50) NOT NULL,
    user_id character varying(50) NOT NULL,
    album_id character varying(50) NOT NULL
);


ALTER TABLE public.user_album_likes OWNER TO developer;

--
-- Name: users; Type: TABLE; Schema: public; Owner: developer
--

CREATE TABLE public.users (
    id character varying(50) NOT NULL,
    username character varying(50) NOT NULL,
    password text NOT NULL,
    fullname text NOT NULL
);


ALTER TABLE public.users OWNER TO developer;

--
-- Name: pgmigrations id; Type: DEFAULT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.pgmigrations ALTER COLUMN id SET DEFAULT nextval('public.pgmigrations_id_seq'::regclass);


--
-- Data for Name: albums; Type: TABLE DATA; Schema: public; Owner: developer
--

COPY public.albums (id, name, year, "coverUrl") FROM stdin;
album-f1w8fqbPQYc_QNu3	Viva la vida	2008	\N
album-c4i7wwP9ASiGSqG1	Viva la vida	2008	\N
album-5Jn54Fpw8HqFGqWr	Viva la vida	2008	\N
\.


--
-- Data for Name: authentications; Type: TABLE DATA; Schema: public; Owner: developer
--

COPY public.authentications (token) FROM stdin;
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItV2NZOGJCem1hLVlLWkxNRCIsImlhdCI6MTY4ODEzMTAwM30.6f3bIE-tewzKq0U5MrEfSC9I9gVnpU5VH_DldFx6IYs
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItOHhJRU93ZjlXQmlqV1BxWSIsImlhdCI6MTY4ODEzMTAwNH0.6EK94oiPH9RkKsJ-m4qmftmevOTdXfhyj863aAb2d-U
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItV2NZOGJCem1hLVlLWkxNRCIsImlhdCI6MTY4ODEzMTAwNn0.epSCDVyA6F7QpcWRhUJPl4AT0t3oi0cJgoPbZQKWK08
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItOHhJRU93ZjlXQmlqV1BxWSIsImlhdCI6MTY4ODEzMTAwN30.T8FF5Gnbwg9N4mnbKMM5I5qHshq6uID5e141WEPxhlY
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItV2NZOGJCem1hLVlLWkxNRCIsImlhdCI6MTY4ODEzMTAwOH0.u4frq5v2jELXtilrTkOU_isgw6W5CzJYvtldbay6HcE
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItNGxVWkVILXFMcS00ZllDNCIsImlhdCI6MTY4ODEzMTAwOH0.YWUf713Ks1brMkilAPKxwYcXgE7hrsARWPb9Z_Wi-1k
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItV2NZOGJCem1hLVlLWkxNRCIsImlhdCI6MTY4ODEzMTAxMX0.XzmLTtNaD43iTllFPQQz1nGSg585BWs1lM8kX59vajk
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItOHhJRU93ZjlXQmlqV1BxWSIsImlhdCI6MTY4ODEzMTAxMX0.xeNbUtae_Spb4AFRfSoaYM7BtzT9RYlcTjPknwdn96w
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItV2NZOGJCem1hLVlLWkxNRCIsImlhdCI6MTY4ODEzMTAxOH0.B4LdWUkrCFcTHHgoHf-H-_-3AtHmi_k0vSCE_Aezm08
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItOHhJRU93ZjlXQmlqV1BxWSIsImlhdCI6MTY4ODEzMTAxOH0.VMla35cnIgbk5EF4cASOzJKuEXG6RuT6uoSsM8jDOs8
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItV2NZOGJCem1hLVlLWkxNRCIsImlhdCI6MTY4ODEzMTQxN30.wbfqfac-rRnfobejL5JiF1LBUjh8_TlznWjKIdpNdXs
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItV2NZOGJCem1hLVlLWkxNRCIsImlhdCI6MTY4ODEzMTg1NX0.fP9nqcqy4skz2GxnfPiJAWm1aFaTyVKfWwqAuyiKKUY
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItV2NZOGJCem1hLVlLWkxNRCIsImlhdCI6MTY4ODEzMzc2M30.kfg7RlZqn6IEbzLU0H0I_ndBR-4uhtXUKdC2VdPtB9Y
\.


--
-- Data for Name: collaborations; Type: TABLE DATA; Schema: public; Owner: developer
--

COPY public.collaborations (id, playlist_id, user_id) FROM stdin;
\.


--
-- Data for Name: pgmigrations; Type: TABLE DATA; Schema: public; Owner: developer
--

COPY public.pgmigrations (id, name, run_on) FROM stdin;
1	1686995839669_create-albums-table	2023-06-30 13:15:20.272426
2	1686995986215_create-songs-table	2023-06-30 13:15:20.272426
3	1687264364775_create-users-table	2023-06-30 13:15:20.272426
4	1687264449582_create-playlists-table	2023-06-30 13:15:20.272426
5	1687264512807_create-playlist-songs-table	2023-06-30 13:15:20.272426
6	1687264585041_create-collaborations-table	2023-06-30 13:15:20.272426
7	1687264612971_create-playlist-song-activities-table	2023-06-30 13:15:20.272426
8	1687264627412_create-playlist-authentications-table	2023-06-30 13:15:20.272426
9	1687959813887_add-fk-to-albums-in-songs-table	2023-06-30 13:15:20.272426
10	1687959870309_create-album-likes-table	2023-06-30 13:15:20.272426
\.


--
-- Data for Name: playlist_song_activities; Type: TABLE DATA; Schema: public; Owner: developer
--

COPY public.playlist_song_activities (id, playlist_id, song_id, user_id, action, "time") FROM stdin;
song-activities-BYq5I2SJeUOtiwvU	playlist-K6qkP3vjsSFSIeSt	song-iTY8I7KjUhOfwjC-	user-WcY8bBzma-YKZLMD	add	2023-06-30 13:23:46.764
song-activities-0OZsfIvhDNmZeyIb	playlist-K6qkP3vjsSFSIeSt	song-Tj94oI6xnWCe1vsm	user-WcY8bBzma-YKZLMD	add	2023-06-30 13:23:53.101
\.


--
-- Data for Name: playlist_songs; Type: TABLE DATA; Schema: public; Owner: developer
--

COPY public.playlist_songs (id, playlist_id, song_id) FROM stdin;
playlist-song-MmikSMqoPNi_K6LH	playlist-K6qkP3vjsSFSIeSt	song-iTY8I7KjUhOfwjC-
playlist-song-Re4wXXhE50Hc1S7S	playlist-K6qkP3vjsSFSIeSt	song-Tj94oI6xnWCe1vsm
\.


--
-- Data for Name: playlists; Type: TABLE DATA; Schema: public; Owner: developer
--

COPY public.playlists (id, name, owner) FROM stdin;
playlist-K6qkP3vjsSFSIeSt	Lagu Indie Hits Indonesia	user-WcY8bBzma-YKZLMD
playlist-v8ySoPuTE-A9L9px	Lagu Indie Hits Indonesia	user-WcY8bBzma-YKZLMD
\.


--
-- Data for Name: songs; Type: TABLE DATA; Schema: public; Owner: developer
--

COPY public.songs (id, title, year, genre, performer, duration, album_id) FROM stdin;
song-T7Sm1fSxbIg0JbqN	Fix You	2008	Pop	Coldplay	120	\N
song-1GLX-xIbSk7p3Pf-	Life in Technicolor	2008	Pop	Coldplay	120	album-f1w8fqbPQYc_QNu3
song-5JJQesnJbR_vOsze	Fix you	2008	Pop	Coldplay	120	album-f1w8fqbPQYc_QNu3
song-O9U9BT93cSIgyhzz	Life in Technicolor	2008	Pop	Coldplay	120	\N
song-i0rzcgqE5WcyeNHL	Fix You	2008	Pop	Coldplay	120	\N
song-lVQd_aZ01jY0bEnL	Life in Technicolor	2008	Pop	Coldplay	120	\N
song-QsW5Obr_6k5jeYsT	Fix You	2008	Pop	Coldplay	120	\N
song-yB58A_bujz3bUFl3	Life in Technicolor	2008	Pop	Coldplay	120	\N
song-0gt2w6QbIzOdmnhF	Fix You	2008	Pop	Coldplay	120	\N
song-9j8Co7hPLEwz2HRZ	Life in Technicolor	2008	Pop	Coldplay	120	\N
song-dMrFJRd36tui2DT8	Fix You	2008	Pop	Coldplay	120	\N
song-iTY8I7KjUhOfwjC-	Life in Technicolor	2008	Pop	Coldplay	120	\N
song-Tj94oI6xnWCe1vsm	Fix You	2008	Pop	Coldplay	120	\N
\.


--
-- Data for Name: user_album_likes; Type: TABLE DATA; Schema: public; Owner: developer
--

COPY public.user_album_likes (id, user_id, album_id) FROM stdin;
likes-i074ubEaLYYNmQt2	user-WcY8bBzma-YKZLMD	album-5Jn54Fpw8HqFGqWr
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: developer
--

COPY public.users (id, username, password, fullname) FROM stdin;
user-9MAvsaErPrf2_T1s	dicoding_1688131000	$2b$10$te5zT1uMbLijRdxEtXIGJeFH4XSXjOyt5or6npKqmg6bm8I6SfjOy	Dicoding Indonesia
user-nPiuW1TRmkQHohGJ	dicoding	$2b$10$1X9px.leKYLtiRzAHZ1ndO/.6uPleEB4iwj8fonl/2/r3FuyP74mS	Dicoding Indonesia
user-WcY8bBzma-YKZLMD	john	$2b$10$OQuzTBpMQmbDmQnT64GZaeB0VSQ7BI9Yrbbhr07qDvvds1esYXjwW	John Doe
user-8xIEOwf9WBijWPqY	jane	$2b$10$cAtAmmcpyh5C50QGDak6RuddFU4.8AUH9xIA9tRRlRaukiZxQguWm	John Doe
user-4lUZEH-qLq-4fYC4	tom1688131008283	$2b$10$VEaSG7ibueLMfAvB5CDulO5jjh2GYXLTjgolhFONMJQxDlxVapMwe	Tom Riddle
\.


--
-- Name: pgmigrations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: developer
--

SELECT pg_catalog.setval('public.pgmigrations_id_seq', 10, true);


--
-- Name: albums albums_pkey; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.albums
    ADD CONSTRAINT albums_pkey PRIMARY KEY (id);


--
-- Name: collaborations collaborations_pkey; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.collaborations
    ADD CONSTRAINT collaborations_pkey PRIMARY KEY (id);


--
-- Name: pgmigrations pgmigrations_pkey; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.pgmigrations
    ADD CONSTRAINT pgmigrations_pkey PRIMARY KEY (id);


--
-- Name: playlist_song_activities playlist_song_activities_pkey; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.playlist_song_activities
    ADD CONSTRAINT playlist_song_activities_pkey PRIMARY KEY (id);


--
-- Name: playlist_songs playlist_songs_pkey; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.playlist_songs
    ADD CONSTRAINT playlist_songs_pkey PRIMARY KEY (id);


--
-- Name: playlists playlists_pkey; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.playlists
    ADD CONSTRAINT playlists_pkey PRIMARY KEY (id);


--
-- Name: songs songs_pkey; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.songs
    ADD CONSTRAINT songs_pkey PRIMARY KEY (id);


--
-- Name: playlist_songs unique_playlist_id_and_song_id; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.playlist_songs
    ADD CONSTRAINT unique_playlist_id_and_song_id UNIQUE (playlist_id, song_id);


--
-- Name: collaborations unique_playlist_id_and_user_id; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.collaborations
    ADD CONSTRAINT unique_playlist_id_and_user_id UNIQUE (playlist_id, user_id);


--
-- Name: user_album_likes user_album_likes_pkey; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.user_album_likes
    ADD CONSTRAINT user_album_likes_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: collaborations fk_collaborations.playlist_id_playlists.id; Type: FK CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.collaborations
    ADD CONSTRAINT "fk_collaborations.playlist_id_playlists.id" FOREIGN KEY (playlist_id) REFERENCES public.playlists(id) ON DELETE CASCADE;


--
-- Name: collaborations fk_collaborations.user_id_users.id; Type: FK CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.collaborations
    ADD CONSTRAINT "fk_collaborations.user_id_users.id" FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: playlist_song_activities fk_playlist_song_activities.playlist_id_playlists.id; Type: FK CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.playlist_song_activities
    ADD CONSTRAINT "fk_playlist_song_activities.playlist_id_playlists.id" FOREIGN KEY (playlist_id) REFERENCES public.playlists(id) ON DELETE CASCADE;


--
-- Name: playlist_songs fk_playlist_songs.playlist_id_playlists.id; Type: FK CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.playlist_songs
    ADD CONSTRAINT "fk_playlist_songs.playlist_id_playlists.id" FOREIGN KEY (playlist_id) REFERENCES public.playlists(id) ON DELETE CASCADE;


--
-- Name: playlist_songs fk_playlist_songs.song_id_songs.id; Type: FK CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.playlist_songs
    ADD CONSTRAINT "fk_playlist_songs.song_id_songs.id" FOREIGN KEY (song_id) REFERENCES public.songs(id) ON DELETE CASCADE;


--
-- Name: playlists fk_playlists.owner_users.id; Type: FK CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.playlists
    ADD CONSTRAINT "fk_playlists.owner_users.id" FOREIGN KEY (owner) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: songs fk_songs.album_id_albums.id; Type: FK CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.songs
    ADD CONSTRAINT "fk_songs.album_id_albums.id" FOREIGN KEY (album_id) REFERENCES public.albums(id) ON DELETE CASCADE;


--
-- Name: user_album_likes fk_user_album_likes.album_id_albums.id; Type: FK CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.user_album_likes
    ADD CONSTRAINT "fk_user_album_likes.album_id_albums.id" FOREIGN KEY (album_id) REFERENCES public.albums(id) ON DELETE CASCADE;


--
-- Name: user_album_likes fk_user_album_likes.user_id_users.id; Type: FK CONSTRAINT; Schema: public; Owner: developer
--

ALTER TABLE ONLY public.user_album_likes
    ADD CONSTRAINT "fk_user_album_likes.user_id_users.id" FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

