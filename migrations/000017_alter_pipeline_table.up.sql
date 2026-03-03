ALTER TABLE pipelines
    ADD IF NOT EXISTS drm bool NOT NULL default false,
    ADD IF NOT EXISTS key_id text NOT NULL default '';
