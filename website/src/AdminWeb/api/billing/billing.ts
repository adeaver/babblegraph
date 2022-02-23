import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

import { PremiumNewsletterSubscription } from 'common/api/billing/billing';
import { SubscriptionLevel } from 'common/api/useraccounts/useraccounts';

export type UserBillingInformation = {
    userId: string;
    userAccountStatus: SubscriptionLevel | undefined;
    externalIdType: string;
    subscriptions: Array<PremiumNewsletterSubscription>;
}

export type GetBillingInformationForEmailAddressRequest = {
    emailAddress: string;
}

export type GetBillingInformationForEmailAddressResponse = {
    billingInformation: UserBillingInformation | undefined;
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
