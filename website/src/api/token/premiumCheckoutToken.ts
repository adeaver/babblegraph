import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

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
