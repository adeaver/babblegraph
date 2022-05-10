import React, { useState } from 'react';

import { Heading3 } from 'common/typography/Heading';
import Link, { LinkTarget } from 'common/components/Link/Link';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

import {
    PaymentState,
    PremiumNewsletterSubscription
} from 'common/api/billing/billing';

import ResolvePaymentIntentForm from './stripe/ResolvePaymentIntentForm';
import ResolveSetupIntentForm from './stripe/ResolveSetupIntentForm';

type PremiumNewsletterSubscriptionCardFormProps = {
    premiumNewsletterSusbcription: PremiumNewsletterSubscription;

    subscriptionManagementToken?: string;
}

const PremiumNewsletterSubscriptionCardForm = (props: PremiumNewsletterSubscriptionCardFormProps) => {
    const { paymentState } = props.premiumNewsletterSusbcription;

    const [ shouldShowSuccessPage, setShouldShowSuccessPage ] = useState<boolean>(false);

    if (shouldShowSuccessPage) {
        return (
            <div>
                <Heading3>
                    Your payment details have been saved successfully it may take up to 10 minutes for your details to successfully appear on your account.
                </Heading3>
                {
                    !!props.subscriptionManagementToken && (
                        <Link href={`/manage/${props.subscriptionManagementToken}`} target={LinkTarget.Self}>
                            Manage your subscription settings
                        </Link>
                    )
                }
            </div>
        );
    }

    switch (paymentState) {
        case PaymentState.CreatedUnpaid:
            if (!!props.premiumNewsletterSusbcription.stripePaymentIntentId) {
                return (
                    <ResolvePaymentIntentForm
                            premiumNewsletterSubscriptionID={props.premiumNewsletterSusbcription.id}
                            clientSecret={props.premiumNewsletterSusbcription.stripePaymentIntentId}
                            toggleSuccessMessage={setShouldShowSuccessPage}
                            redirectURL={!!props.subscriptionManagementToken ? `https://www.babblegraph.com/manage/${props.subscriptionManagementToken}/payment-settings` : `https://www.babblegraph.com/login`} />
                );
            }
            throw new Error("Payment intent ID is not set")
        case PaymentState.TrialNoPaymentMethod:
        case PaymentState.TrialPaymentMethodAdded:
        case PaymentState.Active:
        case PaymentState.Errored:
            return (
                <ResolveSetupIntentForm
                    premiumNewsletterSubscriptionID={props.premiumNewsletterSusbcription.id}
                    toggleSuccessMessage={setShouldShowSuccessPage}
                    redirectURL={!!props.subscriptionManagementToken ? `https://www.babblegraph.com/manage/${props.subscriptionManagementToken}/payment-settings` : `https://www.babblegraph.com/login`} />
            );
        case PaymentState.Terminated:
            // TODO: add redirect to premium information page
            return (
                <p>Your subscription has already ended.</p>
            );
        default:
            throw new Error(`Unrecognized payment state ${paymentState}`)
    }
}

export default PremiumNewsletterSubscriptionCardForm;
