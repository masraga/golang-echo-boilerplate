ALTER TABLE IF EXISTS public.auth_access_token
    ADD COLUMN refresh_token TEXT;

COMMENT ON COLUMN public.auth_access_token.refresh_token IS 'The refresh token to regenerate access token';

