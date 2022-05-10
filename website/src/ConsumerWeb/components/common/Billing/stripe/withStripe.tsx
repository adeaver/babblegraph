import React from 'react';

import {
    Elements,
    ElementsConsumer,
} from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";

declare const window: any;

export type WithStripeProps = {
    stripe: any;
    elements: any;
}

export type StripeComponentProps = {
    clientSecret: string;
}

export function withStripe<P extends StripeComponentProps>(WrappedComponent: React.ComponentType<P & WithStripeProps>) {
    return (props: P) => {
        const stripePromise = loadStripe(window.initialData["stripe_public_key"]);
        const options = {
            clientSecret: props.clientSecret,
        };
        return (
            // Unclear why, but Elements doesn't think it has children
            // @ts-ignore
            <Elements stripe={stripePromise} options={options}>
                <ElementsConsumer>
                    {({stripe, elements}) => (
                        <WrappedComponent stripe={stripe} elements={elements} {...props} />
                    )}
                </ElementsConsumer>
            </Elements>
        )
    }
}
