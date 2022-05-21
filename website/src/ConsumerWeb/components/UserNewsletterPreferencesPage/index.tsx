import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import {
    withUserProfileInformation,
    UserProfileComponentProps,
} from 'ConsumerWeb/base/UserProfile/withUserProfile';

import {
    DateTimeSelector
} from './components';

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
    LoginRedirectKey.NewsletterPreferences,
    (props: UserNewsletterPreferencesPageProps & UserProfileComponentProps) => {
        return (
            <DateTimeSelector />
        );
    }
);

export default UserNewsletterPreferencesPage;
