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

export enum PaymentState {
    CreatedUnpaid = 0,
    TrialNoPaymentMethod = 1,
    TrialPaymentMethodAdded = 2,
    Active = 3,
    Errored = 4,
    Terminated = 5,
}

export type PremiumNewsletterSubscription = {
    paymentState: PaymentState;
}

export type GetOrCreatePremiumNewsletterSubscriptionRequest = {
    premiumSubscriptionCheckoutToken: string;
}

export type GetOrCreatePremiumNewsletterSubscriptionResponse = {
    premiumNewsletterSubscription: PremiumNewsletterSubscription;
}

export function getOrCreatePremiumNewsletterSubscription(
    req: GetOrCreatePremiumNewsletterSubscriptionRequest,
    onSuccess: (resp: GetOrCreatePremiumNewsletterSubscriptionResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetOrCreatePremiumNewsletterSubscriptionRequest, GetOrCreatePremiumNewsletterSubscriptionResponse>(
        '/api/billing/get_or_create_premium_newsletter_subscription_1',
        req,
        onSuccess,
        onError,
    );
}
