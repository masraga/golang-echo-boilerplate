CREATE TABLE IF NOT EXISTS public.auth_otp (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id uuid NOT NULL,
    otp_code varchar(6) NOT NULL,
    note varchar(255),
    expired_at_utc0 int8 NOT NULL,
    created_at_utc0 int8 DEFAULT extract(epoch FROM now()),
    updated_at_utc0 int8 DEFAULT extract(epoch FROM now()),
    is_active boolean DEFAULT TRUE,
    is_verified boolean DEFAULT FALSE
);

COMMENT ON COLUMN public.auth_otp.created_at_utc0 IS 'Created at in UTC+0 epoch time';

COMMENT ON COLUMN public.auth_otp.otp_code IS 'OTP code for verification';

COMMENT ON COLUMN public.auth_otp.expired_at_utc0 IS 'Expired at in UTC+0 epoch time';

COMMENT ON COLUMN public.auth_otp.expired_at_utc0 IS 'Expired at in UTC+0 epoch time';

COMMENT ON COLUMN public.auth_otp.updated_at_utc0 IS 'Updated at in UTC+0 epoch time';

COMMENT ON COLUMN public.auth_otp.is_active IS 'Flag indicating the data not deleted';

COMMENT ON COLUMN public.auth_otp.note IS 'Additional note for the OTP record';

COMMENT ON COLUMN public.auth_otp.is_verified IS 'Flag indicating whether the OTP has been verified';