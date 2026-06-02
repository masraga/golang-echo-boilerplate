CREATE TABLE IF NOT EXISTS public.auth_roles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name varchar(255) NOT NULL,
    description text NOT NULL,
    owner_id uuid NOT NULL REFERENCES public.auth(id),
    created_at_utc0 int8 NOT NULL DEFAULT (floor(extract(epoch FROM now()) * 1000)::bigint),
    updated_at_utc0 int8 NOT NULL DEFAULT (floor(extract(epoch FROM now()) * 1000)::bigint),
    created_by varchar(255) NOT NULL,
    is_active boolean NOT NULL DEFAULT TRUE
);

COMMENT ON COLUMN public.auth_roles.id IS 'Unique identifier for an auth role';
COMMENT ON COLUMN public.auth_roles.role_name IS 'Human readable role name used to group API contracts';
COMMENT ON COLUMN public.auth_roles.description IS 'Business description of the role purpose';
COMMENT ON COLUMN public.auth_roles.owner_id IS 'Auth user id that owns or created the role';
COMMENT ON COLUMN public.auth_roles.created_at_utc0 IS 'Created at in UTC+0 Unix millisecond time';
COMMENT ON COLUMN public.auth_roles.updated_at_utc0 IS 'Updated at in UTC+0 Unix millisecond time';
COMMENT ON COLUMN public.auth_roles.created_by IS 'Creator identifier for audit purposes';
COMMENT ON COLUMN public.auth_roles.is_active IS 'Flag indicating the role is active and not deleted';

CREATE UNIQUE INDEX IF NOT EXISTS uq_auth_roles_active_role_name
    ON public.auth_roles (role_name)
    WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_auth_roles_owner
    ON public.auth_roles (owner_id)
    WHERE is_active = TRUE;

ALTER TABLE public.auth
    ADD COLUMN IF NOT EXISTS role_name varchar(255),
    ADD COLUMN IF NOT EXISTS role_id uuid REFERENCES public.auth_roles(id),
    ADD COLUMN IF NOT EXISTS created_by varchar(255);

COMMENT ON COLUMN public.auth.role_name IS 'Current assigned auth role name copied from auth_roles.role_name';
COMMENT ON COLUMN public.auth.role_id IS 'Current assigned auth role id from auth_roles';
COMMENT ON COLUMN public.auth.created_by IS 'Creator identifier for audit purposes';

CREATE TABLE IF NOT EXISTS public.auth_roles_contract_api (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id uuid NOT NULL REFERENCES public.auth_roles(id),
    auth_api_contract_id varchar(255) NOT NULL REFERENCES public.auth_api_contract(id),
    created_at_utc0 int8 NOT NULL DEFAULT (floor(extract(epoch FROM now()) * 1000)::bigint),
    updated_at_utc0 int8 NOT NULL DEFAULT (floor(extract(epoch FROM now()) * 1000)::bigint),
    created_by varchar(255) NOT NULL,
    is_active boolean NOT NULL DEFAULT TRUE
);

COMMENT ON COLUMN public.auth_roles_contract_api.id IS 'Unique identifier for a role API contract mapping';
COMMENT ON COLUMN public.auth_roles_contract_api.role_id IS 'Auth role id that receives the API contract';
COMMENT ON COLUMN public.auth_roles_contract_api.auth_api_contract_id IS 'API contract id assigned to the role';
COMMENT ON COLUMN public.auth_roles_contract_api.created_at_utc0 IS 'Created at in UTC+0 Unix millisecond time';
COMMENT ON COLUMN public.auth_roles_contract_api.updated_at_utc0 IS 'Updated at in UTC+0 Unix millisecond time';
COMMENT ON COLUMN public.auth_roles_contract_api.created_by IS 'Creator identifier for audit purposes';
COMMENT ON COLUMN public.auth_roles_contract_api.is_active IS 'Flag indicating the role API contract mapping is active and not deleted';

CREATE UNIQUE INDEX IF NOT EXISTS uq_auth_roles_contract_api_active
    ON public.auth_roles_contract_api (role_id, auth_api_contract_id)
    WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_auth_roles_contract_api_role
    ON public.auth_roles_contract_api (role_id)
    WHERE is_active = TRUE;

INSERT INTO public.auth_api_contract (id, endpoint_path, endpoint_method, description)
VALUES
    ('ListAuthRoles', '/api/v1/auth/roles', 'get', 'List active auth roles'),
    ('CreateAuthRole', '/api/v1/auth/roles', 'post', 'Create an auth role'),
    ('GetAuthRole', '/api/v1/auth/roles/:id', 'get', 'Read an auth role by id'),
    ('UpdateAuthRole', '/api/v1/auth/roles/:id', 'put', 'Update an auth role by id'),
    ('DeleteAuthRole', '/api/v1/auth/roles/:id', 'delete', 'Soft delete an auth role by id'),
    ('ListAuthRoleContractApis', '/api/v1/auth/roles/:roleId/api-contracts', 'get', 'List API contracts assigned to an auth role'),
    ('CreateAuthRoleContractApi', '/api/v1/auth/roles/:roleId/api-contracts', 'post', 'Assign an API contract to an auth role'),
    ('DeleteAuthRoleContractApi', '/api/v1/auth/roles/:roleId/api-contracts/:id', 'delete', 'Remove an API contract from an auth role'),
    ('AssignAuthUserRole', '/api/v1/auth/users/:userId/role', 'put', 'Assign an auth role to an auth user and sync API grants'),
    ('DeleteAuthUserRole', '/api/v1/auth/users/:userId/role', 'delete', 'Remove the assigned role from an auth user and clear API grants')
ON CONFLICT (id) DO UPDATE SET
    endpoint_path = EXCLUDED.endpoint_path,
    endpoint_method = EXCLUDED.endpoint_method,
    description = EXCLUDED.description,
    updated_at_utc0 = floor(extract(epoch FROM now()) * 1000)::bigint,
    is_active = TRUE;
