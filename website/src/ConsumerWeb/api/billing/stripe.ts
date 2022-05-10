import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type StripeBeginPaymentMethodSetupRequest = {}

export type StripeBeginPaymentMethodSetupResponse = {
    clientSecret: string;
}

export function stripeBeginPaymentMethodSetup(
    req: StripeBeginPaymentMethodSetupRequest,
    onSuccess: (resp: StripeBeginPaymentMethodSetupResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<StripeBeginPaymentMethodSetupRequest, StripeBeginPaymentMethodSetupResponse>(
        '/api/billing/stripe_begin_payment_method_setup_1',
        req,
        onSuccess,
        onError,
    );
}
