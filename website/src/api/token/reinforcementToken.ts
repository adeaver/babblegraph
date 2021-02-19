import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type GetReinforcementTokenRequest = {
    token: string;
}

export type GetReinforcementTokenResponse = {
    token: string;
}

export function getReinforcementToken(
    req: GetReinforcementTokenRequest,
    onSuccess: (resp: GetReinforcementTokenResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetReinforcementTokenRequest, GetReinforcementTokenResponse>(
        '/api/token/get_reinforcement_token_1',
        req,
        onSuccess,
        onError,
    );
}
