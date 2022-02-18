import React, { useState } from 'react';

import {
    CardNumberElement,
} from "@stripe/react-stripe-js";

import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';
import Alert from 'common/components/Alert/Alert';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

import {
    StripeBeginPaymentMethodSetupResponse,
    stripeBeginPaymentMethodSetup,
} from 'ConsumerWeb/api/billing/stripe';

import { withStripe, WithStripeProps } from './withStripe';
import GenericCardForm, { StripeError } from './GenericCardForm';

type StripeSetupIntentResult = {
    error?: StripeError | null;
    setupIntent?: any; // TODO(fill this in?)
}

type ResolveSetupIntentFormOwnProps = {
    toggleSuccessMessage: (shouldShowSuccessMessage) => void;
}

const ResolveSetupIntentForm = asBaseComponent<StripeBeginPaymentMethodSetupResponse, ResolveSetupIntentFormOwnProps>(
    withStripe<ResolveSetupIntentFormOwnProps & StripeBeginPaymentMethodSetupResponse & BaseComponentProps>(
        (props: ResolveSetupIntentFormOwnProps & StripeBeginPaymentMethodSetupResponse & BaseComponentProps & WithStripeProps) => {
            const [ cardholderName, setCardholderName ] = useState<string>(null);
            const [ postalCode, setPostalCode ] = useState<string>(null);

            const [ isLoading, setIsLoading ] = useState<boolean>(false);

            const [ errorMessage, setErrorMessage ] = useState<string>(null);

            const handleSubmit = () => {
                setIsLoading(true);
                const cardElement = props.elements.getElement(CardNumberElement);
                props.stripe.confirmCardSetup(props.setupIntentClientSecret, {
                    payment_method: {
                        card: cardElement,
                        billing_details: {
                            name: cardholderName,
                            address: {
                                postal_code: postalCode,
                            },
                        },
                    }
                }).then((result: StripeSetupIntentResult) => {
                    setIsLoading(false);
                    if (!!result.setupIntent && result.setupIntent.status === "succeeded") {
                        props.toggleSuccessMessage(true);
                    } else if (!!result.error) {
                        setErrorMessage(result.error.message);
                    } else {
                        setErrorMessage("There was an error setting up your card");
                    }
                });
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
                        Add Payment Method
                    </PrimaryButton>
                    {
                        isLoading && <LoadingSpinner />
                    }
                    <Snackbar open={!!errorMessage} autoHideDuration={6000} onClose={() => setErrorMessage(null)}>
                        <Alert severity="error">{errorMessage}</Alert>
                    </Snackbar>
                </Form>
            );
        }
    ),
    (
        ownProps: ResolveSetupIntentFormOwnProps,
        onSuccess: (resp: StripeBeginPaymentMethodSetupResponse) => void,
        onError: (err: Error) => void,
    ) => stripeBeginPaymentMethodSetup({}, onSuccess, onError),
    false,
);

export default ResolveSetupIntentForm;
