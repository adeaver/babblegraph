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
import { PrimarySwitch } from 'common/components/Switch/Switch';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import PaymentMethodDisplay from 'ConsumerWeb/components/common/Billing/PaymentMethodDisplay';
import Alert from 'common/components/Alert/Alert';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import { withStripe, WithStripeProps } from './withStripe';
import GenericCardForm, { StripeError } from './GenericCardForm';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

import {
    PaymentMethod,

    GetPaymentMethodsForUserResponse,
    getPaymentMethodsForUser,

    PremiumNewsletterSubscriptionUpdateType,
    PreparePremiumNewsletterSubscriptionSyncResponse,
    preparePremiumNewsletterSubscriptionSync,
} from 'ConsumerWeb/api/billing/billing';

const styleClasses = makeStyles({
    paymentMethodDisplayContainer: {
        padding: '5px',
    },
    submitButton: {
        width: '100%',
    },
});

type StripePaymentIntentResult = {
    error?: StripeError | null;
    paymentIntent?: any;
}

type ResolvePaymentIntentFormOwnProps = {
    premiumNewsletterSubscriptionID: string;
    clientSecret: string;
    toggleSuccessMessage: (shouldShowSuccessMessage) => void;
}

const ResolvePaymentIntentForm = asBaseComponent<GetPaymentMethodsForUserResponse, ResolvePaymentIntentFormOwnProps>(
    withStripe<GetPaymentMethodsForUserResponse & ResolvePaymentIntentFormOwnProps & BaseComponentProps>(
        (props: GetPaymentMethodsForUserResponse & ResolvePaymentIntentFormOwnProps & BaseComponentProps & WithStripeProps) => {
            const [ existingPaymentMethodIDToUse, setExistingPaymentMethodIDToUse ] = useState<string>(null);
            const [ showCardForm, setShouldShowCardForm ] = useState<boolean>(true);
            const handleSelectPaymentMethodID = (paymentMethodExternalID: string) => {
                if (!isLoading) {
                    setShouldShowCardForm(false);
                    setExistingPaymentMethodIDToUse(paymentMethodExternalID);
                }
            }
            const handleToggleShowCardForm = () => {
                if (props.paymentMethods.length) {
                    if (showCardForm) {
                        setShouldShowCardForm(false);
                    } else {
                        setShouldShowCardForm(true);
                        setExistingPaymentMethodIDToUse(null);
                    }
                }
            }

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
                    props.stripe.confirmPayment({
                        ...props.elements,
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

            const classes = styleClasses();
            return (
                <Form handleSubmit={handleSubmit}>
                    <Grid container>
                    {
                        (props.paymentMethods || []).map((paymentMethod: PaymentMethod) => (
                            <Grid item xs={12} md={6}
                                className={classes.paymentMethodDisplayContainer}
                                key={paymentMethod.externalId}>
                                <PaymentMethodDisplay onClick={handleSelectPaymentMethodID} paymentMethod={paymentMethod} />
                            </Grid>
                        ))
                    }
                    </Grid>
                    <FormControlLabel
                        control={
                            <PrimarySwitch
                                checked={showCardForm}
                                onChange={handleToggleShowCardForm}
                                disabled={isLoading}
                                name="checkbox-show-card-form" />
                        }
                        label="Use new card" />
                    {
                        showCardForm && (
                            <PaymentElement />
                        )
                    }
                    <CenteredComponent>
                        <PrimaryButton
                            className={classes.submitButton}
                            type='submit'
                            disabled={isLoading}>
                            Pay
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
        ownProps: ResolvePaymentIntentFormOwnProps,
        onSuccess: (resp: GetPaymentMethodsForUserResponse) => void,
        onError: (err: Error) => void,
    ) => getPaymentMethodsForUser({}, onSuccess, onError),
    false
);

export default ResolvePaymentIntentForm;
