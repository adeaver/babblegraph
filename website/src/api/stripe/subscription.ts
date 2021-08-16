import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type Subscription = {
    stripeSubscriptionId: string;
    paymentState: PaymentState;
    currentPeriodEnd: string;
    cancelAtPeriodEnd: boolean;
    paymentIntentClientSecret: string | undefined;
    subscriptionType: SubscriptionType;
    trialInfo: SubscriptionTrialInfo;
}

export type SubscriptionTrialInfo = {
    isCurrentlyTrialing: boolean;
    trialEligibilityDays: number;
}

export enum SubscriptionType {
    Yearly = 'yearly',
    Monthly = 'monthly',
}

export enum PaymentState {
    CreatedUnpaid = 0,
    TrialNoPaymentMethod = 1,
    TrialPaymentMethodAdded = 2,
    Active = 3,
    Errored = 4,
    Terminated = 5,
}

export type CreateUserSubscriptionRequest = {
    subscriptionType: SubscriptionType;
}

export type CreateUserSubscriptionResponse = {
    subscription: Subscription | undefined;
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

export type GetActiveSubscriptionForUserRequest = {}

export type GetActiveSubscriptionForUserResponse = {
    subscription: Subscription | undefined;
}

export function getActiveSubscriptionForUser(
    req: GetActiveSubscriptionForUserRequest,
    onSuccess: (resp: GetActiveSubscriptionForUserResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetActiveSubscriptionForUserRequest, GetActiveSubscriptionForUserResponse>(
        '/api/stripe/get_active_subscription_for_user_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateSubscriptionOptions = {
    subscriptionType?: SubscriptionType;
    cancelAtPeriodEnd?: boolean;
}

export type UpdateStripeSubscriptionForUserRequest = {
    options: UpdateSubscriptionOptions,
}

export type UpdateStripeSubscriptionForUserResponse = {
    subscription: Subscription | undefined;
}

export function updateStripeSubscriptionForUser(
    req: UpdateStripeSubscriptionForUserRequest,
    onSuccess: (resp: UpdateStripeSubscriptionForUserResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateStripeSubscriptionForUserRequest, UpdateStripeSubscriptionForUserResponse>(
        '/api/stripe/update_stripe_subscription_for_user_1',
        req,
        onSuccess,
        onError,
    );
}
