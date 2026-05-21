CREATE TABLE IF NOT EXISTS public.auth (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_no varchar(15) NOT NULL UNIQUE,
    pin varchar(255),
    is_verified boolean DEFAULT FALSE,
    otp_code varchar(6),
    created_at_utc0 int8 DEFAULT (extract(epoch FROM now())),
    updated_at_utc0 int8,
    is_active boolean DEFAULT true
);

COMMENT ON COLUMN public.auth.created_at_utc0 IS 'Created at in UTC+0 epoch time';

COMMENT ON COLUMN public.auth.updated_at_utc0 IS 'Updated at in UTC+0 epoch time';

COMMENT ON COLUMN public.auth.phone_no IS 'Phone number of the user';

COMMENT ON COLUMN public.auth.pin IS 'PIN of the user';

COMMENT ON COLUMN public.auth.otp_code IS 'OTP code for verification';

COMMENT ON COLUMN public.auth.is_verified IS 'Flag indicating if the user is verified';

COMMENT ON COLUMN public.auth.is_active IS 'Flag indicating the data not deleted';

