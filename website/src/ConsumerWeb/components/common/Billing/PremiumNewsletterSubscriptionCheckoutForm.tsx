import React from 'react';

import Grid from '@material-ui/core/Grid';

import {
    PaymentState,
    PremiumNewsletterSubscription,
} from 'ConsumerWeb/api/billing/billing';

import ResolvePaymentIntentForm from './stripe/ResolvePaymentIntentForm';

type PremiumNewsletterSubscriptionCardFormProps = {
    premiumNewsletterSusbcription: PremiumNewsletterSubscription;
}

const PremiumNewsletterSubscriptionCardForm = (props: PremiumNewsletterSubscriptionCardFormProps) => {
    const { paymentState } = props.premiumNewsletterSusbcription;
    switch (paymentState) {
        case PaymentState.CreatedUnpaid:
            if (!!props.premiumNewsletterSusbcription.stripePaymentIntentId) {
                return (
                    <Grid container>
                        <Grid item xs={false} md={3}>
                            &nbsp;
                        </Grid>
                        <Grid item xs={12} md={6}>
                            <ResolvePaymentIntentForm
                                stripePaymentIntentClientSecret={props.premiumNewsletterSusbcription.stripePaymentIntentId} />
                        </Grid>
                    </Grid>
                )
            }
            throw new Error("Payment intent ID is not set")
        case PaymentState.TrialNoPaymentMethod:
        case PaymentState.TrialPaymentMethodAdded:
        case PaymentState.Active:
        case PaymentState.Errored:
            return (
                <p>This will show a form to add a payment method</p>
            );
        case PaymentState.Terminated:
            return (
                <p>Your subscription has already ended.</p>
            );
        default:
            throw new Error(`Unrecognized payment state ${paymentState}`)
    }
}

export default PremiumNewsletterSubscriptionCardForm;
