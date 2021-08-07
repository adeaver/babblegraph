import React, { useState, useEffect } from 'react';

import {
    CardElement,
    Elements,
    ElementsConsumer,
} from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";

import { GenericCardForm, StripeError } from 'common/components/Stripe/common';
import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';

declare const window: any;

type AddPaymentMethodFormProps = {
    customerID: string;
    isDefault?: boolean;

    handleIsStripeRequestLoading: (isLoading: boolean) => void;
    handleSuccess: (paymentMethodID: string) => void;
    handleFailure: (displayableErrorMessage: string) => void;
    handleError: (err: Error) => void;
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

    const handleIsLoadingStripeRequest = (isLoading: boolean) => {
        setIsLoadingStripeRequest(isLoading);
        props.handleIsStripeRequestLoading(isLoading);
    }

    useEffect(() => {
        // TODO: Create Setup Intent
    }, []);

    const handleSubmit = (cardElement: typeof CardElement, cardholderName: string, postalCode: string) => {
        handleIsLoadingStripeRequest(true);
        props.stripe.confirmCardSetup("", {
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
                // TODO: Send request to Babblegraph API for payment method
                props.handleSuccess("");
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
    return (
        <div>
            <GenericCardForm
                actionTitle="Add a payment method"
                elements={props.elements}
                isLoading={isLoadingStripeRequest}
                handleSubmit={handleSubmit} />
        </div>
    );
}

export default AddPaymentMethodForm;
