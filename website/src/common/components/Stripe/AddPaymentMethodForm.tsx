import React, { useState, useEffect } from 'react';

import {
    CardElement,
    Elements,
    ElementsConsumer,
} from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";

import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { GenericCardForm, StripeError } from 'common/components/Stripe/common';
import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';

import {
    getSetupIntentForUser,
    GetSetupIntentForUserResponse,

    insertNewPaymentMethodForUser,
    InsertNewPaymentMethodForUserResponse,
} from 'api/stripe/payment_method';

declare const window: any;

type AddPaymentMethodFormProps = {
    handleIsStripeRequestLoading: (isLoading: boolean) => void;
    handleSuccess: (paymentMethodID: string) => void;
    handleFailure: (displayableErrorMessage: string) => void;
    handleError: (err: Error) => void;

    isDefault?: boolean;
}

const AddPaymentMethodForm = (props: AddPaymentMethodFormProps) => {
    const stripePromise = loadStripe(window.initialData["stripe_public_key"]);
    return (
        // Unclear why, but Elements doesn't think it has children
        // @ts-ignore
        <Elements stripe={stripePromise}>
            <InjectedAddPaymentMethodForm {...props} />
        </Elements>
    );
}

const InjectedAddPaymentMethodForm = (props: AddPaymentMethodFormProps) => (
    <ElementsConsumer>
        {({stripe, elements}) => (
            <AddPaymentMethodFormAction stripe={stripe} elements={elements} {...props} />
        )}
    </ElementsConsumer>
);

type StripeSetupIntentResult = {
    error?: StripeError | null;
    setupIntent?: any; // TODO(fill this in?)
}

type AddPaymentMethodFormActionProps = {
    stripe: any;
    elements: any;
} & AddPaymentMethodFormProps;

const AddPaymentMethodFormAction = (props: AddPaymentMethodFormActionProps) => {
    const [ isLoadingStripeRequest, setIsLoadingStripeRequest ] = useState<boolean>(false);
    const [ isLoadingComponent, setIsLoadingComponent ] = useState<boolean>(true);

    const [ clientSecret, setClientSecret ] = useState<string | null>(null);
    const [ setupIntentID, setSetupIntentID ] = useState<string | null>(null);
    const [ error, setError ] = useState<Error>(null);

    const handleIsLoadingStripeRequest = (isLoading: boolean) => {
        setIsLoadingStripeRequest(isLoading);
        props.handleIsStripeRequestLoading(isLoading);
    }

    useEffect(() => {
        getSetupIntentForUser(
            {},
            (resp: GetSetupIntentForUserResponse) => {
                setIsLoadingComponent(false);
                setClientSecret(resp.clientSecret);
                setSetupIntentID(resp.setupIntentId);
            },
            (err: Error) => {
                setIsLoadingComponent(false);
                setError(err);
            }
        );
    }, []);

    const handleSubmit = (cardElement: typeof CardElement, cardholderName: string, postalCode: string) => {
        handleIsLoadingStripeRequest(true);
        props.stripe.confirmCardSetup(clientSecret, {
            payment_method: {
                card: cardElement,
                billing_details: {
                    name: cardholderName,
                    address: {
                        postal_code: postalCode,
                    },
                }
            }
        }).then((result: StripeSetupIntentResult) => {
            if (!!result.setupIntent && result.setupIntent.status === "succeeded") {
                insertNewPaymentMethodForUser({
                    stripePaymentMethodId: result.setupIntent.payment_method,
                },
                (resp: InsertNewPaymentMethodForUserResponse) => {
                    handleIsLoadingStripeRequest(false);
                    props.handleSuccess(resp.paymentMethod.stripePaymentMethodId);
                },
                (err: Error) => {
                    handleIsLoadingStripeRequest(false);
                    props.handleError(err)
                });
            } else if (!!result.error) {
                handleIsLoadingStripeRequest(false);
                props.handleFailure(result.error.message);
            } else {
                handleIsLoadingStripeRequest(false);
                props.handleFailure("There was an error setting up your card");
            }
        }).catch((err: Error) => {
            props.handleError(err)
        });
    }
    if (isLoadingComponent) {
        return <LoadingSpinner />;
    } else if (!!clientSecret && !!setupIntentID) {
        return (
            <div>
                <GenericCardForm
                    actionTitle="Add a payment method"
                    elements={props.elements}
                    isLoading={isLoadingStripeRequest}
                    handleSubmit={handleSubmit} />
            </div>
        );
    } else {
        return (
            <Heading3 color={TypographyColor.Primary}>
                Something went wrong preparing the payment.
            </Heading3>
        );
    }
}

export default AddPaymentMethodForm;
