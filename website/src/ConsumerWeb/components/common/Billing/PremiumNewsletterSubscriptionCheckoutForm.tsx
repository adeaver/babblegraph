import React from 'react';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

import {
    PaymentState,
    PremiumNewsletterSubscription,
} from 'ConsumerWeb/api/billing/billing';

import ResolvePaymentIntentForm from './stripe/ResolvePaymentIntentForm';
import ResolveSetupIntentForm from './stripe/ResolveSetupIntentForm';

type PremiumNewsletterSubscriptionCardFormProps = {
    premiumNewsletterSusbcription: PremiumNewsletterSubscription;
}

const PremiumNewsletterSubscriptionCardForm = (props: PremiumNewsletterSubscriptionCardFormProps) => {
    const { paymentState } = props.premiumNewsletterSusbcription;
    switch (paymentState) {
        case PaymentState.CreatedUnpaid:
            if (!!props.premiumNewsletterSusbcription.stripePaymentIntentId) {
                return (
                    <ResolvePaymentIntentForm
                            stripePaymentIntentClientSecret={props.premiumNewsletterSusbcription.stripePaymentIntentId} />
                );
            }
            throw new Error("Payment intent ID is not set")
        case PaymentState.TrialNoPaymentMethod:
        case PaymentState.TrialPaymentMethodAdded:
        case PaymentState.Active:
        case PaymentState.Errored:
            return <ResolveSetupIntentForm />;
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
