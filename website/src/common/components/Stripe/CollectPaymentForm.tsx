import React, { useState, useEffect } from 'react';

import {
    CardElement,
    Elements,
    ElementsConsumer,
} from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";

import { GenericCardForm, StripeError } from 'common/components/Stripe/common';
import { Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import { TypographyColor } from 'common/typography/common';

declare const window: any;

type CollectPaymentFormProps = {
    paymentIntentClientSecret: string;

    handleIsStripeRequestLoading: (isLoading: boolean) => void;
    handleSuccess: () => void;
    handleFailure: (displayableErrorMessage: string) => void;
    handleError: (err: Error) => void;
}

const CollectPaymentForm = (props: CollectPaymentFormProps) => {
    const stripePromise = loadStripe(window.initialData["stripe_public_key"]);
    return (
        // Unclear why, but Elements doesn't think it has children
        // @ts-ignore
        <Elements stripe={stripePromise}>
            <InjectedCollectPaymentForm {...props} />
        </Elements>
    )
}

const InjectedCollectPaymentForm = (props: CollectPaymentFormProps) => (
    <ElementsConsumer>
        {({stripe, elements}) => (
            <CollectPaymentFormAction stripe={stripe} elements={elements} {...props} />
        )}
    </ElementsConsumer>
);

type StripePaymentIntentResult = {
    error?: StripeError | null;
    paymentIntent?: any;
}

type CollectPaymentFormActionProps = CollectPaymentFormProps & {
    stripe: any;
    elements: any;
}

const CollectPaymentFormAction = (props: CollectPaymentFormActionProps) => {
    const [ isLoadingStripeRequest, setIsLoadingStripeRequest ] = useState<boolean>(false);
    const [ isLoadingComponent, setIsLoadingComponent ] = useState<boolean>(true);

    const handleIsLoadingStripeRequest = (isLoading: boolean) => {
        setIsLoadingStripeRequest(isLoading);
        props.handleIsStripeRequestLoading(isLoading);
    }

    const handleSubmit = (cardElement: typeof CardElement, cardholderName: string, postalCode: string) => {
        handleIsLoadingStripeRequest(true);
        props.stripe.confirmCardPayment(props.paymentIntentClientSecret, {
            payment_method: {
                card: cardElement,
                billing_details: {
                    name: cardholderName,
                    address: {
                        postal_code: postalCode,
                    },
                }
            }
        }).then((result: StripePaymentIntentResult) => {
            handleIsLoadingStripeRequest(false);
            if (!!result.error) {
                props.handleFailure(result.error.message);
            } else if (!!result.paymentIntent && result.paymentIntent.status === "succeeded") {
                props.handleSuccess();
            } else {
                 props.handleFailure("There was an error completing your payment");
            }
        }).catch((err: Error) => {
            props.handleError(err)
        });
    }
    return (
        <div>
            <GenericCardForm
                actionTitle="Pay"
                elements={props.elements}
                isLoading={isLoadingStripeRequest}
                handleSubmit={handleSubmit} />
        </div>
    );
}

export default CollectPaymentForm;
