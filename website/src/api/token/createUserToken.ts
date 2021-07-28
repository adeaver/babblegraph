import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type GetCreateUserTokenRequest = {
    token: string;
}

export type GetCreateUserTokenResponse = {
    token: string;
}

export function getCreateUserToken(
    req: GetCreateUserTokenRequest,
    onSuccess: (resp: GetCreateUserTokenResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetCreateUserTokenRequest, GetCreateUserTokenResponse>(
        '/api/token/get_create_user_token_1',
        req,
        onSuccess,
        onError,
    );
}
