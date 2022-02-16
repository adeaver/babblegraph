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

import {
    PaymentState,
    PremiumNewsletterSubscription,

    GetOrCreateBillingInformationResponse,
    getOrCreateBillingInformation,

    GetOrCreatePremiumNewsletterSubscriptionResponse,
    getOrCreatePremiumNewsletterSubscription,
} from 'ConsumerWeb/api/billing/billing';

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

        const [ shouldShowCheckoutForm, setShouldShowCheckoutForm ] = useState<boolean>(false);

        if (!props.userProfile.hasAccount) {
            setLocation(`/signup/${createUserToken}`);
            return;
        } else if (!!props.userProfile.subscriptionLevel) {
            setLocation(`/manage/${subscriptionManagementToken}`);
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
                            title="Babblegraph Premium Checkout"
                            backArrowDestination={`/manage/${subscriptionManagementToken}`} />
                        <OrderDetailsSection
                            trialEligibilityDays={props.userProfile.trialEligibilityDays}
                            premiumSubscriptionCheckoutToken={token}
                            isButtonDisabled={shouldShowCheckoutForm}
                            handleProceedToCheckout={() => setShouldShowCheckoutForm(true)} />
                        {
                            shouldShowCheckoutForm && (
                                <PaymentSection
                                    premiumSubscriptionCheckoutToken={token} />
                            )
                        }
                    </DisplayCard>
                </Grid>
            </Grid>
        );
    }
);

type OrderDetailsSectionProps = {
    premiumSubscriptionCheckoutToken: string;
    trialEligibilityDays: number | undefined;
    isButtonDisabled: boolean;

    handleProceedToCheckout: () => void;
}

const OrderDetailsSection = asBaseComponent<GetOrCreateBillingInformationResponse, OrderDetailsSectionProps>(
    (props: GetOrCreateBillingInformationResponse & OrderDetailsSectionProps & BaseComponentProps) => {
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
                        1-year Babblegraph Premium Subscription
                    </Paragraph>
                </Grid>
                <Grid item xs={4}>
                    <Paragraph align={Alignment.Right}>
                        US$29.00
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
                        { !!props.trialEligibilityDays ? "US$0.00" : "US$29.00"}
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
                                US$29.00
                            </Paragraph>
                        </Grid>
                    )
                }
                <Grid item xs={12}>
                    <PrimaryButton
                        onClick={props.handleProceedToCheckout}
                        className={classes.callToActionButton}
                        disabled={props.isButtonDisabled}
                        size="large">
                        {
                            !props.trialEligibilityDays ? "Proceed to checkout" : "Start your trial"
                        }
                    </PrimaryButton>
                </Grid>
            </Grid>
        );
    },
    (
        ownProps: OrderDetailsSectionProps,
        onSuccess: (GetOrCreateBillingInformationResponse) => void,
        onError: (err: Error) => void,
    ) => {
        getOrCreateBillingInformation({
            premiumSubscriptionCheckoutToken: ownProps.premiumSubscriptionCheckoutToken,
        },
        onSuccess,
        onError);
    },
    false,
);

type PaymentSectionProps = {
    premiumSubscriptionCheckoutToken: string;
}

const PaymentSection = asBaseComponent<GetOrCreatePremiumNewsletterSubscriptionResponse, PaymentSectionProps>(
    (props: GetOrCreatePremiumNewsletterSubscriptionResponse & PaymentSectionProps & BaseComponentProps) => {
        return <div />
    },
    (
        ownProps: PaymentSectionProps,
        onSuccess: (GetOrCreatePremiumNewsletterSubscriptionResponse) => void,
        onError: (err: Error) => void,
    ) => {
        getOrCreatePremiumNewsletterSubscription({
            premiumSubscriptionCheckoutToken: ownProps.premiumSubscriptionCheckoutToken,
        },
        onSuccess,
        onError);
    },
    false,
)

export default PremiumNewsletterSubscriptionCheckoutPage;
