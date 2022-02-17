import React, { useState, useRef, useImperativeHandle } from 'react';

import {
    CardNumberElement,
    CardExpiryElement,
    CardCvcElement,
    CardElement,
} from "@stripe/react-stripe-js";

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import Form from 'common/components/Form/Form';

const styleClasses = makeStyles({
    checkoutFormObject: {
        width: "100%",
        margin: "10px 0",
        paddingLeft: "5px",
        paddingRight: "5px",
        boxSizing: "border-box",
    },
    stripeBadge: {
        backgroundImage: "url('https://static.babblegraph.com/assets/powered-by-stripe.svg')",
        height: "25px",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
        marginBottom: "10px",
    },
});

const StripeInput = (
     { component: Component, inputRef, ...props }
) => {
    const elementRef = useRef();
    useImperativeHandle(inputRef, () => ({
        // @ts-ignore
        focus: () => elementRef.current.focus
    }));
    return (
        <Component
            onReady={element => (elementRef.current = element)}
        {...props} />
    );
}

export type StripeError = {
    type: string;
    code: string;
    decline_code: string;
    message: string;
    param: string;
    payment_intent: string;
}

type GenericCardFormProps = {
    actionTitle: string;
    isLoading: boolean;
    elements: any;

    handleSubmit: (cardElement: typeof CardElement, cardholderName: string, postalCode: string) => void;
}

const GenericCardForm = (props: GenericCardFormProps) => {
    const [ cardholderName, setCardholderName ] = useState<string>("");
    const [ postalCode, setPostalCode ] = useState<string>("");

    const handleCardholderNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setCardholderName((event.target as HTMLInputElement).value);
    }
    const handlePostalCodeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPostalCode((event.target as HTMLInputElement).value);
    }
    const handleSubmit = () => {
        const cardElement = props.elements.getElement(CardNumberElement);
        props.handleSubmit(cardElement, cardholderName, postalCode);
    }

    const classes = styleClasses();
    return (
        <Form handleSubmit={handleSubmit}>
            <Grid container>
                <Grid item xs={12}>
                    <PrimaryTextField
                        className={classes.checkoutFormObject}
                        id="cardholder-name"
                        label="Cardholder Name"
                        variant="outlined"
                        defaultValue={cardholderName}
                        onChange={handleCardholderNameChange}
                        disabled={props.isLoading}
                        required
                        fullWidth />
                    </Grid>
                    <Grid item xs={12}>
                        <PrimaryTextField
                            className={classes.checkoutFormObject}
                            id="credit-card-number"
                            label="Credit Card Number"
                            variant="outlined"
                            InputLabelProps={{ shrink: true }}
                            InputProps={{
                                inputComponent: StripeInput,
                                inputProps: {
                                    component: CardNumberElement
                                },
                            }}
                            disabled={props.isLoading}
                            required
                            fullWidth />
                </Grid>
                <Grid item xs={6}>
                    <PrimaryTextField
                        className={classes.checkoutFormObject}
                        id="credit-card-expiration"
                        label="Expiration Date"
                        variant="outlined"
                        InputLabelProps={{ shrink: true }}
                        InputProps={{
                            inputComponent: StripeInput,
                            inputProps: {
                                component: CardExpiryElement
                            },
                        }}
                        disabled={props.isLoading}
                        required
                        fullWidth />
                </Grid>
                <Grid item xs={6}>
                    <PrimaryTextField
                        className={classes.checkoutFormObject}
                        id="credit-card-cvc"
                        label="CVC"
                        variant="outlined"
                        InputLabelProps={{ shrink: true }}
                        InputProps={{
                            inputComponent: StripeInput,
                            inputProps: {
                                component: CardCvcElement
                            },
                        }}
                        disabled={props.isLoading}
                        required
                        fullWidth />
                </Grid>
                <Grid item xs={12}>
                    <PrimaryTextField
                        id="zip"
                        className={classes.checkoutFormObject}
                        label="Postal Code"
                        variant="outlined"
                        defaultValue={postalCode}
                        onChange={handlePostalCodeChange}
                        disabled={props.isLoading}
                        required
                        fullWidth />
                </Grid>
                <Grid item xs={12}>
                    <PrimaryButton
                        type="submit"
                        className={classes.checkoutFormObject}
                        disabled={!postalCode || !cardholderName || props.isLoading}>
                        { props.actionTitle }
                    </PrimaryButton>
                </Grid>
                <Grid item xs={12}>
                    <a className={classes.stripeLink} href="https://stripe.com/" target="_blank">
                        <div className={classes.stripeBadge} />
                    </a>
                </Grid>
            </Grid>
        </Form>
    );
}

export default GenericCardForm;
