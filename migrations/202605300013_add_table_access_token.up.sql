CREATE TABLE IF NOT EXISTS public.access_token (
    id varchar(1024) PRIMARY KEY,
    user_id uuid NOT NULL,
    expired_at_utc0 int8 NOT NULL,
    is_active boolean DEFAULT TRUE
);

COMMENT ON COLUMN public.access_token.id IS 'JWT access token';

COMMENT ON COLUMN public.access_token.user_id IS 'Auth user id';

COMMENT ON COLUMN public.access_token.expired_at_utc0 IS 'Expired at in UTC+0 epoch time';

COMMENT ON COLUMN public.access_token.is_active IS 'Flag indicating the data not deleted';
