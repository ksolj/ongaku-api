CREATE INDEX IF NOT EXISTS tracks_name_idx ON tracks USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS tracks_artists_idx ON tracks USING GIN (artists);