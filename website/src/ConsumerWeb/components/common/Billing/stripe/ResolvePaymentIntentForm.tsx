import React, { useState } from 'react';

import {
    CardNumberElement,
} from "@stripe/react-stripe-js";

import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';

import { withStripe, WithStripeProps } from './withStripe';
import GenericCardForm, { StripeError } from './GenericCardForm';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

import {
    GetPaymentMethodsForUserResponse,
    getPaymentMethodsForUser,

    PremiumNewsletterSubscriptionUpdateType,
    PreparePremiumNewsletterSubscriptionSyncResponse,
    preparePremiumNewsletterSubscriptionSync,
} from 'ConsumerWeb/api/billing/billing';

type StripePaymentIntentResult = {
    error?: StripeError | null;
    paymentIntent?: any;
}

type ResolvePaymentIntentFormOwnProps = {
    premiumNewsletterSubscriptionID: string;
    stripePaymentIntentClientSecret: string;
    toggleSuccessMessage: (shouldShowSuccessMessage) => void;
}

const ResolvePaymentIntentForm = asBaseComponent<GetPaymentMethodsForUserResponse, ResolvePaymentIntentFormOwnProps>(
    withStripe<GetPaymentMethodsForUserResponse & ResolvePaymentIntentFormOwnProps & BaseComponentProps>(
        (props: GetPaymentMethodsForUserResponse & ResolvePaymentIntentFormOwnProps & BaseComponentProps & WithStripeProps) => {
            const [ cardholderName, setCardholderName ] = useState<string>(null);
            const [ postalCode, setPostalCode ] = useState<string>(null);

            const [ isLoading, setIsLoading ] = useState<boolean>(false);

            const [ errorMessage, setErrorMessage ] = useState<string>(null);

            const handleSubmit = () => {
                setIsLoading(true);
                preparePremiumNewsletterSubscriptionSync({
                    id: props.premiumNewsletterSubscriptionID,
                    updateType: PremiumNewsletterSubscriptionUpdateType.TransitionToActive,
                },
                (resp: PreparePremiumNewsletterSubscriptionSyncResponse) => {
                    if (!resp.success) {
                        setIsLoading(false);
                        props.setError(new Error("something went wrong"));
                        return
                    }
                    const cardElement = props.elements.getElement(CardNumberElement);
                    props.stripe.confirmCardPayment(props.stripePaymentIntentClientSecret, {
                        payment_method: {
                            card: cardElement,
                            billing_details: {
                                name: cardholderName,
                                address: {
                                    postal_code: postalCode,
                                },
                            },
                        },
                    }).then((result: StripePaymentIntentResult) => {
                        setIsLoading(false);
                        if (!!result.paymentIntent && result.paymentIntent.status === "succeeded") {
                            props.toggleSuccessMessage(true);
                        } else if (!!result.error) {
                            setErrorMessage(result.error.message);
                        } else {
                            setErrorMessage("There was an error setting up your card");
                        }
                    }).catch((err: Error) => {
                        setErrorMessage("There was an error setting up your card");
                    });
                },
                props.setError);
            }

            return (
                <Form handleSubmit={handleSubmit}>
                    <GenericCardForm
                        cardholderName={cardholderName}
                        postalCode={postalCode}
                        isDisabled={isLoading}
                        setCardholderName={setCardholderName}
                        setPostalCode={setPostalCode} />
                    <PrimaryButton
                        type='submit'
                        disabled={isLoading}>
                        Pay
                    </PrimaryButton>
                </Form>
            );
        }
    ),
    (
        ownProps: ResolvePaymentIntentFormOwnProps,
        onSuccess: (resp: GetPaymentMethodsForUserResponse) => void,
        onError: (err: Error) => void,
    ) => getPaymentMethodsForUser({}, onSuccess, onError),
    false
);

export default ResolvePaymentIntentForm;
