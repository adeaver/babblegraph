import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { setLocation } from 'util/window/Location';

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

type CreateUserAccountPageOwnProps = RouteComponentProps<Params>;

const CreateUserAccountPage = withUserProfileInformation<CreateUserAccountPageOwnProps>(
    RouteEncryptionKey.CreateUser,
    [RouteEncryptionKey.SubscriptionManagement, RouteEncryptionKey.PremiumSubscriptionCheckout],
    (ownProps: CreateUserAccountPageOwnProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.CheckoutPage,
    (props: CreateUserAccountPageOwnProps & UserProfileComponentProps) => {
        const [ subscriptionManagementToken, premiumSubscriptionCheckoutToken ] = props.userProfile.nextTokens;
        if (!!props.userProfile.subscriptionLevel) {
            setLocation(`/manage/${subscriptionManagementToken}`);
            return <div />;
        } else if (props.userProfile.hasAccount) {
            setLocation(`/checkout/${premiumSubscriptionCheckoutToken}`);
            return <div />;
        }
        return (
            <div />
        );
    }
);

export default CreateUserAccountPage;
