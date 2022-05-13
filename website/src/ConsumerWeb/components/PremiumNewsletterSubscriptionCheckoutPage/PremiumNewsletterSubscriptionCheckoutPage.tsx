import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { setLocation } from 'util/window/Location';
import { PrimaryButton } from 'common/components/Button/Button';
import Link, { LinkTarget } from 'common/components/Link/Link';

import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import {
    withUserProfileInformation,
    UserProfileComponentProps
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';
import PremiumNewsletterSubscriptionCardForm from 'ConsumerWeb/components/common/Billing/PremiumNewsletterSubscriptionCheckoutForm';

import {
    PaymentState,
    PremiumNewsletterSubscription
} from 'common/api/billing/billing';
import {
    GetOrCreateBillingInformationResponse,
    getOrCreateBillingInformation,

    GetOrCreatePremiumNewsletterSubscriptionResponse,
    getOrCreatePremiumNewsletterSubscription,
} from 'ConsumerWeb/api/billing/billing';

import { asRoundedFixedDecimal } from 'util/string/NumberString';

const styleClasses = makeStyles({
    callToActionButton: {
        margin: '15px 0',
        width: '100%',
    },
});

type Params = {
    token: string;
}

type PremiumNewsletterSubscriptionCheckoutPageProps = RouteComponentProps<Params>;

const PremiumNewsletterSubscriptionCheckoutPage = withUserProfileInformation<PremiumNewsletterSubscriptionCheckoutPageProps>(
    RouteEncryptionKey.PremiumSubscriptionCheckout,
    [RouteEncryptionKey.SubscriptionManagement, RouteEncryptionKey.CreateUser],
    (ownProps: PremiumNewsletterSubscriptionCheckoutPageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.CheckoutPage,
    (props: PremiumNewsletterSubscriptionCheckoutPageProps & UserProfileComponentProps) => {
        const { token } = props.match.params;
        const [ subscriptionManagementToken, createUserToken ] = props.userProfile.nextTokens;

        if (!props.userProfile.hasAccount) {
            setLocation(`/signup/${createUserToken}`);
            return;
        }

        return (
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <DisplayCardHeader
                            title="Subscription Checkout"
                            backArrowDestination={`/manage/${subscriptionManagementToken}`} />
                        <OrderDetailsSection
                            trialEligibilityDays={props.userProfile.trialEligibilityDays}
                            premiumSubscriptionCheckoutToken={token}
                            subscriptionManagementToken={subscriptionManagementToken} />
                    </DisplayCard>
                </Grid>
            </Grid>
        );
    }
);

type OrderDetailsSectionProps = {
    premiumSubscriptionCheckoutToken: string;
    trialEligibilityDays: number | undefined;
    subscriptionManagementToken: string;
}

type OrderDetailsSectionAPIProps = GetOrCreatePremiumNewsletterSubscriptionResponse & GetOrCreateBillingInformationResponse;

const OrderDetailsSection = asBaseComponent<OrderDetailsSectionAPIProps, OrderDetailsSectionProps>(
    (props: OrderDetailsSectionAPIProps & OrderDetailsSectionProps & BaseComponentProps) => {
        const [ shouldShowCheckoutForm, setShouldShowCheckoutForm ] = useState<boolean>(false);

        const classes = styleClasses();
        return (
            <Grid container>
                <Grid item xs={12}>
                    <Heading3 align={Alignment.Left}>
                        Your Order
                    </Heading3>
                </Grid>
                <Grid item xs={12}>
                    <Divider />
                </Grid>
                <Grid item xs={8}>
                    <Paragraph align={Alignment.Left}>
                        1-year Babblegraph Subscription
                    </Paragraph>
                </Grid>
                <Grid item xs={4}>
                    <Paragraph align={Alignment.Right}>
                        US${asRoundedFixedDecimal(props.premiumNewsletterSubscription.priceCents / 100.0, 2)}
                    </Paragraph>
                </Grid>
                <Grid item xs={12}>
                    <Divider />
                </Grid>
                <Grid item xs={8}>
                    <Paragraph align={Alignment.Left}>
                        Total Due Now
                    </Paragraph>
                </Grid>
                <Grid item xs={4}>
                    <Paragraph align={Alignment.Right}>
                        { !!props.trialEligibilityDays ? "US$0.00" : `$${asRoundedFixedDecimal(props.premiumNewsletterSubscription.priceCents / 100.0, 2)}`}
                    </Paragraph>
                </Grid>
                {
                    !!props.trialEligibilityDays && (
                        <Grid item xs={8}>
                            <Paragraph align={Alignment.Left}>
                                Total Due In {props.trialEligibilityDays} Days
                            </Paragraph>
                        </Grid>
                    )
                }
                {
                    !!props.trialEligibilityDays && (
                        <Grid item xs={4}>
                            <Paragraph align={Alignment.Right}>
                                US${asRoundedFixedDecimal(props.premiumNewsletterSubscription.priceCents / 100.0, 2)}
                            </Paragraph>
                        </Grid>
                    )
                }
                {
                    props.premiumNewsletterSubscription.hasValidDiscount && (
                        <Grid item xs={12}>
                            <Paragraph size={Size.Small} align={Alignment.Center}>
                                Renews at US$29.00 after the first year
                            </Paragraph>
                        </Grid>
                    )
                }
                <Grid item xs={12}>
                    <PrimaryButton
                        onClick={() => {setShouldShowCheckoutForm(true)}}
                        className={classes.callToActionButton}
                        disabled={shouldShowCheckoutForm}
                        size="large">
                        {
                            !props.trialEligibilityDays ? "Proceed to pay" : "Add a payment method"
                        }
                    </PrimaryButton>
                </Grid>
                {
                    shouldShowCheckoutForm && (
                        <Grid item xs={12}>
                            <PaymentSection
                                premiumSubscriptionCheckoutToken={props.premiumSubscriptionCheckoutToken}
                                subscriptionManagementToken={props.subscriptionManagementToken}
                                premiumNewsletterSubscription={props.premiumNewsletterSubscription} />
                        </Grid>
                    )
                }
            </Grid>
        );
    },
    (
        ownProps: OrderDetailsSectionProps,
        onSuccess: (OrderDetailsSectionAPIProps) => void,
        onError: (err: Error) => void,
    ) => {
        getOrCreateBillingInformation({
            premiumSubscriptionCheckoutToken: ownProps.premiumSubscriptionCheckoutToken,
        },
        (resp: GetOrCreateBillingInformationResponse) => {
            getOrCreatePremiumNewsletterSubscription({
                premiumSubscriptionCheckoutToken: ownProps.premiumSubscriptionCheckoutToken,
            },
            (resp2: GetOrCreatePremiumNewsletterSubscriptionResponse) => {
                onSuccess({
                    ...resp,
                    ...resp2,
                });
            },
            onError);
        },
        onError);
    },
    false,
);

type PaymentSectionProps = {
    premiumSubscriptionCheckoutToken: string;
    subscriptionManagementToken: string;
    premiumNewsletterSubscription: PremiumNewsletterSubscription;
}

const PaymentSection = (props: PaymentSectionProps) => {
    switch (props.premiumNewsletterSubscription.paymentState) {
        case PaymentState.CreatedUnpaid:
        case PaymentState.TrialNoPaymentMethod:
            return (
                <PremiumNewsletterSubscriptionCardForm
                    premiumNewsletterSusbcription={props.premiumNewsletterSubscription}
                    subscriptionManagementToken={props.subscriptionManagementToken} />
            );
        case PaymentState.TrialPaymentMethodAdded:
        case PaymentState.Active:
        case PaymentState.Errored:
            return (
                <div>
                    <Heading3 color={TypographyColor.Warning}>
                        It looks like we’ve already collected your payment information!
                    </Heading3>
                    <Link href={`/manage/${props.subscriptionManagementToken}/payment-settings`} target={LinkTarget.Self}>
                        You can make changes to your payment information by clicking here.
                    </Link>
                </div>
            );
    }
    return (
        <div>
            <Heading3 color={TypographyColor.Warning}>
                Something went wrong. Try again later or contact hello@babblegraph.com for help.
            </Heading3>
            <Link href={`/manage/${props.subscriptionManagementToken}`} target={LinkTarget.Self}>
                Go back to subscription management
            </Link>
        </div>
    )
}

export default PremiumNewsletterSubscriptionCheckoutPage;
