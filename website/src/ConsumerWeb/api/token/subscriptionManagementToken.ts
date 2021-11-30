import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type GetManageTokenForReinforcementTokenRequest = {
    token: string;
}

export type GetManageTokenForReinforcementTokenResponse = {
    token: string;
}

export function getManageTokenForReinforcementToken(
    req: GetManageTokenForReinforcementTokenRequest,
    onSuccess: (resp: GetManageTokenForReinforcementTokenResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetManageTokenForReinforcementTokenRequest, GetManageTokenForReinforcementTokenResponse>(
        '/api/token/get_manage_token_for_reinforcement_token_1',
        req,
        onSuccess,
        onError,
    );
}
