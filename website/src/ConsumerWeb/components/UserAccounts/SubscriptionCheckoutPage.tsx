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
import Link, { LinkTarget } from 'common/components/Link/Link';

import {
    Subscription,
    SubscriptionTrialInfo,
    PaymentState,
    createUserSubscription,
    CreateUserSubscriptionResponse,
    getActiveSubscriptionForUser,
    GetActiveSubscriptionForUserResponse,
    UpdateStripeSubscriptionForUserResponse,
    updateStripeSubscriptionForUser,
} from 'ConsumerWeb/api/stripe/subscription';
import {
    getManageTokenForPremiumCheckoutToken,
    GetManageTokenForPremiumCheckoutTokenResponse,
} from 'ConsumerWeb/api/token/premiumCheckoutToken';

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
    const [ isLoadingManagementToken, setIsLoadingManagementToken ] = useState<boolean>(true);


    const [ subscription, setSubscription ] = useState<Subscription | null>(null);
    const [ subscriptionTrialInfo, setSubscriptionTrialInfo ] = useState<SubscriptionTrialInfo | null>(null);
    const [ isLoadingCreateSubscription, setIsLoadingCreateSubscription ] = useState<boolean>(false);
    const [ error, setError ] = useState<Error>(null);

    const [ isPaymentConfirmationLoading, setIsPaymentConfirmationLoading ] = useState<boolean>(false);
    const [ successType, setSuccessType ] = useState<PaymentType | null>(null);
    const [ paymentError, setPaymentError ] = useState<string | null>(null);

    const [ subscriptionManagementToken, setSubscriptionManagementToken ] = useState<string | null>(null);

    useEffect(() => {
        getActiveSubscriptionForUser({},
        (resp: GetActiveSubscriptionForUserResponse) => {
            if (!!resp.subscription) {
                setIsLoadingUserSubscription(false);
                setSubscription(resp.subscription);
                setSubscriptionTrialInfo(resp.subscription.trialInfo);
            } else {
                createUserSubscription({},
                    (resp: CreateUserSubscriptionResponse) => {
                        setIsLoadingUserSubscription(false);
                        if (!!resp.subscription) {
                            setSubscription(resp.subscription);
                            setSubscriptionTrialInfo(resp.subscription.trialInfo);
                        } else {
                            setError(new Error("Something went wrong with your request"));
                        }
                },
                (err: Error) => {
                    setIsLoadingUserSubscription(false);
                    setError(err);
                });
            }
        },
        (err: Error) => {
            setIsLoadingUserSubscription(false);
            setError(err);
        });
        getManageTokenForPremiumCheckoutToken({
            token: token,
        },
        (resp: GetManageTokenForPremiumCheckoutTokenResponse) => {
            setIsLoadingManagementToken(false);
            setSubscriptionManagementToken(resp.token);
        },
        (err: Error) => {
            setIsLoadingManagementToken(false);
            setError(err);
        });
    }, []);

    const handleSubmit = () => {
        setIsLoadingCreateSubscription(true);
    }
    const isSubscriptionAlreadySetup = (
        subscription &&
        subscription.paymentState !== PaymentState.CreatedUnpaid &&
        subscription.paymentState !== PaymentState.TrialNoPaymentMethod
    );
    const isPageLoading = isLoadingUserSubscription || isLoadingManagementToken;
    const isSelectorLoading = isLoadingCreateSubscription || isPaymentConfirmationLoading;
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
                <Link href={`/manage/${subscriptionManagementToken}/payment-settings`} target={LinkTarget.Self}>
                    Click here to modify your subscription
                </Link>
            </div>
        );
    } else if (!!successType) {
        const headerMessage = successType === PaymentType.Setup ? (
            "Your payment method was successfully setup."
        ) : (
            "Your payment was successful!"
        );
        const bodyMessage = successType === PaymentType.Setup ? (
            "You will be automatically charged at the end of your trial! You will get an email before you’re charged. Your subscription will activate in a few minutes."
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
                <Link href={`/manage/${subscriptionManagementToken}`} target={LinkTarget.Self}>
                    Click here to your settings
                </Link>
            </div>
        );
    } else if (!!error) {
        body = (
            <Heading3 color={TypographyColor.Primary}>
                Something went wrong processing your request. You have not been charged. Try again later, or reach out to hello@babblegraph.com
            </Heading3>
        );
    } else {
        const shouldShowCheckoutForm = !!subscription;
        const isEligibleForTrial = !subscriptionTrialInfo || !!subscriptionTrialInfo.trialEligibilityDays;
        const currentPeriodEndDate = !subscription ? null : new Date(subscription.currentPeriodEnd);
        body = (
            <div>
                <Heading1 color={TypographyColor.Primary}>
                    Babblegraph Premium Subscription Checkout
                </Heading1>
                {
                    isEligibleForTrial ? (
                        <Paragraph color={TypographyColor.Confirmation}>
                            Your free trial of Babblegraph Premium has begun!
                        </Paragraph>
                    ) : (
                        <Paragraph color={TypographyColor.Warning}>
                            According to our records, you are not eligible for the 14 day free trial of Babblegraph Premium. If you believe this is an error, reach out via email at hello@babblegraph.com
                        </Paragraph>
                    )
                }
                {
                    shouldShowCheckoutForm && (
                        subscription.paymentIntentClientSecret ? (
                            <div>
                                <Paragraph>
                                    Enter your payment details below. You will be charged $20 USD. Your subscription to premium will not begin until your payment is completed.
                                </Paragraph>
                                <CollectPaymentForm
                                    paymentIntentClientSecret={subscription.paymentIntentClientSecret}
                                    handleIsStripeRequestLoading={setIsPaymentConfirmationLoading}
                                    handleSuccess={() => setSuccessType(PaymentType.Payment)}
                                    handleFailure={setPaymentError}
                                    handleError={setError} />
                            </div>
                        ) : (
                            <div>
                                <Heading3 color={TypographyColor.Primary}>
                                    Add a payment method to your account
                                </Heading3>
                                <Paragraph>
                                    {
                                        !!currentPeriodEndDate ? (
                                            `You will be charged $20 USD on ${currentPeriodEndDate.toLocaleDateString()}, the end of your trial. You will not be charged before then.`
                                        ) : (
                                            `You will be charged $20 USD at the end of your trial. You will not be charged before then.`
                                        )
                                    }
                                </Paragraph>
                                <Paragraph>
                                    If you do not add a payment method now, your subscription will automatically be canceled at the end of the trial.
                                </Paragraph>
                                <AddPaymentMethodForm
                                    handleIsStripeRequestLoading={setIsPaymentConfirmationLoading}
                                    handleSuccess={(paymentMethodID: string) => setSuccessType(PaymentType.Setup)}
                                    handleFailure={setPaymentError}
                                    handleError={setError}
                                    isDefault />
                            </div>
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

export default SubscriptionCheckoutPage;
