CREATE TABLE IF NOT EXISTS public.oauth (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    username varchar(255),
    email varchar(255),
    refresh_token text NOT NULL,
    access_token text NOT NULL,
    expiry_at_utc0 int8,
    created_at_utc0 int8 NOT NULL DEFAULT (floor(extract(epoch FROM now()) * 1000)::bigint),
    updated_at_utc0 int8,
    is_active boolean DEFAULT TRUE
);

COMMENT ON COLUMN public.oauth.username IS 'Username of authenticate user';

COMMENT ON COLUMN public.oauth.email IS 'Email of authenticate user';

COMMENT ON COLUMN public.oauth.refresh_token IS 'Refresh token of authenticate user';

COMMENT ON COLUMN public.oauth.access_token IS 'Access token of authenticate user';

COMMENT ON COLUMN public.oauth.expiry_at_utc0 IS 'Expiry at unix time of authenticate user';

COMMENT ON COLUMN public.oauth.created_at_utc0 IS 'Created at in UTC+0 epoch time';

