import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export enum PaymentState {
    CreatedUnpaid = 0,
    TrialNoPaymentMethod = 1,
    TrialPaymentMethodAdded = 2,
    Active = 3,
    Errored = 4,
    Terminated = 5,
}

export type CreateUserSubscriptionRequest = {
    subscriptionCreationToken: string;
    isYearlySubscription: boolean;
}

export type CreateUserSubscriptionResponse = {
    stripeSubscriptionId: string;
    stripeClientSecret: string;
    stripePaymentState: PaymentState;
}

export function createUserSubscription(
    req: CreateUserSubscriptionRequest,
    onSuccess: (resp: CreateUserSubscriptionResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<CreateUserSubscriptionRequest, CreateUserSubscriptionResponse>(
        '/api/stripe/create_user_subscription_1',
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
    stripePaymentState: PaymentState | undefined;
    isEligibleForTrial: boolean;
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

export type UpdateStripeSubscriptionFrequencyForUserRequest = {
    stripeSubscriptionId: string;
    isYearlySubscription: boolean;
}

export type UpdateStripeSubscriptionFrequencyForUserResponse = {
    success: boolean;
}

export function updateStripeSubscriptionFrequencyForUser(
    req: UpdateStripeSubscriptionFrequencyForUserRequest,
    onSuccess: (resp: UpdateStripeSubscriptionFrequencyForUserResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateStripeSubscriptionFrequencyForUserRequest, UpdateStripeSubscriptionFrequencyForUserResponse>(
        '/api/stripe/update_stripe_subscription_for_user_1',
        req,
        onSuccess,
        onError,
    );
}


export type DeleteStripeSubscriptionForUserRequest = {
    stripeSubscriptionId: string;
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

