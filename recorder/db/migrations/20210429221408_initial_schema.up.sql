CREATE TABLE IF NOT EXISTS websites (
       id TEXT PRIMARY KEY,
       url TEXT NOT NULL,
       method TEXT NOT NULL,
       match_regexp TEXT
);

CREATE TABLE IF NOT EXISTS websites_results (
       website_id TEXT REFERENCES websites(id),
       elapsed_time DOUBLE PRECISION,
       status INT,
       matched BOOLEAN,
       unreachable BOOLEAN DEFAULT FALSE,
       at TIMESTAMP WITHOUT TIME ZONE,

       PRIMARY KEY (website_id, at)
);

-- Get latest results.
CREATE INDEX IF NOT EXISTS index_websites_results_on_at_desc ON websites_results(at DESC);

-- Get first results.
CREATE INDEX IF NOT EXISTS index_websites_results_on_at_asc ON websites_results(at ASC);
