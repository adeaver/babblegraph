import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';

import { WordsmithLanguageCode } from 'common/model/language/language';
import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import {
    withUserProfileInformation,
    UserProfileComponentProps,
} from 'ConsumerWeb/base/UserProfile/withUserProfile';

import InterestSelector from './InterestSelector';

type Params = {
    token: string;
}

type InterestSelectionPageProps = RouteComponentProps<Params>;

const InterestSelectionPage = withUserProfileInformation<InterestSelectionPageProps>(
    RouteEncryptionKey.SubscriptionManagement,
    [ RouteEncryptionKey.SubscriptionManagement ],
    (ownProps: InterestSelectionPageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.ContentTopics,
    (props: InterestSelectionPageProps & UserProfileComponentProps) => {
        const { token } = props.match.params;

        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Select interests"
                        backArrowDestination={`/manage/${token}`} />
                    <InterestSelector
                        subscriptionManagementToken={token}
                        languageCode={WordsmithLanguageCode.Spanish}
                        omitEmailAddress={props.userProfile.hasAccount} />
                </DisplayCard>
            </CenteredComponent>
        );
    }
);

export default InterestSelectionPage;
