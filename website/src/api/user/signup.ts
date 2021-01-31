import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type SignupUserRequest = {
    emailAddress: string;
    captchaToken: string;
}

export type SignupUserResponse = {
    success: boolean;
    errorMessage: SignupErrorMessage | undefined;
}

export enum SignupErrorMessage {
    InvalidEmailAddress = 'invalid-email',
    IncorrectStatus = 'invalid-account-status',
    RateLimited = 'rate-limited',
}

export function signupUser(
    req: SignupUserRequest,
    onSuccess: (resp: SignupUserResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<SignupUserRequest, SignupUserResponse>(
        '/api/user/signup_user_1',
        req,
        onSuccess,
        onError,
    );
}
