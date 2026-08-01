ALTER TABLE IF EXISTS public.auth
    ADD COLUMN email varchar(255);

COMMENT ON COLUMN public.auth.email IS 'User email that registered with gmail';

