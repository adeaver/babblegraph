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
