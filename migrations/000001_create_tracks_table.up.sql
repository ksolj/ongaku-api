CREATE TABLE IF NOT EXISTS tracks (
    id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    duration integer NOT NULL,
    artists text[] NOT NULL,
    album text NOT NULL,
    tabs text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);