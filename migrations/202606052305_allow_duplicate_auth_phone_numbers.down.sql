ALTER TABLE public.auth
    ADD CONSTRAINT auth_phone_no_key UNIQUE (phone_no);
