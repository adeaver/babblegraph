import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import AddPaymentMethodForm from 'common/components/Stripe/AddPaymentMethodForm';
import CollectPaymentForm from 'common/components/Stripe/CollectPaymentForm';

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
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PrimaryRadio } from 'common/components/Radio/Radio';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import StripeInput from 'common/util/stripe/StripeInput';

import {
    Subscription,
    SubscriptionType,
    PaymentState,
    createUserSubscription,
    CreateUserSubscriptionResponse,
    getActiveSubscriptionForUser,
    GetActiveSubscriptionForUserResponse,
    UpdateStripeSubscriptionFrequencyForUserResponse,
    updateStripeSubscriptionFrequencyForUser,
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
    checkIcon: {
        color: Color.Confirmation,
        display: 'block',
        margin: '0 auto',
        fontSize: '48px',
    },
    checkoutFormObject: {
        width: "100%",
        margin: "10px 0",
        paddingLeft: "5px",
        paddingRight: "5px",
        boxSizing: "border-box",
    },
})

enum PaymentType {
    Setup = 'setup',
    Payment = 'payment',
}

type Params = {
    token: string
}

type SubscriptionCheckoutPageProps = RouteComponentProps<Params>

const SubscriptionCheckoutPage = (props: SubscriptionCheckoutPageProps) => {
    const { token } = props.match.params;

    const [ isLoadingUserSubscription, setIsLoadingUserSubscription ] = useState<boolean>(true);

    const [ isLoadingUpdateSubscription, setIsLoadingUpdateSubscription ] = useState<boolean>(false);
    const [ subscriptionType, setSubscriptionType ] = useState<SubscriptionType>(SubscriptionType.Monthly);

    const [ subscription, setSubscription ] = useState<Subscription | null>(null);
    const [ isLoadingCreateSubscription, setIsLoadingCreateSubscription ] = useState<boolean>(false);
    const [ error, setError ] = useState<Error>(null);

    const [ isPaymentConfirmationLoading, setIsPaymentConfirmationLoading ] = useState<boolean>(false);
    const [ successType, setSuccessType ] = useState<PaymentType | null>(null);
    const [ paymentError, setPaymentError ] = useState<string | null>(null);

    useEffect(() => {
        getActiveSubscriptionForUser({},
        (resp: GetActiveSubscriptionForUserResponse) => {
            setIsLoadingUserSubscription(false);
            !!resp.subscription && setSubscription(resp.subscription);
        },
        (err: Error) => {
            setIsLoadingUserSubscription(false);
            setError(err);
        });
    }, []);

    const handleUpdateSubscription = () => {
        setIsLoadingUpdateSubscription(true);
        updateStripeSubscriptionFrequencyForUser({
            isYearlySubscription: subscriptionType === "yearly",
            // stripeSubscriptionId: stripeSubscriptionID,
        },
        (resp: UpdateStripeSubscriptionFrequencyForUserResponse)  => {
            setIsLoadingUpdateSubscription(false);
        },
        (err: Error) => {
            setIsLoadingUpdateSubscription(false);
            setError(err);
        });
    }
    const handleSubmit = () => {
        setIsLoadingCreateSubscription(true);
        createUserSubscription({
            subscriptionCreationToken: token,
            isYearlySubscription: subscriptionType === "yearly",
        },
        (resp: CreateUserSubscriptionResponse) => {
            setIsLoadingCreateSubscription(false);
        },
        (err: Error) => {
            setIsLoadingCreateSubscription(false);
            setError(err);
        });
    }
    const isSubscriptionAlreadySetup = (
        subscription &&
        subscription.paymentState !== PaymentState.CreatedUnpaid &&
        subscription.paymentState !== PaymentState.TrialNoPaymentMethod
    );
    const isPageLoading = isLoadingUserSubscription;
    const isSelectorLoading = isLoadingCreateSubscription || isPaymentConfirmationLoading || isLoadingUpdateSubscription;
    const classes = styleClasses();
    let body;
    if (isPageLoading) {
        body = <LoadingSpinner />;
    } else if (isSubscriptionAlreadySetup) {
        body = (
            <div>
                <VerifiedUserIcon className={classes.checkIcon} />
                <Heading1 color={TypographyColor.Primary}>
                    Your subscription was already setup!
                </Heading1>
                <Paragraph>
                    If you need to modify it, you can go to the subscription settings page! Thanks again for subscribing to Babblegraph.
                </Paragraph>
            </div>
        );
    } else if (!!successType) {
        const headerMessage = successType === PaymentType.Setup ? (
            "Your payment method was successfully setup."
        ) : (
            "Your payment was successful!"
        );
        const bodyMessage = successType === PaymentType.Setup ? (
            "You will be automatically charged at the end of your trial! You will get an email before you’re charged."
        ) : (
            "Your subscription is currently processing. Your premium subscription will become active in the next 10 minutes. You’ll receive an email when it's active!"
        )
        body = (
            <div>
                <VerifiedUserIcon className={classes.checkIcon} />
                <Heading1 color={TypographyColor.Primary}>
                    { headerMessage }
                </Heading1>
                <Paragraph>
                    { bodyMessage }
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
        const shouldShowCheckoutForm = !!subscription && !isLoadingUpdateSubscription;
        body = (
            <div>
                <SubscriptionSelector
                    subscriptionType={subscriptionType}
                    isPaymentConfirmationLoading={isPaymentConfirmationLoading}
                    isLoadingUpdateSubscription={isLoadingUpdateSubscription}
                    isCheckoutFormVisible={shouldShowCheckoutForm}
                    isEligibleForTrial={!!subscription.trialInfo.trialEligibilityDays}
                    handleUpdateSubscriptionType={setSubscriptionType}
                    handleUpdateSubscription={handleUpdateSubscription}
                    handleSubmit={handleSubmit} />
                {
                    shouldShowCheckoutForm && (
                        subscription.paymentIntentClientSecret ? (
                            <CollectPaymentForm
                                paymentIntentClientSecret={subscription.paymentIntentClientSecret}
                                handleIsStripeRequestLoading={setIsPaymentConfirmationLoading}
                                handleSuccess={() => setSuccessType(PaymentType.Payment)}
                                handleFailure={setPaymentError}
                                handleError={setError} />
                        ) : (
                            <AddPaymentMethodForm
                                handleIsStripeRequestLoading={setIsPaymentConfirmationLoading}
                                handleSuccess={(paymentMethodID: string) => setSuccessType(PaymentType.Setup)}
                                handleFailure={setPaymentError}
                                handleError={setError}
                                isDefault />
                        )
                    )
                }
                {
                    isSelectorLoading && <LoadingSpinner />
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
    isLoadingUpdateSubscription: boolean;
    isEligibleForTrial: boolean;

    handleUpdateSubscription: () => void;
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
            {
                props.isEligibleForTrial ? (
                    <Paragraph color={TypographyColor.Confirmation}>
                        You are eligible for the free trial of Babblegraph Premium
                    </Paragraph>
                ) : (
                    <Paragraph color={TypographyColor.Warning}>
                        According to our records, you are not eligible for the 14 day free trial of Babblegraph Premium. If you believe this is an error, reach out via email at hello@babblegraph.com
                    </Paragraph>
                )
            }
            <FormControl className={classes.subscriptionSelector} component="fieldset">
                <RadioGroup aria-label="subscription-type" name="subscription-type1" value={props.subscriptionType} onChange={handleRadioFormChange}>
                    <Grid container
                        className={classes.checkoutFormObject}>
                        <Grid item xs={false} md={3}>
                            &nbsp;
                        </Grid>
                        <Grid item className={classes.subscriptionOption} xs={12} md={3}>
                            <FormControlLabel value="monthly" control={<PrimaryRadio />} label="Monthly ($3/month)" />
                        </Grid>
                        <Grid item className={classes.subscriptionOption} xs={12} md={3}>
                            <FormControlLabel value="yearly" control={<PrimaryRadio />} label="Yearly ($34/year)" />
                        </Grid>
                    </Grid>
                </RadioGroup>
            </FormControl>
            {
                (props.isEligibleForTrial && !props.isCheckoutFormVisible) && (
                    <Paragraph size={Size.Small}>
                        You can change this later!
                    </Paragraph>
                )
            }
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item className={classes.submitButton} xs={12} md={6}>
                    {
                        !props.isCheckoutFormVisible ? (
                            <PrimaryButton
                                className={classes.checkoutFormObject}
                                disabled={props.isPaymentConfirmationLoading || props.isLoadingUpdateSubscription}
                                onClick={props.handleSubmit}>
                                { props.isEligibleForTrial ? "Confirm and start your trial" : "Confirm your selection" }
                            </PrimaryButton>
                        ) : (
                            <PrimaryButton
                                className={classes.checkoutFormObject}
                                disabled={props.isPaymentConfirmationLoading || props.isLoadingUpdateSubscription}
                                onClick={props.handleUpdateSubscription}>
                                Update Subscription Selection
                            </PrimaryButton>
                        )
                    }
                </Grid>
            </Grid>
        </div>
    );
}

export default SubscriptionCheckoutPage;
