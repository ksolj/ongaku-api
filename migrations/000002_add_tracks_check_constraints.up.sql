ALTER TABLE tracks ADD CONSTRAINT tracks_duration_check CHECK (duration > 0);

ALTER TABLE tracks ADD CONSTRAINT artists_length_check CHECK (array_length(artists, 1) >= 1);