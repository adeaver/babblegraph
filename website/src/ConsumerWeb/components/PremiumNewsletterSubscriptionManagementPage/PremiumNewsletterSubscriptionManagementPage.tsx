import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Link, { LinkTarget } from 'common/components/Link/Link';
import Paragraph, { Size } from 'common/typography/Paragraph';

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
    PremiumNewsletterSubscription,
    PaymentState,
    LookupActivePremiumNewsletterSubscriptionResponse,
    lookupActivePremiumNewsletterSubscription,
} from 'ConsumerWeb/api/billing/billing';

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

export default PremiumNewsletterSubscriptionManagementPage;
