import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type LoginUserRequest = {
    emailAddress: string;
    password: string;
}

export type LoginUserResponse = {
    managementToken: string | null;
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
    PasswordRequirements = 'password-requirements',
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
