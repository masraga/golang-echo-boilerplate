ALTER TABLE public.auth
    ADD COLUMN is_otp_valid boolean NOT NULL DEFAULT FALSE;

COMMENT ON COLUMN public.auth.is_otp_valid IS 'One-time authentication gate indicating that the user has validated the current OTP before PIN authentication';
