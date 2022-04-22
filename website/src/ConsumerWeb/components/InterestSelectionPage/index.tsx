import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';

import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import {
    withUserProfileInformation,
    UserProfileComponentProps,
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import { WordsmithLanguageCode } from 'common/model/language/language';

import {
    TopicSelector,
} from './components';

type Params = {
    token: string
}

type InterestSelectionPageProps = RouteComponentProps<Params>;

const InterestSelectionPage = withUserProfileInformation<InterestSelectionPageProps>(
    RouteEncryptionKey.SubscriptionManagement,
    [],
    (ownProps: InterestSelectionPageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.ContentTopics,
    (props: InterestSelectionPageProps & UserProfileComponentProps) => {
        const subscriptionManagementToken = props.match.params.token;
        return (
            <CenteredComponent useLargeVersion>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Manage your interests"
                        backArrowDestination={`/manage/${props.match.params.token}`} />
                    <TopicSelector languageCode={WordsmithLanguageCode.Spanish} />

                </DisplayCard>
            </CenteredComponent>
        );
    }
);

export default InterestSelectionPage;
