import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import {
    withUserProfileInformation,
    UserProfileComponentProps
} from 'ConsumerWeb/base/UserProfile/withUserProfile';

type Params = {
    token: string;
}

type PremiumInformationPageOwnProps = RouteComponentProps<Params>

const PremiumInformationPage = withUserProfileInformation<PremiumInformationPageOwnProps>(
    "subscription-management",
    undefined,
    (ownProps: PremiumInformationPageOwnProps) => {
        return ownProps.match.params.token;
    },
    "cts",
    (props: PremiumInformationPageOwnProps & UserProfileComponentProps) => {
        return (
            <div>
                { `${props.userProfile.isLoggedIn}` }
            </div>
        )
    }
);

export default PremiumInformationPage;
