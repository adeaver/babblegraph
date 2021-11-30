import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type RequestPasswordResetLinkRequest = {
    emailAddress: string;
    captchaToken: string;
}

export type RequestPasswordResetLinkResponse = {
    success: boolean;
}

export function requestPasswordResetLink(
    req: RequestPasswordResetLinkRequest,
    onSuccess: (resp: RequestPasswordResetLinkResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<RequestPasswordResetLinkRequest, RequestPasswordResetLinkResponse>(
        '/api/user/handle_request_password_reset_link_1',
        req,
        onSuccess,
        onError,
    );
}
