import React from 'react';

import {
    asBaseComponent,
    BaseComponentProps,
    InitialDataFunc,
} from 'common/base/BaseComponent';
import { setLocation } from 'util/window/Location';

import {
    SubscriptionLevel
} from 'common/api/useraccounts/useraccounts';
import {
    UserProfileInformationError,
    UserProfileInformation,
    GetUserProfileInformationResponse,
    getUserProfileInformation,
} from 'ConsumerWeb/api/useraccounts2/useraccounts';
import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';

export type UserProfileComponentProps = {
    userProfile: UserProfileInformation;
} & BaseComponentProps;

export function withUserProfileInformation<P>(
    key: RouteEncryptionKey,
    nextKeys: Array<RouteEncryptionKey>,
    getToken: (ownProps: P) => string,
    loginRedirectKey: LoginRedirectKey | undefined,
    WrappedComponent: React.ComponentType<P & UserProfileComponentProps>,
) {
    return asBaseComponent<GetUserProfileInformationResponse, P>(
        (props: P & GetUserProfileInformationResponse & BaseComponentProps) => {
            if (!props.error && !props.userProfile) {
                return <div />;
            } else if (props.error) {
                return (
                    <p>Error</p>
                )
            }
            if (props.userProfile.hasAccount && !props.userProfile.isLoggedIn) {
                setLocation(`/login${!!loginRedirectKey ? "?d=" + loginRedirectKey : ""}`);
                return <div />;
            }
            // @ts-ignore
            return <WrappedComponent {...props as P & UserProfileComponentProps} />;
        },
        (
            ownProps: P,
            onSuccess: (apiProps: GetUserProfileInformationResponse) => void,
            onError: (err: Error) => void,
        ) => {
            getUserProfileInformation({
                key: key,
                token: getToken(ownProps),
                nextKeys: nextKeys,
            },
            onSuccess,
            onError)
        },
        true,
    )
}
