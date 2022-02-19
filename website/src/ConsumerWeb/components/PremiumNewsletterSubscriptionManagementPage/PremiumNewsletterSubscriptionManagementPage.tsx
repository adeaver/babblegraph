import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Link, { LinkTarget } from 'common/components/Link/Link';

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

import {
    PremiumNewsletterSubscription,
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
        return (
            <Heading3 color={TypographyColor.Primary}>
                Cool, cool
            </Heading3>
        );
    },
    (
        ownProps: PremiumSubscriptionManagementComponentOwnProps,
        onSuccess: (resp: LookupActivePremiumNewsletterSubscriptionResponse) => void,
        onError: (err: Error) => void,
    ) => lookupActivePremiumNewsletterSubscription({subscriptionManagementToken: ownProps.subscriptionManagementToken}, onSuccess, onError),
    false
);

export default PremiumNewsletterSubscriptionManagementPage;
