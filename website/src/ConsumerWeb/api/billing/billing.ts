import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

import {
    PaymentState,
    PremiumNewsletterSubscription
} from 'common/api/billing/billing';

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

export type LookupActivePremiumNewsletterSubscriptionRequest = {
    subscriptionManagementToken: string;
}

export type LookupActivePremiumNewsletterSubscriptionResponse = {
    premiumNewsletterSubscription: PremiumNewsletterSubscription | undefined;
}

export function lookupActivePremiumNewsletterSubscription(
    req: LookupActivePremiumNewsletterSubscriptionRequest,
    onSuccess: (resp: LookupActivePremiumNewsletterSubscriptionResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<LookupActivePremiumNewsletterSubscriptionRequest, LookupActivePremiumNewsletterSubscriptionResponse>(
        '/api/billing/lookup_active_premium_newsletter_subscription_1',
        req,
        onSuccess,
        onError,
    );
}

export type SetPremiumNewsletterSubscriptionAutoRenewRequest = {
    isAutoRenewEnabled: boolean;
    subscriptionManagementToken: string;
}

export type SetPremiumNewsletterSubscriptionAutoRenewResponse = {
    success: boolean;
}

export function setPremiumNewsletterSubscriptionAutoRenew(
    req: SetPremiumNewsletterSubscriptionAutoRenewRequest,
    onSuccess: (resp: SetPremiumNewsletterSubscriptionAutoRenewResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<SetPremiumNewsletterSubscriptionAutoRenewRequest, SetPremiumNewsletterSubscriptionAutoRenewResponse>(
        '/api/billing/set_premium_newsletter_subscription_auto_renew_1',
        req,
        onSuccess,
        onError,
    );
}

export enum PremiumNewsletterSubscriptionUpdateType {
    TransitionToActive = 'transition-to-active',
    PaymentMethodAdded = 'payment-method-added',
}

export type PreparePremiumNewsletterSubscriptionSyncRequest = {
    id: string;
    updateType: PremiumNewsletterSubscriptionUpdateType;
}

export type PreparePremiumNewsletterSubscriptionSyncResponse = {
    success: boolean;
}

export function preparePremiumNewsletterSubscriptionSync(
    req: PreparePremiumNewsletterSubscriptionSyncRequest,
    onSuccess: (resp: PreparePremiumNewsletterSubscriptionSyncResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<PreparePremiumNewsletterSubscriptionSyncRequest, PreparePremiumNewsletterSubscriptionSyncResponse>(
        '/api/billing/prepare_premium_newsletter_subscription_sync_1',
        req,
        onSuccess,
        onError,
    );
}

export enum CardType {
    Amex = 'amex',
    Visa = 'visa',
    Mastercard = 'mc',
    Discover = 'discover',
    Other = 'other',
}

export type PaymentMethod = {
    externalId: string;
    displayMask: string;
    cardExpiration: string;
    cardType: CardType;
    isDefault: boolean;
}

export type GetPaymentMethodsForUserRequest = {}

export type GetPaymentMethodsForUserResponse = {
    paymentMethods: Array<PaymentMethod>;
}

export function getPaymentMethodsForUser(
    req: GetPaymentMethodsForUserRequest,
    onSuccess: (resp: GetPaymentMethodsForUserResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetPaymentMethodsForUserRequest, GetPaymentMethodsForUserResponse>(
        '/api/billing/get_payment_methods_for_user_1',
        req,
        onSuccess,
        onError,
    );
}

export type MarkPaymentMethodAsDefaultRequest = {
    paymentMethodId: string;
}

export type MarkPaymentMethodAsDefaultResponse = {
    success: boolean;
}

export function markPaymentMethodAsDefault(
    req: MarkPaymentMethodAsDefaultRequest,
    onSuccess: (resp: MarkPaymentMethodAsDefaultResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<MarkPaymentMethodAsDefaultRequest, MarkPaymentMethodAsDefaultResponse>(
        '/api/billing/mark_payment_method_as_default_1',
        req,
        onSuccess,
        onError,
    );
}

export type DeletePaymentMethodForUserRequest = {
    paymentMethodId: string;
}

export type DeletePaymentMethodForUserResponse = {
    success: boolean;
}

export function deletePaymentMethodForUser(
    req: DeletePaymentMethodForUserRequest,
    onSuccess: (resp: DeletePaymentMethodForUserResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<DeletePaymentMethodForUserRequest, DeletePaymentMethodForUserResponse>(
        '/api/billing/delete_payment_method_for_user_1',
        req,
        onSuccess,
        onError,
    );
}
