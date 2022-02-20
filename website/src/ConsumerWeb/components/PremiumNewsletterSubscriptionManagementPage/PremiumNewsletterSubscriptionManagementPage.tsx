import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Link, { LinkTarget } from 'common/components/Link/Link';
import Paragraph, { Size } from 'common/typography/Paragraph';
import PaymentMethodDisplay from 'ConsumerWeb/components/common/Billing/PaymentMethodDisplay';
import { PrimaryButton, WarningButton } from 'common/components/Button/Button';
import Alert from 'common/components/Alert/Alert';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';
import {
    withUserProfileInformation,
    UserProfileComponentProps
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import PremiumNewsletterSubscriptionCardForm from 'ConsumerWeb/components/common/Billing/PremiumNewsletterSubscriptionCheckoutForm';

import {
    PaymentMethod,
    PremiumNewsletterSubscription,
    PaymentState,

    LookupActivePremiumNewsletterSubscriptionResponse,
    lookupActivePremiumNewsletterSubscription,

    GetPaymentMethodsForUserResponse,
    getPaymentMethodsForUser,

    DeletePaymentMethodForUserResponse,
    deletePaymentMethodForUser,

    MarkPaymentMethodAsDefaultResponse,
    markPaymentMethodAsDefault,
} from 'ConsumerWeb/api/billing/billing';

const styleClasses = makeStyles({
    paymentMethodDisplayContainer: {
        padding: '5px',
    },
    paymentMethodButton: {
        width: '100%',
        margin: '10px 0',
    },
});

type Params = {
    token: string;
}

type PremiumNewsletterSubscriptionManagementPageProps = RouteComponentProps<Params>;

const PremiumNewsletterSubscriptionManagementPage = withUserProfileInformation<PremiumNewsletterSubscriptionManagementPageProps>(
    RouteEncryptionKey.SubscriptionManagement,
    [],
    (ownProps: PremiumNewsletterSubscriptionManagementPageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.PaymentSettings,
    (props: PremiumNewsletterSubscriptionManagementPageProps & UserProfileComponentProps) => {
        const { token } = props.match.params;

        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Premium Subscription and Payment Settings"
                        backArrowDestination={`/manage/${token}`} />
                    <PremiumSubscriptionManagementComponent
                        subscriptionManagementToken={token} />
                    <PaymentMethodManagementComponent />
                </DisplayCard>
            </CenteredComponent>
        )
    }
);

type PremiumSubscriptionManagementComponentOwnProps = {
    subscriptionManagementToken: string;
}

const PremiumSubscriptionManagementComponent = asBaseComponent<LookupActivePremiumNewsletterSubscriptionResponse, PremiumSubscriptionManagementComponentOwnProps>(
    (props: LookupActivePremiumNewsletterSubscriptionResponse & PremiumSubscriptionManagementComponentOwnProps & BaseComponentProps) => {
        if (!props.premiumNewsletterSubscription) {
            return (
                <div>
                    <Heading3 color={TypographyColor.Primary}>
                        You don’t have an active Babblegraph Premium subscription
                    </Heading3>
                    <Link href={`/manage/${props.subscriptionManagementToken}/premium`} target={LinkTarget.Self}>
                        If you’d like to start or restart one, click here.
                    </Link>
                </div>
            );
        }
        const { paymentState } = props.premiumNewsletterSubscription;
        switch (paymentState) {
            case PaymentState.CreatedUnpaid:
                return (
                    <div>
                        <Heading3 color={TypographyColor.Primary}>
                            It looks like you recently tried to start a Babblegraph Premium subscription...
                        </Heading3>
                        <Paragraph>
                            You currently haven’t paid for your subscription, so it’s not active. You can complete the payment using the form below.
                        </Paragraph>
                        <Paragraph>
                            If you didn’t try to start a subscription or no longer want one, then no action is required on your part. Our system should update within a couple of days.
                        </Paragraph>
                        <PremiumNewsletterSubscriptionCardForm
                            premiumNewsletterSusbcription={props.premiumNewsletterSubscription} />
                    </div>
                );
            case PaymentState.TrialNoPaymentMethod:
                return (
                    <div>
                        <Heading3 color={TypographyColor.Primary}>
                            You’re currently trialing Babblegraph Premium until {new Date(props.premiumNewsletterSubscription.currentPeriodEnd).toLocaleDateString()}
                        </Heading3>
                        <Paragraph>
                            However, you haven’t added a payment method, so you are set to lose Babblegraph Premium features on that date. You can add a payment method with the form below.
                        </Paragraph>
                        <PremiumNewsletterSubscriptionCardForm
                            premiumNewsletterSusbcription={props.premiumNewsletterSubscription} />
                    </div>
                );
            case PaymentState.TrialPaymentMethodAdded:
                return (
                    <div>
                        <Heading3 color={TypographyColor.Primary}>
                            You’re currently trialing Babblegraph Premium until {new Date(props.premiumNewsletterSubscription.currentPeriodEnd).toLocaleDateString()}
                        </Heading3>
                        { /* TODO: add auto-renew information here */ }
                        <Paragraph>
                            You will be charged on that date. You can add a new payment method below.
                        </Paragraph>
                        <PremiumNewsletterSubscriptionCardForm
                            premiumNewsletterSusbcription={props.premiumNewsletterSubscription} />
                    </div>
                );
            case PaymentState.Active:
                return (
                    <div>
                        <Heading3 color={TypographyColor.Primary}>
                            You’re currently on Babblegraph Premium
                        </Heading3>
                        { /* TODO: add auto-renew information here */ }
                        <Paragraph>
                            You will be charged next on {new Date(props.premiumNewsletterSubscription.currentPeriodEnd).toLocaleDateString()} You can add a new payment method below.
                        </Paragraph>
                        <PremiumNewsletterSubscriptionCardForm
                            premiumNewsletterSusbcription={props.premiumNewsletterSubscription} />
                    </div>
                );
            case PaymentState.Errored:
                return (
                    <div>
                        <Heading3 color={TypographyColor.Primary}>
                            We’ve encountered an error processing payment for your Babblegraph Premium subscription
                        </Heading3>
                        <Paragraph>
                            If you’d like to continue using Babblegraph Premium, please designate a new default payment method. If you would not like to continue, then no action is required.
                        </Paragraph>
                        <Paragraph>
                            It may take a day or two for this issue to be resolved with your account.
                        </Paragraph>
                        <PremiumNewsletterSubscriptionCardForm
                            premiumNewsletterSusbcription={props.premiumNewsletterSubscription} />
                    </div>
                );
            case PaymentState.Terminated:
                throw new Error("Active subscription endpoint returned terminated subscription");
            default:
                throw new Error(`Unrecognized payment state ${paymentState}`);
        }
    },
    (
        ownProps: PremiumSubscriptionManagementComponentOwnProps,
        onSuccess: (resp: LookupActivePremiumNewsletterSubscriptionResponse) => void,
        onError: (err: Error) => void,
    ) => lookupActivePremiumNewsletterSubscription({subscriptionManagementToken: ownProps.subscriptionManagementToken}, onSuccess, onError),
    false
);

type PaymentMethodManagementComponentOwnProps = {}

const PaymentMethodManagementComponent = asBaseComponent<GetPaymentMethodsForUserResponse, PaymentMethodManagementComponentOwnProps>(
    (props: GetPaymentMethodsForUserResponse & PaymentMethodManagementComponentOwnProps & BaseComponentProps) => {
        const [ selectedPaymentMethodID, setSelectedPaymentMethodID ] = useState<string>(null);
        const [ paymentMethods, setPaymentMethods ] = useState<Array<PaymentMethod>>(props.paymentMethods);

        const [ wasSuccessful, setWasSuccessful ] = useState<boolean>(false);

        const handleMarkPaymentMethodAsDefault = () => {
            props.setIsLoading(true);
            markPaymentMethodAsDefault({
                paymentMethodId: selectedPaymentMethodID,
            },
            (resp: MarkPaymentMethodAsDefaultResponse) => {
                props.setIsLoading(false);
                setWasSuccessful(true);
            },
            (err: Error) => {
                props.setIsLoading(false);
                props.setError(err);
            });
        }

        const handleDeletePaymentMethod = () => {
            props.setIsLoading(true);
            deletePaymentMethodForUser({
                paymentMethodId: selectedPaymentMethodID,
            },
            (resp: DeletePaymentMethodForUserResponse) => {
                props.setIsLoading(false);
                setWasSuccessful(true);
                setPaymentMethods(paymentMethods.filter((p: PaymentMethod) => p.externalId !== selectedPaymentMethodID));
            },
            (err: Error) => {
                props.setIsLoading(false);
                props.setError(err);
            });
        }

        const classes = styleClasses();
        return (
            <div>
                <Heading3 color={TypographyColor.Primary}>
                    Existing Payment Methods
                </Heading3>
                {
                    (!paymentMethods || !paymentMethods.length) ? (
                        <Paragraph>
                            You currently have no payment methods
                        </Paragraph>
                    ) : (
                        <div>
                            <Grid container>
                            {
                                paymentMethods.map((paymentMethod: PaymentMethod) => (
                                    <Grid item xs={12} md={6}
                                        className={classes.paymentMethodDisplayContainer}
                                        key={paymentMethod.externalId}>
                                        <PaymentMethodDisplay onClick={setSelectedPaymentMethodID} paymentMethod={paymentMethod} isHighlighted={paymentMethod.externalId === selectedPaymentMethodID} />
                                    </Grid>
                                ))
                            }
                            </Grid>
                            <CenteredComponent>
                                <PrimaryButton
                                    onClick={handleMarkPaymentMethodAsDefault}
                                    className={classes.paymentMethodButton}
                                    disabled={!selectedPaymentMethodID}>
                                    Make Default Payment Method
                                </PrimaryButton>
                            </CenteredComponent>
                            <CenteredComponent>
                                <WarningButton
                                    onClick={handleDeletePaymentMethod}
                                    className={classes.paymentMethodButton}
                                    disabled={!selectedPaymentMethodID}>
                                    Delete Payment Method
                                </WarningButton>
                            </CenteredComponent>
                        </div>
                    )
                }
                <Snackbar open={wasSuccessful} autoHideDuration={6000} onClose={() => setWasSuccessful(false)}>
                    <Alert severity="success">Your request was successful! Please allow a few minutes for our systems to update.</Alert>
                </Snackbar>
            </div>
        );
    },
    (
        ownProps: PaymentMethodManagementComponentOwnProps,
        onSuccess: (resp: GetPaymentMethodsForUserResponse) => void,
        onError: (err: Error) => void,
    ) => getPaymentMethodsForUser({}, onSuccess, onError),
    false
);

export default PremiumNewsletterSubscriptionManagementPage;
