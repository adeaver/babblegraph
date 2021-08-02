import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import {
    Elements,
    CardNumberElement,
    CardExpiryElement,
    CardCvcElement,
    ElementsConsumer,
} from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";

import { makeStyles } from '@material-ui/core/styles';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import RadioGroup from '@material-ui/core/RadioGroup';
import VerifiedUserIcon from '@material-ui/icons/VerifiedUser';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import Color from 'common/styles/colors';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PrimaryRadio } from 'common/components/Radio/Radio';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import StripeInput from 'common/util/stripe/StripeInput';

import {
    getOrCreateUserSubscription,
    GetOrCreateUserSubscriptionResponse,
    getUserNonTerminatedStripeSubscription,
    GetUserNonTerminatedStripeSubscriptionResponse,
    DeleteStripeSubscriptionForUserResponse,
    deleteStripeSubscriptionForUser
} from 'api/stripe/subscription';

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
    checkoutFormObject: {
        width: "100%",
        margin: "10px 0",
        paddingLeft: "5px",
        paddingRight: "5px",
        boxSizing: "border-box",
    },
    checkIcon: {
        color: Color.Confirmation,
        display: 'block',
        margin: '0 auto',
        fontSize: '48px',
    }
})

const stripeElementsOptions = {
    style: {
        base: {
            fontSize: '16px',
            fontFamily: "'Roboto', sans-serif",
            color: Color.TextGray,
            '::placeholder': {
                color: Color.TextGray
            }
        },
        invalid: {
            color: Color.Warning,
        },
    }
}

declare const window: any;
const stripePromise = loadStripe(window.initialData["stripe_public_key"]);

type Params = {
    token: string
}

type SubscriptionCheckoutPageProps = RouteComponentProps<Params>

const SubscriptionCheckoutPage = (props: SubscriptionCheckoutPageProps) => {
    const { token } = props.match.params;

    const [ isLoadingUserSubscription, setIsLoadingUserSubscription ] = useState<boolean>(true);

    const [ subscriptionType, setSubscriptionType ] = useState<string>("monthly");

    const [ stripeSubscriptionID, setStripeSubscriptionID ] = useState<string | null>(null);
    const [ stripeClientSecret, setStripeClientSecret ] = useState<string | null>(null);
    const [ stripePaymentState, setStripePaymentState ] = useState<number | null>(null);
    const [ isLoadingCreateSubscription, setIsLoadingCreateSubscription ] = useState<boolean>(false);
    const [ error, setError ] = useState<Error>(null);

    const [ isPaymentConfirmationLoading, setIsPaymentConfirmationLoading ] = useState<boolean>(false);
    const [ wasPaymentSuccessful, setWasPaymentSuccessful ] = useState<boolean>(false);
    const [ paymentError, setPaymentError ] = useState<string | null>(null);

    useEffect(() => {
        getUserNonTerminatedStripeSubscription({
            subscriptionCreationToken: token,
        },
        (resp: GetUserNonTerminatedStripeSubscriptionResponse) => {
            setIsLoadingUserSubscription(false);
            console.log(resp);
            resp.stripeSubscriptionId != null && setStripeSubscriptionID(resp.stripeSubscriptionId);
            resp.stripeClientSecret != null && setStripeClientSecret(resp.stripeClientSecret);
            resp.stripePaymentState != null && setStripePaymentState(resp.stripePaymentState);
            !!resp.isYearlySubscription && setSubscriptionType("yearly");
        },
        (err: Error) => {
            setIsLoadingUserSubscription(false);
            setError(err);
        });
    }, []);

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

    const isLoading = isLoadingUserSubscription;
    const classes = styleClasses();
    let body;
    if (isLoading) {
        body = <LoadingSpinner />;
    } else if (wasPaymentSuccessful) {
        body = (
            <div>
                <VerifiedUserIcon className={classes.checkIcon} />
                <Heading1 color={TypographyColor.Primary}>
                    Your payment method was successfully setup.
                </Heading1>
                <Paragraph>
                    You will be automatically charged at the end of your trial! You will get an email before you’re charged.
                </Paragraph>
            </div>
        );
    } else if (!!error) {
        body = (
            <Heading3 color={TypographyColor.Primary}>
                Something went wrong processing your request. You have not been charged. Try again later, or reach out to hello@babblegraph.com
            </Heading3>
        );
    } else {
        const shouldShowCheckoutForm = stripeClientSecret != null && !!stripeSubscriptionID;
        console.log(stripeSubscriptionID);
        console.log(stripeClientSecret);
        body = (
            <div>
                <SubscriptionSelector
                    subscriptionType={subscriptionType}
                    isCheckoutFormVisible={shouldShowCheckoutForm}
                    isPaymentConfirmationLoading={isPaymentConfirmationLoading}
                    handleUpdateSubscriptionType={setSubscriptionType}
                    handleSubmit={handleSubmit} />
                {
                    shouldShowCheckoutForm && (
                        // Unclear why, but Elements doesn't think it has children
                        // @ts-ignore
                        <Elements stripe={stripePromise}>
                            <InjectedSubscriptionCheckoutFormProps
                                stripeClientSecret={stripeClientSecret}
                                stripeSubscriptionID={stripeSubscriptionID}
                                isPaymentConfirmationLoading={isPaymentConfirmationLoading}
                                setIsPaymentConfirmationLoading={setIsPaymentConfirmationLoading}
                                handleSuccessfulPayment={() => setWasPaymentSuccessful(true)}
                                handlePaymentError={setPaymentError}
                                handleGenericRequestError={setError} />
                        </Elements>
                    )
                }
                {
                    (isLoadingCreateSubscription || isPaymentConfirmationLoading) && <LoadingSpinner />
                }
            </div>
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
                    <Snackbar open={paymentError != null} autoHideDuration={6000} onClose={() => setPaymentError(null)}>
                        <Alert severity="error">{paymentError}</Alert>
                    </Snackbar>
                </Grid>
            </Grid>
        </Page>
    );
}

type SubscriptionSelectorProps = {
    subscriptionType: string;
    isCheckoutFormVisible: boolean;
    isPaymentConfirmationLoading: boolean;

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
                    <Grid container
                        className={classes.checkoutFormObject}>
                        <Grid item xs={false} md={3}>
                            &nbsp;
                        </Grid>
                        <Grid item className={classes.subscriptionOption} xs={12} md={3}>
                            <FormControlLabel value="monthly" control={<PrimaryRadio disabled={props.isCheckoutFormVisible} />} label="Monthly ($3/month)" />
                        </Grid>
                        <Grid item className={classes.subscriptionOption} xs={12} md={3}>
                            <FormControlLabel value="yearly" control={<PrimaryRadio disabled={props.isCheckoutFormVisible} />} label="Yearly ($34/year)" />
                        </Grid>
                    </Grid>
                </RadioGroup>
            </FormControl>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item className={classes.submitButton} xs={12} md={6}>
                    {
                        !props.isCheckoutFormVisible ? (
                            <PrimaryButton
                                className={classes.checkoutFormObject}
                                disabled={props.isPaymentConfirmationLoading}
                                onClick={props.handleSubmit}>
                                Confirm Selection
                            </PrimaryButton>
                        ) : (
                            <PrimaryButton
                                className={classes.checkoutFormObject}
                                disabled={props.isPaymentConfirmationLoading}
                                onClick={props.handleSubmit}>
                                Update Subscription Selection
                            </PrimaryButton>
                        )
                    }
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
    isPaymentConfirmationLoading: boolean;

    setIsPaymentConfirmationLoading: (isLoading: boolean) => void;
    handleSuccessfulPayment: () => void;
    handlePaymentError: (msg: string) => void;
    handleGenericRequestError: (err: Error) => void;
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
    const [ cardholderName, setCardholderName ] = useState<string>("");
    const [ postalCode, setPostalCode ] = useState<string>("");

    const handleCardholderNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setCardholderName((event.target as HTMLInputElement).value);
    }
    const handlePostalCodeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPostalCode((event.target as HTMLInputElement).value);
    }

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        props.setIsPaymentConfirmationLoading(true);
        const cardElement = props.elements.getElement(CardNumberElement);
        // TODO: handle case in which this isn't a setup card
        props.stripe.confirmCardSetup(props.stripeClientSecret, {
            payment_method: {
                card: cardElement,
                billing_details: {
                    name: cardholderName,
                }
            }
        }).then((result: StripeResult) => {
            props.setIsPaymentConfirmationLoading(false);
            if (!!result.error) {
                props.handlePaymentError(result.error.message);
            } else if (!!result.setupIntent && result.setupIntent.status === "succeeded") {
                props.handleSuccessfulPayment();
            } else {
                props.handlePaymentError(result.setupIntent.status);
            }
            console.log(result);
        }).catch((err: Error) => {
            props.handleGenericRequestError(err)
        });
    }

    const classes = styleClasses();
    return (
        <form onSubmit={handleSubmit} noValidate autoComplete="off">
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <PrimaryTextField
                        className={classes.checkoutFormObject}
                        id="cardholder-name"
                        label="Cardholder Name"
                        variant="outlined"
                        defaultValue={cardholderName}
                        onChange={handleCardholderNameChange}
                        disabled={props.isPaymentConfirmationLoading}
                        required
                        fullWidth />
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
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
                        disabled={props.isPaymentConfirmationLoading}
                        required
                        fullWidth />

                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={6} md={3}>
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
                        disabled={props.isPaymentConfirmationLoading}
                        required
                        fullWidth />

                </Grid>
                <Grid item xs={6} md={3}>
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
                        disabled={props.isPaymentConfirmationLoading}
                        required
                        fullWidth />

                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <PrimaryTextField
                        id="zip"
                        className={classes.checkoutFormObject}
                        label="Postal Code"
                        variant="outlined"
                        defaultValue={postalCode}
                        onChange={handlePostalCodeChange}
                        disabled={props.isPaymentConfirmationLoading}
                        required
                        fullWidth />

                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <PrimaryButton
                        type="submit"
                        className={classes.checkoutFormObject}
                        disabled={!postalCode || !cardholderName || props.isPaymentConfirmationLoading}>
                        Pay
                    </PrimaryButton>
                </Grid>
            </Grid>
        </form>
    );
}

export default SubscriptionCheckoutPage;
