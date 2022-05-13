import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Link, { LinkTarget } from 'common/components/Link/Link';
import Paragraph, { Size } from 'common/typography/Paragraph';
import PaymentMethodDisplay from 'ConsumerWeb/components/common/Billing/PaymentMethodDisplay';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { PrimaryButton, WarningButton } from 'common/components/Button/Button';
import { SubscriptionLevel } from 'common/api/useraccounts/useraccounts';

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
    PaymentState,
    PremiumNewsletterSubscription
} from 'common/api/billing/billing';
import {
    PaymentMethod,

    LookupActivePremiumNewsletterSubscriptionResponse,
    lookupActivePremiumNewsletterSubscription,

    SetPremiumNewsletterSubscriptionAutoRenewResponse,
    setPremiumNewsletterSubscriptionAutoRenew,

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
    autoRenewButton: {
        width: '100%',
        margin: '10px 0',
    }
});

type Params = {
    token: string;
}

type PremiumNewsletterSubscriptionManagementPageProps = RouteComponentProps<Params>;

const PremiumNewsletterSubscriptionManagementPage = withUserProfileInformation<PremiumNewsletterSubscriptionManagementPageProps>(
    RouteEncryptionKey.SubscriptionManagement,
    [RouteEncryptionKey.PremiumSubscriptionCheckout],
    (ownProps: PremiumNewsletterSubscriptionManagementPageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.PaymentSettings,
    (props: PremiumNewsletterSubscriptionManagementPageProps & UserProfileComponentProps) => {
        const { token } = props.match.params;
        const [ checkoutToken ] = props.userProfile.nextTokens;

        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Payment and Subscription Settings"
                        backArrowDestination={`/manage/${token}`} />
                    <PremiumSubscriptionManagementComponent
                        subscriptionManagementToken={token}
                        checkoutToken={checkoutToken}
                        subscriptionLevel={props.userProfile.subscriptionLevel} />
                    <PaymentMethodManagementComponent />
                    <Paragraph color={TypographyColor.Primary}>
                        Need help? Just reach out to hello@babblegraph.com
                    </Paragraph>
                </DisplayCard>
            </CenteredComponent>
        )
    }
);

type PremiumSubscriptionManagementComponentOwnProps = {
    subscriptionManagementToken: string;
    checkoutToken: string;

    subscriptionLevel: SubscriptionLevel | null;
}

const PremiumSubscriptionManagementComponent = asBaseComponent<LookupActivePremiumNewsletterSubscriptionResponse, PremiumSubscriptionManagementComponentOwnProps>(
    (props: LookupActivePremiumNewsletterSubscriptionResponse & PremiumSubscriptionManagementComponentOwnProps & BaseComponentProps) => {
        if (!props.premiumNewsletterSubscription) {
            return (
                <div>
                    <Heading3 color={TypographyColor.Warning}>
                        You don’t have an active Babblegraph subscription
                    </Heading3>
                    <Link
                        href={props.subscriptionLevel === SubscriptionLevel.Legacy ? (
                                `/manage/${props.subscriptionManagementToken}/premium`
                            ) : (
                                `/checkout/${props.checkoutToken}`
                            )
                        }
                        target={LinkTarget.Self}>
                        If you’d like to start or restart one, click here.
                    </Link>
                </div>
            );
        }

        const [ isLoadingAutoRenew, setIsLoadingAutoRenew ] = useState<boolean>(false);
        const [ isAutoRenewEnabled, setIsAutoRenewEnabled ] = useState<boolean>(props.premiumNewsletterSubscription.isAutoRenewEnabled);
        const handleUpdateAutoRenew = (isAutoRenewEnabled: boolean) => {
            return () => {
                setIsLoadingAutoRenew(true);
                setPremiumNewsletterSubscriptionAutoRenew({
                    subscriptionManagementToken: props.subscriptionManagementToken,
                    isAutoRenewEnabled: isAutoRenewEnabled,
                },
                (resp: SetPremiumNewsletterSubscriptionAutoRenewResponse) => {
                    setIsLoadingAutoRenew(false);
                    setIsAutoRenewEnabled(isAutoRenewEnabled);
                },
                (err: Error) => {
                    setIsLoadingAutoRenew(false);
                    props.setError(err);
                });
            }
        }
        const { paymentState } = props.premiumNewsletterSubscription;

        const classes = styleClasses();
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
                            You’re currently trialing Babblegraph until {new Date(props.premiumNewsletterSubscription.currentPeriodEnd).toLocaleDateString()}
                        </Heading3>
                        <Paragraph>
                            However, you haven’t added a payment method, so you are set to lose access to Babblegraph on that date. You can add a payment method with the form below.
                        </Paragraph>
                        <PremiumNewsletterSubscriptionCardForm
                            premiumNewsletterSusbcription={props.premiumNewsletterSubscription} />
                    </div>
                );
            case PaymentState.TrialPaymentMethodAdded:
                return (
                    <div>
                        <Heading3 color={TypographyColor.Primary}>
                            You’re currently trialing Babblegraph until {new Date(props.premiumNewsletterSubscription.currentPeriodEnd).toLocaleDateString()}
                        </Heading3>
                        {
                            isAutoRenewEnabled ? (
                                <div>
                                    <Paragraph>
                                        You will be charged on that date. If you do not want to be charged, you can disable your auto-renew below
                                    </Paragraph>
                                    {
                                        isLoadingAutoRenew ? (
                                            <LoadingSpinner />
                                        ) : (
                                            <CenteredComponent>
                                                <WarningButton className={classes.autoRenewButton} onClick={handleUpdateAutoRenew(false)}>
                                                    Disable auto-renew
                                                </WarningButton>
                                            </CenteredComponent>
                                        )
                                    }
                                </div>
                            ) : (
                                <div>
                                    <Paragraph color={TypographyColor.Warning}>
                                        Auto-renew is currently disabled, so you’ll lose access to Babblegraph at that date. You can enable it below.
                                    </Paragraph>
                                    {
                                        isLoadingAutoRenew ? (
                                            <LoadingSpinner />
                                        ) : (
                                            <CenteredComponent>
                                                <PrimaryButton className={classes.autoRenewButton} onClick={handleUpdateAutoRenew(true)}>
                                                    Enable auto-renew
                                                </PrimaryButton>
                                            </CenteredComponent>
                                        )
                                    }
                                </div>
                            )
                        }
                        <Paragraph>
                             You can also add a new payment method below.
                        </Paragraph>
                        <PremiumNewsletterSubscriptionCardForm
                            premiumNewsletterSusbcription={props.premiumNewsletterSubscription} />
                    </div>
                );
            case PaymentState.Active:
                return (
                    <div>
                        <Heading3 color={TypographyColor.Primary}>
                            You’re currently subscribed to Babblegraph
                        </Heading3>
                        {
                            isAutoRenewEnabled ? (
                                <div>
                                    <Paragraph>
                                        You will be charged next on {new Date(props.premiumNewsletterSubscription.currentPeriodEnd).toLocaleDateString()}. If you do not want to be charged, you can disable your auto-renew below
                                    </Paragraph>
                                    {
                                        isLoadingAutoRenew ? (
                                            <LoadingSpinner />
                                        ) : (
                                            <CenteredComponent>
                                                <WarningButton className={classes.autoRenewButton} onClick={handleUpdateAutoRenew(false)}>
                                                    Disable auto-renew
                                                </WarningButton>
                                            </CenteredComponent>
                                        )
                                    }
                                </div>
                            ) : (
                                <div>
                                    <Paragraph color={TypographyColor.Warning}>
                                        Auto-renew is currently disabled, so you’ll lose access to Babblegraph on {new Date(props.premiumNewsletterSubscription.currentPeriodEnd).toLocaleDateString()}. You can enable it below.
                                    </Paragraph>
                                    {
                                        isLoadingAutoRenew ? (
                                            <LoadingSpinner />
                                        ) : (
                                            <CenteredComponent>
                                                <PrimaryButton className={classes.autoRenewButton} onClick={handleUpdateAutoRenew(true)}>
                                                    Enable auto-renew
                                                </PrimaryButton>
                                            </CenteredComponent>
                                        )
                                    }
                                </div>
                            )
                        }
                        <Paragraph>
                             You can also add a new payment method below.
                        </Paragraph>
                        <PremiumNewsletterSubscriptionCardForm
                            premiumNewsletterSusbcription={props.premiumNewsletterSubscription} />
                    </div>
                );
            case PaymentState.Errored:
                return (
                    <div>
                        <Heading3 color={TypographyColor.Primary}>
                            We’ve encountered an error processing payment for your Babblegraph subscription
                        </Heading3>
                        <Paragraph>
                            If you’d like to continue using Babblegraph, please designate a new default payment method. If you would not like to continue, then no action is required.
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

        const [ isLoading, setIsLoading ] = useState<boolean>(false);
        const [ wasSuccessful, setWasSuccessful ] = useState<boolean>(false);

        const handleMarkPaymentMethodAsDefault = () => {
            setIsLoading(true);
            markPaymentMethodAsDefault({
                paymentMethodId: selectedPaymentMethodID,
            },
            (resp: MarkPaymentMethodAsDefaultResponse) => {
                setIsLoading(false);
                setWasSuccessful(true);
                setPaymentMethods(paymentMethods.map((p: PaymentMethod) => ({
                    ...p,
                    isDefault: p.externalId === selectedPaymentMethodID
                })));
            },
            (err: Error) => {
                setIsLoading(false);
                props.setError(err);
            });
        }

        const handleDeletePaymentMethod = () => {
            setIsLoading(true);
            deletePaymentMethodForUser({
                paymentMethodId: selectedPaymentMethodID,
            },
            (resp: DeletePaymentMethodForUserResponse) => {
                setIsLoading(false);
                setWasSuccessful(true);
                setPaymentMethods(paymentMethods.filter((p: PaymentMethod) => p.externalId !== selectedPaymentMethodID));
            },
            (err: Error) => {
                setIsLoading(false);
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
                                    disabled={!selectedPaymentMethodID || isLoading}>
                                    Make Default Payment Method
                                </PrimaryButton>
                            </CenteredComponent>
                            <CenteredComponent>
                                <WarningButton
                                    onClick={handleDeletePaymentMethod}
                                    className={classes.paymentMethodButton}
                                    disabled={!selectedPaymentMethodID || isLoading}>
                                    Delete Payment Method
                                </WarningButton>
                            </CenteredComponent>
                            {
                                isLoading && <LoadingSpinner />
                            }
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
