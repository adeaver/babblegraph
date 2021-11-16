import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type GetPremiumCheckoutTokenRequest = {
    token: string;
}

export type GetPremiumCheckoutTokenResponse = {
    token: string;
}

export function getPremiumCheckoutToken(
    req: GetPremiumCheckoutTokenRequest,
    onSuccess: (resp: GetPremiumCheckoutTokenResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetPremiumCheckoutTokenRequest, GetPremiumCheckoutTokenResponse>(
        '/api/token/get_premium_checkout_token_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetManageTokenForPremiumCheckoutTokenRequest = {
    token: string;
}

export type GetManageTokenForPremiumCheckoutTokenResponse = {
    token: string;
}

export function getManageTokenForPremiumCheckoutToken(
    req: GetManageTokenForPremiumCheckoutTokenRequest,
    onSuccess: (resp: GetManageTokenForPremiumCheckoutTokenResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetManageTokenForPremiumCheckoutTokenRequest, GetManageTokenForPremiumCheckoutTokenResponse>(
        '/api/token/get_manage_token_for_premium_checkout_token_1',
        req,
        onSuccess,
        onError,
    );
}
