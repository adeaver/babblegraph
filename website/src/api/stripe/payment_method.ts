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

export type SetDefaultPaymentMethodForUserRequest = {
    stripePaymentMethodId: string;
}

export type SetDefaultPaymentMethodForUserResponse = {}

export function setDefaultPaymentMethodForUser(
    req: SetDefaultPaymentMethodForUserRequest,
    onSuccess: (resp: SetDefaultPaymentMethodForUserResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<SetDefaultPaymentMethodForUserRequest, SetDefaultPaymentMethodForUserResponse>(
        '/api/stripe/set_default_payment_method_for_user_1',
        req,
        onSuccess,
        onError,
    );
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
        '/api/stripe/get_payment_methods_for_user_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetPaymentMethodByIDRequest = {
    stripePaymentMethodId: string;
}

export type GetPaymentMethodByIDResponse = {
    paymentMethod: PaymentMethod | undefined;
}

export function getPaymentMethodByID(
    req: GetPaymentMethodByIDRequest,
    onSuccess: (resp: GetPaymentMethodByIDResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetPaymentMethodByIDRequest, GetPaymentMethodByIDResponse>(
        '/api/stripe/get_payment_method_by_id_1',
        req,
        onSuccess,
        onError,
    );
}

export enum DeletePaymentMethodError {
    DeleteDefault = 'no-delete-default',
    OnlyCard = 'only-card',
}

export type DeletePaymentMethodForUserRequest = {
    stripePaymentMethodId: string;
}

export type DeletePaymentMethodForUserResponse = {
    error: DeletePaymentMethodError | undefined;
}

export function deletePaymentMethodForUser(
    req: DeletePaymentMethodForUserRequest,
    onSuccess: (resp: DeletePaymentMethodForUserResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<DeletePaymentMethodForUserRequest, DeletePaymentMethodForUserResponse>(
        '/api/stripe/delete_payment_method_for_user_1',
        req,
        onSuccess,
        onError,
    );
}
