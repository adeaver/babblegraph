import React, { useState } from 'react';

import {
    PaymentElement,
} from "@stripe/react-stripe-js";

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
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
import { StripeError } from './GenericCardForm';

type StripeSetupIntentResult = {
    error?: StripeError | null;
    setupIntent?: any; // TODO(fill this in?)
}

const styleClasses = makeStyles({
    submitButton: {
        width: '100%',
    },
});

type ResolveSetupIntentFormOwnProps = {
    toggleSuccessMessage: (shouldShowSuccessMessage) => void;
    redirectURL: string;
}

const ResolveSetupIntentForm = asBaseComponent<StripeBeginPaymentMethodSetupResponse, ResolveSetupIntentFormOwnProps>(
    withStripe<ResolveSetupIntentFormOwnProps & StripeBeginPaymentMethodSetupResponse & BaseComponentProps>(
        (props: ResolveSetupIntentFormOwnProps & StripeBeginPaymentMethodSetupResponse & BaseComponentProps & WithStripeProps) => {
            const [ isLoading, setIsLoading ] = useState<boolean>(false);
            const [ errorMessage, setErrorMessage ] = useState<string>(null);

            const handleSubmit = () => {
                setIsLoading(true);
                props.stripe.confirmSetup({
                    elements: props.elements,
                    confirmParams: {
                        return_url: props.redirectURL,
                    },
                    redirect: 'if_required',
                }).then((result: StripeSetupIntentResult) => {
                    setIsLoading(false);
                    if (!!result.setupIntent && result.setupIntent.status === "succeeded") {
                        props.toggleSuccessMessage(true);
                    } else if (!!result.error) {
                        setErrorMessage(result.error.message);
                    } else {
                        setErrorMessage("There was an error setting up your card");
                    }
                }).catch((err: Error) => {
                    setIsLoading(false);
                    setErrorMessage("There was an error setting up your card");
                });
            }

            const classes = styleClasses();
            return (
                <Form handleSubmit={handleSubmit}>
                    <PaymentElement />
                    <CenteredComponent>
                        <PrimaryButton
                            className={classes.submitButton}
                            type='submit'
                            disabled={isLoading}>
                            Add Payment Method
                        </PrimaryButton>
                    </CenteredComponent>
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
