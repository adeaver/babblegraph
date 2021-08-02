import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type GetOrCreateUserSubscriptionRequest = {
    subscriptionCreationToken: string;
    isYearlySubscription: boolean;
}

export type GetOrCreateUserSubscriptionResponse = {
    stripeSubscriptionId: string;
    stripeClientSecret: string;
    stripePaymentState: number;
}

export function getOrCreateUserSubscription(
    req: GetOrCreateUserSubscriptionRequest,
    onSuccess: (resp: GetOrCreateUserSubscriptionResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetOrCreateUserSubscriptionRequest, GetOrCreateUserSubscriptionResponse>(
        '/api/stripe/get_or_create_user_subscription_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetUserNonTerminatedStripeSubscriptionRequest = {
    subscriptionCreationToken: string;
}

export type GetUserNonTerminatedStripeSubscriptionResponse = {
    isYearlySubscription: boolean | undefined;
    stripeSubscriptionId: string | undefined;
    stripeClientSecret: string | undefined;
    stripePaymentState: number | undefined;
}

export function getUserNonTerminatedStripeSubscription(
    req: GetUserNonTerminatedStripeSubscriptionRequest,
    onSuccess: (resp: GetUserNonTerminatedStripeSubscriptionResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserNonTerminatedStripeSubscriptionRequest, GetUserNonTerminatedStripeSubscriptionResponse>(
        '/api/stripe/get_user_nonterm_stripe_subscription_1',
        req,
        onSuccess,
        onError,
    );
}

export type DeleteStripeSubscriptionForUserRequest = {
    stripeSubscriptionID: string;
}

export type DeleteStripeSubscriptionForUserResponse = {
    success: boolean;
}

export function deleteStripeSubscriptionForUser(
    req: DeleteStripeSubscriptionForUserRequest,
    onSuccess: (resp: DeleteStripeSubscriptionForUserResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<DeleteStripeSubscriptionForUserRequest, DeleteStripeSubscriptionForUserResponse>(
        '/api/stripe/delete_stripe_subscription_for_user_1',
        req,
        onSuccess,
        onError,
    );
}

