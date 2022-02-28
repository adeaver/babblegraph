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

export function withStripe<P>(WrappedComponent: React.ComponentType<P & WithStripeProps>) {
    return (props: P) => {
        const stripePromise = loadStripe(window.initialData["stripe_public_key"]);
        return (
            // Unclear why, but Elements doesn't think it has children
            // @ts-ignore
            <Elements stripe={stripePromise}>
                <ElementsConsumer>
                    {({stripe, elements}) => (
                        <WrappedComponent stripe={stripe} elements={elements} {...props} />
                    )}
                </ElementsConsumer>
            </Elements>
        )
    }
}
