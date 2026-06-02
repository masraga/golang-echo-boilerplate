DELETE FROM public.auth_api_contract
WHERE id IN (
    'ListAuthRoles',
    'CreateAuthRole',
    'GetAuthRole',
    'UpdateAuthRole',
    'DeleteAuthRole',
    'ListAuthRoleContractApis',
    'CreateAuthRoleContractApi',
    'DeleteAuthRoleContractApi',
    'AssignAuthUserRole',
    'DeleteAuthUserRole'
);

DROP INDEX IF EXISTS public.idx_auth_roles_contract_api_role;
DROP INDEX IF EXISTS public.uq_auth_roles_contract_api_active;
DROP TABLE IF EXISTS public.auth_roles_contract_api;

ALTER TABLE public.auth
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS role_id,
    DROP COLUMN IF EXISTS role_name;

DROP INDEX IF EXISTS public.idx_auth_roles_owner;
DROP INDEX IF EXISTS public.uq_auth_roles_active_role_name;
DROP TABLE IF EXISTS public.auth_roles;
