import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';

import {
    withUserProfileInformation,
    UserProfileComponentProps
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';

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
                </DisplayCard>
            </CenteredComponent>
        )
    }
);

export default PremiumNewsletterSubscriptionManagementPage;
