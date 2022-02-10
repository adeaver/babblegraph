import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

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
    // TODO: type this
    key: string;
    toKey: string | undefined;
}

export type GetUserProfileInformationResponse = {
    error: UserProfileInformationError | undefined;
    userProfile: UserProfileInformation | undefined;
}

export type UserProfileInformation = {
	hasAccount: boolean;
    isLoggedIn: boolean;
	subscriptionLevel: SubscriptionLevel | undefined;
	nextToken: string | undefined;
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
