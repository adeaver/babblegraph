import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type GetSetupIntentForUserRequest = {}

export type GetSetupIntentForUserResponse = {
    setupIntentId: string;
    clientSecret: string;
}

export function getSetupIntentForUser(
    req: GetSetupIntentForUserRequest,
    onSuccess: (resp: GetSetupIntentForUserResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetSetupIntentForUserRequest, GetSetupIntentForUserResponse>(
        '/api/stripe/get_setup_intent_for_user_1',
        req,
        onSuccess,
        onError,
    );
}

export type PaymentMethod = {
    stripePaymentMethodId: string;
    cardType: string;
    lastFourDigits: string;
    expirationMonth: string;
    expirationYear: string;
    isDefault: boolean;
}

export type InsertNewPaymentMethodForUserRequest = {
    stripePaymentMethodId: string;
}

export type InsertNewPaymentMethodForUserResponse = {
    paymentMethod: PaymentMethod;
}

export function insertNewPaymentMethodForUser(
    req: InsertNewPaymentMethodForUserRequest,
    onSuccess: (resp: InsertNewPaymentMethodForUserResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<InsertNewPaymentMethodForUserRequest, InsertNewPaymentMethodForUserResponse>(
        '/api/stripe/insert_payment_method_for_user_1',
        req,
        onSuccess,
        onError,
    );
}
