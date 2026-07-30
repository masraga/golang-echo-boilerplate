ALTER TABLE IF EXISTS public.auth
    ADD COLUMN login_mode varchar(255) NOT NULL DEFAULT 'PHONE_NUMBER';

COMMENT ON COLUMN public.auth.login_mode IS 'Default login mode that selected by user when register app';
