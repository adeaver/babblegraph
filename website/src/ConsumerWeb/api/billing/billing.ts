import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type GetOrCreateBillingInformationRequest = {
    premiumSubscriptionCheckoutToken: string;
}

export type GetOrCreateBillingInformationResponse = {
    stripeCustomerId: string;
}

export function getOrCreateBillingInformation(
    req: GetOrCreateBillingInformationRequest,
    onSuccess: (resp: GetOrCreateBillingInformationResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetOrCreateBillingInformationRequest, GetOrCreateBillingInformationResponse>(
        '/api/billing/get_or_create_billing_information_1',
        req,
        onSuccess,
        onError,
    );
}
