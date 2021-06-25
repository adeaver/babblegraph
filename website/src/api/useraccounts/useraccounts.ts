import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type LoginUserRequest = {
    emailAddress: string;
    password: string;
    redirectKey: string;
}

export type LoginUserResponse = {
    location: string | undefined;
    loginError: LoginError | null;
}

export enum LoginError {
    InvalidCredentials = 'invalid-creds',
}

export function loginUser(
    req: LoginUserRequest,
    onSuccess: (resp: LoginUserResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<LoginUserRequest, LoginUserResponse>(
        '/api/useraccounts/login_user_1',
        req,
        onSuccess,
        onError,
    );
}

export type CreateUserRequest = {
    createUserToken: string;
    emailAddress: string;
    password: string;
    confirmPassword: string;
}

export type CreateUserResponse = {
    managementToken: string | null;
    createUserError: CreateUserError | null;
}

export enum CreateUserError {
    AlreadyExists = 'already-exists',
    InvalidToken = 'invalid-token',
    PasswordRequirements = 'pass-requirements',
    NoSubscription = 'no-subscription',
    PasswordsNoMatch = 'passwords-no-match',
}

export function createUser(
    req: CreateUserRequest,
    onSuccess: (resp: CreateUserResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<CreateUserRequest, CreateUserResponse>(
        '/api/useraccounts/create_user_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetUserProfileRequest = {
    subscriptionManagementToken: string
}

export type GetUserProfileResponse = {
    emailAddress: string | undefined;
    subscriptionLevel: string | undefined;
}

export function getUserProfile(
    req: GetUserProfileRequest,
    onSuccess: (resp: GetUserProfileResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserProfileRequest, GetUserProfileResponse>(
        '/api/useraccounts/get_user_profile_1',
        req,
        onSuccess,
        onError,
    );
}

export enum ResetPasswordError {
    InvalidToken = 'invalid-token',
    TokenExpired = 'token-expired',
    PasswordRequirements = 'pass-requirements',
    PasswordsNoMatch = 'passwords-no-match',
    NoAccount = 'no-account',
}

export type ResetPasswordRequest = {
    resetPasswordToken: string;
    emailAddress: string;
    password: string;
    confirmPassword: string;
}

export type ResetPasswordResponse = {
    managementToken: string | null;
    resetPasswordError: ResetPasswordError | null;
}

export function resetPassword(
    req: ResetPasswordRequest,
    onSuccess: (resp: ResetPasswordResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<ResetPasswordRequest, ResetPasswordResponse>(
        '/api/useraccounts/reset_password_1',
        req,
        onSuccess,
        onError,
    );
}
