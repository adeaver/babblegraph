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

import TimeSelector from './TimeSelector';

type Params = {
    token: string;
}

type UserNewsletterPreferencesPageProps = RouteComponentProps<Params>;

const UserNewsletterPreferencesPage = withUserProfileInformation<UserNewsletterPreferencesPageProps>(
    RouteEncryptionKey.SubscriptionManagement,
    [],
    (ownProps: UserNewsletterPreferencesPageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.NewsletterSchedule,
    (props: UserNewsletterPreferencesPageProps & UserProfileComponentProps) => {
        const { token } = props.match.params;

        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Schedule settings"
                        backArrowDestination={`/manage/${token}`} />
                    <TimeSelector
                        subscriptionManagementToken={token}
                        languageCode={WordsmithLanguageCode.Spanish} />
                </DisplayCard>
            </CenteredComponent>
        )
    }
);

export default UserNewsletterPreferencesPage;
