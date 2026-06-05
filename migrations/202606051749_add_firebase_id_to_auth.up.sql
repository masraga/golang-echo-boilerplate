ALTER TABLE public.auth
    ADD COLUMN firebase_id text;

COMMENT ON COLUMN public.auth.firebase_id IS 'Firebase Cloud Messaging registration token for the user current device';
