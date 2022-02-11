import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';

export enum UserProfileInformationError {
    InvalidKey = 'invalid-key',
    InvalidToken = 'invalid-token',
}

export enum SubscriptionLevel {
    Premium = 'premium',
    BetaPremium = 'beta-premium',
}

export type GetUserProfileInformationRequest = {
    token: string;
    key: RouteEncryptionKey;
    nextKeys: Array<RouteEncryptionKey>;
}

export type GetUserProfileInformationResponse = {
    error: UserProfileInformationError | undefined;
    userProfile: UserProfileInformation | undefined;
}

export type UserProfileInformation = {
	hasAccount: boolean;
    isLoggedIn: boolean;
	subscriptionLevel: SubscriptionLevel | undefined;
	nextTokens: Array<string>;
}

export function getUserProfileInformation(
    req: GetUserProfileInformationRequest,
    onSuccess: (resp: GetUserProfileInformationResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserProfileInformationRequest, GetUserProfileInformationResponse>(
        '/api/useraccounts/get_user_profile_information_1',
        req,
        onSuccess,
        onError,
    );
}
