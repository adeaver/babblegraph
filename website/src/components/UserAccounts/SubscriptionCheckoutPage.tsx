import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import {
    Elements,
    CardElement,
    ElementsConsumer,
} from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";

import { makeStyles } from '@material-ui/core/styles';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import RadioGroup from '@material-ui/core/RadioGroup';

import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PrimaryRadio } from 'common/components/Radio/Radio';
import { PrimaryButton } from 'common/components/Button/Button';

import {
    getOrCreateUserSubscription,
    GetOrCreateUserSubscriptionResponse
} from 'api/useraccounts/useraccounts';

const styleClasses = makeStyles({
    subscriptionSelector: {
        width: "100%",
    },
    subscriptionOption: {
        display: "flex",
        justifyContent: "center",
    },
    submitButton: {
        display: "flex",
        justifyContent: "center",
    },
})

declare const window: any;
const stripePromise = loadStripe(window.initialData["stripe_public_key"]);

type Params = {
    token: string
}

type SubscriptionCheckoutPageProps = RouteComponentProps<Params>

const SubscriptionCheckoutPage = (props: SubscriptionCheckoutPageProps) => {
    const { token } = props.match.params;

    const [ subscriptionType, setSubscriptionType ] = useState<string>("monthly");

    const [ stripeSubscriptionID, setStripeSubscriptionID ] = useState<string | null>(null);
    const [ stripeClientSecret, setStripeClientSecret ] = useState<string | null>(null);
    const [ stripePaymentState, setStripePaymentState ] = useState<number | null>(null);
    const [ isLoadingCreateSubscription, setIsLoadingCreateSubscription ] = useState<boolean>(false);
    const [ error, setError ] = useState<Error>(null);

    const handleSubmit = () => {
        setIsLoadingCreateSubscription(true);
        getOrCreateUserSubscription({
            subscriptionCreationToken: token,
            isYearlySubscription: subscriptionType === "yearly",
        },
        (resp: GetOrCreateUserSubscriptionResponse) => {
            setIsLoadingCreateSubscription(false);
            setStripeSubscriptionID(resp.stripeSubscriptionId);
            setStripeClientSecret(resp.stripeClientSecret);
            setStripePaymentState(resp.stripePaymentState);
        },
        (err: Error) => {
            setIsLoadingCreateSubscription(false);
            setError(err);
        });
    }

    const classes = styleClasses();
    const isLoading = isLoadingCreateSubscription;
    let body;
    if (isLoading) {
        body = <LoadingSpinner />;
    } else if (!!error) {
        body = (
            <Heading3 color={TypographyColor.Primary}>
                Something went wrong processing your request. You have not been charged. Try again later, or reach out to hello@babblegraph.com
            </Heading3>
        );
    } else if (stripeClientSecret != null && !!stripeSubscriptionID) {
        body = (
            // Unclear why, but Elements doesn't think it has children
            // @ts-ignore
            <Elements stripe={stripePromise}>
                <InjectedSubscriptionCheckoutFormProps
                    stripeClientSecret={stripeClientSecret}
                    stripeSubscriptionID={stripeSubscriptionID}
                    handleSuccessfulPayment={() => console.log("successful")}
                    handlePaymentError={(msg: string) => console.log(msg)} />
            </Elements>
        );
    } else {
        body = (
            <SubscriptionSelector
                subscriptionType={subscriptionType}
                handleUpdateSubscriptionType={setSubscriptionType}
                handleSubmit={handleSubmit} />
        );
    }
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        { body }
                    </DisplayCard>
                </Grid>
            </Grid>
        </Page>
    );
}

type SubscriptionSelectorProps = {
    subscriptionType: string;

    handleUpdateSubscriptionType: (string) => void;
    handleSubmit: () => void;
}

const SubscriptionSelector = (props: SubscriptionSelectorProps) => {
    const handleRadioFormChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.handleUpdateSubscriptionType((event.target as HTMLInputElement).value);
    };

    const classes = styleClasses();
    return (
        <div>
            <Heading1 color={TypographyColor.Primary}>
                Babblegraph Premium Subscription Checkout
            </Heading1>
            <Heading3>
                Choose your subscription
            </Heading3>
            <FormControl className={classes.subscriptionSelector} component="fieldset">
                <RadioGroup aria-label="subscription-type" name="subscription-type1" value={props.subscriptionType} onChange={handleRadioFormChange}>
                    <Grid container>
                        <Grid item xs={false} md={3}>
                            &nbsp;
                        </Grid>
                        <Grid item className={classes.subscriptionOption} xs={12} md={3}>
                            <FormControlLabel value="monthly" control={<PrimaryRadio />} label="Monthly" />
                        </Grid>
                        <Grid item className={classes.subscriptionOption} xs={12} md={3}>
                            <FormControlLabel value="yearly" control={<PrimaryRadio />} label="Yearly" />
                        </Grid>
                    </Grid>
                </RadioGroup>
            </FormControl>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item className={classes.submitButton} xs={12} md={6}>
                    <PrimaryButton onClick={props.handleSubmit}>
                        Continue to Payment
                    </PrimaryButton>
                </Grid>
            </Grid>
        </div>
    );
}

type StripeResult = {
    error?: StripeError | null;
    setupIntent?: any; // TODO(fill this in?)
}

type StripeError = {
    type: string;
    code: string;
    decline_code: string;
    message: string;
    param: string;
    payment_intent: string;
}

type InjectedSubscriptionCheckoutFormProps = {
    stripeClientSecret: string;
    stripeSubscriptionID: string;

    handleSuccessfulPayment: () => void;
    handlePaymentError: (msg: string) => void;
}

const InjectedSubscriptionCheckoutFormProps = (props: InjectedSubscriptionCheckoutFormProps) => (
    <ElementsConsumer>
        {({stripe, elements}) => (
            <SubscriptionCheckoutForm stripe={stripe} elements={elements} {...props} />
        )}
    </ElementsConsumer>
);

type SubscriptionCheckoutFormProps = InjectedSubscriptionCheckoutFormProps & {
    stripe: any;
    elements: any;
}


const SubscriptionCheckoutForm = (props: SubscriptionCheckoutFormProps) => {
    const [ isPaymentConfirmationLoading, setIsPaymentConfirmationLoading ] = useState<boolean>(false);

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setIsPaymentConfirmationLoading(true);
        const cardElement = props.elements.getElement(CardElement);
        props.stripe.confirmCardSetup(props.stripeClientSecret, {
            payment_method: {
                card: cardElement,
            }
        }).then((result: StripeResult) => {
            setIsPaymentConfirmationLoading(false);
            if (!!result.error) {
                props.handlePaymentError(result.error.message);
            } else if (!!result.setupIntent && result.setupIntent.status === "succeeded") {
                props.handleSuccessfulPayment();
            } else {
                props.handlePaymentError(result.setupIntent.status);
            }
        }).catch((err: Error) => {
            props.handlePaymentError(err.message);
        });
    }

    return (
        <form onSubmit={handleSubmit} noValidate autoComplete="off">
            <CardElement />
            {
                isPaymentConfirmationLoading ? (
                    <LoadingSpinner />
                ) : (
                    <PrimaryButton type="submit">Pay</PrimaryButton>
                )
            }
        </form>
    );
}

export default SubscriptionCheckoutPage;
