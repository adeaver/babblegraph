import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

import {
    PremiumNewsletterSubscription,
    Discount,
    PromotionCode,
    PromotionType,
} from 'common/api/billing/billing';
import { SubscriptionLevel } from 'common/api/useraccounts/useraccounts';

export type UserBillingInformation = {
    userId: string;
    externalIdType: string;
    subscriptions: Array<PremiumNewsletterSubscription>;
}

export type GetBillingInformationForEmailAddressRequest = {
    emailAddress: string;
}

export type GetBillingInformationForEmailAddressResponse = {
    billingInformation: UserBillingInformation | undefined;
    userAccountStatus: SubscriptionLevel | undefined;
}

export function getBillingInformationForEmailAddress(
    req: GetBillingInformationForEmailAddressRequest,
    onSuccess: (resp: GetBillingInformationForEmailAddressResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetBillingInformationForEmailAddressRequest, GetBillingInformationForEmailAddressResponse>(
        '/ops/api/billing/get_billing_information_for_email_address_1',
        req,
        onSuccess,
        onError,
    );
}

export type ForceSyncForUserRequest = {
    userId: string;
}

export type ForceSyncForUserResponse = {
    success: boolean;
}

export function forceSyncForUser(
    req: ForceSyncForUserRequest,
    onSuccess: (resp: ForceSyncForUserResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<ForceSyncForUserRequest, ForceSyncForUserResponse>(
        '/ops/api/billing/force_sync_for_user_1',
        req,
        onSuccess,
        onError,
    );
}

export type CreatePromotionCodeRequest = {
   code: string;
   discount: Discount;
   maxRedemptions: number | undefined;
   promotionType: PromotionType;
}

export type CreatePromotionCodeResponse = {
    promotionCode: PromotionCode | undefined;
}

export function createPromotionCode(
    req: CreatePromotionCodeRequest,
    onSuccess: (resp: CreatePromotionCodeResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<CreatePromotionCodeRequest, CreatePromotionCodeResponse>(
        '/ops/api/billing/create_promotion_code_1',
        req,
        onSuccess,
        onError,
    );
}
