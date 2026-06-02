CREATE TABLE IF NOT EXISTS public.auth_api_contract (
    id varchar(255) PRIMARY KEY,
    endpoint_path varchar(255) NOT NULL,
    endpoint_method varchar(16) NOT NULL,
    description text NOT NULL,
    created_at_utc0 int8 NOT NULL DEFAULT (floor(extract(epoch FROM now()) * 1000)::bigint),
    updated_at_utc0 int8 NOT NULL DEFAULT (floor(extract(epoch FROM now()) * 1000)::bigint),
    is_active boolean NOT NULL DEFAULT TRUE
);

COMMENT ON COLUMN public.auth_api_contract.id IS 'OpenAPI operationId that uniquely identifies an API contract';
COMMENT ON COLUMN public.auth_api_contract.endpoint_path IS 'Echo route path protected by this API contract';
COMMENT ON COLUMN public.auth_api_contract.endpoint_method IS 'HTTP method protected by this API contract, stored in lowercase';
COMMENT ON COLUMN public.auth_api_contract.description IS 'Business purpose of the protected API contract';
COMMENT ON COLUMN public.auth_api_contract.created_at_utc0 IS 'Created at in UTC+0 Unix millisecond time';
COMMENT ON COLUMN public.auth_api_contract.updated_at_utc0 IS 'Updated at in UTC+0 Unix millisecond time';
COMMENT ON COLUMN public.auth_api_contract.is_active IS 'Flag indicating the API contract is active and not deleted';

CREATE TABLE IF NOT EXISTS public.auth_user_api_contract (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES public.auth(id),
    api_contract_id varchar(255) NOT NULL REFERENCES public.auth_api_contract(id),
    created_at_utc0 int8 NOT NULL DEFAULT (floor(extract(epoch FROM now()) * 1000)::bigint),
    updated_at_utc0 int8 NOT NULL DEFAULT (floor(extract(epoch FROM now()) * 1000)::bigint),
    is_active boolean NOT NULL DEFAULT TRUE
);

COMMENT ON COLUMN public.auth_user_api_contract.id IS 'Unique identifier for a user API contract grant';
COMMENT ON COLUMN public.auth_user_api_contract.user_id IS 'Auth user id that receives access to an API contract';
COMMENT ON COLUMN public.auth_user_api_contract.api_contract_id IS 'API contract id granted to the auth user';
COMMENT ON COLUMN public.auth_user_api_contract.created_at_utc0 IS 'Created at in UTC+0 Unix millisecond time';
COMMENT ON COLUMN public.auth_user_api_contract.updated_at_utc0 IS 'Updated at in UTC+0 Unix millisecond time';
COMMENT ON COLUMN public.auth_user_api_contract.is_active IS 'Flag indicating the user API contract grant is active and not deleted';

CREATE INDEX IF NOT EXISTS idx_auth_api_contract_endpoint
    ON public.auth_api_contract (endpoint_path, endpoint_method)
    WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_auth_user_api_contract_user
    ON public.auth_user_api_contract (user_id)
    WHERE is_active = TRUE;

CREATE UNIQUE INDEX IF NOT EXISTS uq_auth_user_api_contract_active
    ON public.auth_user_api_contract (user_id, api_contract_id)
    WHERE is_active = TRUE;

INSERT INTO public.auth_api_contract (id, endpoint_path, endpoint_method, description)
VALUES
    ('GetPing', '/api/v1/ping', 'get', 'Check API health'),
    ('RegisterPhoneNumber', '/api/v1/auth/register/phone', 'post', 'Register user phone number'),
    ('VerifyNewAuthUserOTP', '/api/v1/auth/otp/verify', 'post', 'Verify OTP code for a new auth user'),
    ('AuthValidatePin', '/api/v1/auth/validate/pin', 'post', 'Validate user PIN and issue access token'),
    ('CryptoEncryptText', '/api/v1/crypto/encrypt', 'post', 'Encrypt plain text using configured crypto service'),
    ('ListAuthApiContracts', '/api/v1/auth/api-contracts', 'get', 'List active API contracts'),
    ('CreateAuthApiContract', '/api/v1/auth/api-contracts', 'post', 'Create an API contract'),
    ('GetAuthApiContract', '/api/v1/auth/api-contracts/:id', 'get', 'Read an API contract by id'),
    ('UpdateAuthApiContract', '/api/v1/auth/api-contracts/:id', 'put', 'Update an API contract by id'),
    ('DeleteAuthApiContract', '/api/v1/auth/api-contracts/:id', 'delete', 'Soft delete an API contract by id'),
    ('ListAuthUserApiContracts', '/api/v1/auth/user-api-contracts', 'get', 'List active user API contract grants'),
    ('CreateAuthUserApiContract', '/api/v1/auth/user-api-contracts', 'post', 'Create a user API contract grant'),
    ('GetAuthUserApiContract', '/api/v1/auth/user-api-contracts/:id', 'get', 'Read a user API contract grant by id'),
    ('UpdateAuthUserApiContract', '/api/v1/auth/user-api-contracts/:id', 'put', 'Update a user API contract grant by id'),
    ('DeleteAuthUserApiContract', '/api/v1/auth/user-api-contracts/:id', 'delete', 'Soft delete a user API contract grant by id')
ON CONFLICT (id) DO UPDATE SET
    endpoint_path = EXCLUDED.endpoint_path,
    endpoint_method = EXCLUDED.endpoint_method,
    description = EXCLUDED.description,
    updated_at_utc0 = floor(extract(epoch FROM now()) * 1000)::bigint,
    is_active = TRUE;
