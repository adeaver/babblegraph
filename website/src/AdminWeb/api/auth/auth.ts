import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type ValidateLoginCredentialsRequest = {
	emailAddress: string;
	password: string;
}

export type ValidateLoginCredentialsResponse = {
    success: boolean;
}

export function validateLoginCredentials(
    req: ValidateLoginCredentialsRequest,
    onSuccess: (resp: ValidateLoginCredentialsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<ValidateLoginCredentialsRequest, ValidateLoginCredentialsResponse>(
        '/ops/api/auth/validate_login_credentials_1',
        req,
        onSuccess,
        onError,
    );
}

export type ValidateTwoFactorAuthenticationCodeRequest = {
    emailAddress: string;
    twoFactorAuthenticationCode: string;
}

export type ValidateTwoFactorAuthenticationCodeResponse = {
	success: boolean;
}

export function validateTwoFactorAuthenticationCode(
    req: ValidateTwoFactorAuthenticationCodeRequest,
    onSuccess: (resp: ValidateTwoFactorAuthenticationCodeResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<ValidateTwoFactorAuthenticationCodeRequest, ValidateTwoFactorAuthenticationCodeResponse>(
        '/ops/api/auth/validate_two_factor_code_1',
        req,
        onSuccess,
        onError,
    );
}

export type InvalidateCredentialsRequest = {}

export type InvalidateCredentialsResponse = {
	success: boolean;
}

export function invalidateCredentials(
    req: InvalidateCredentialsRequest,
    onSuccess: (resp: InvalidateCredentialsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<InvalidateCredentialsRequest, InvalidateCredentialsResponse>(
        '/ops/api/auth/invalidate_login_credentials_1',
        req,
        onSuccess,
        onError,
    );
}

export type CreateAdminUserPasswordRequest = {
	emailAddress: string;
	password: string;
	token: string;
}

export type CreateAdminUserPasswordResponse = {
	success: boolean;
}

export function createAdminUserPassword(
    req: CreateAdminUserPasswordRequest,
    onSuccess: (resp: CreateAdminUserPasswordResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<CreateAdminUserPasswordRequest, CreateAdminUserPasswordResponse>(
        '/ops/api/auth/create_admin_user_password_1',
        req,
        onSuccess,
        onError,
    );
}

export type ValidateTwoFactorAuthenticationCodeForCreateRequest = {
	token: string;
	twoFactorAuthenticationCode: string;
}

export type ValidateTwoFactorAuthenticationCodeForCreateResponse = {
	success: boolean;
}

export function validateTwoFactorAuthenticationCodeForCreate(
    req: ValidateTwoFactorAuthenticationCodeForCreateRequest,
    onSuccess: (resp: ValidateTwoFactorAuthenticationCodeForCreateResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<ValidateTwoFactorAuthenticationCodeForCreateRequest, ValidateTwoFactorAuthenticationCodeForCreateResponse>(
        '/ops/api/auth/validate_two_factor_code_for_create_1',
        req,
        onSuccess,
        onError,
    );
}

