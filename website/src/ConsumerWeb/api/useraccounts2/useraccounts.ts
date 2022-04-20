import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import { SubscriptionLevel } from 'common/api/useraccounts/useraccounts';

export enum UserProfileInformationError {
    InvalidKey = 'invalid-key',
    InvalidToken = 'invalid-token',
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
    hasPaymentMethod: boolean;
	subscriptionLevel: SubscriptionLevel | undefined;
    trialEligibilityDays: number | undefined;
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

export type CreateUserRequest = {
    createUserToken: string;
    emailAddress: string;
    password: string;
    confirmPassword: string;
}

export type CreateUserResponse = {
    createUserError: CreateUserError | null;
}

export enum CreateUserError {
    AlreadyExists = 'already-exists',
    InvalidToken = 'invalid-token',
    PasswordRequirements = 'pass-requirements',
    PasswordsNoMatch = 'passwords-no-match',
}

export function createUser(
    req: CreateUserRequest,
    onSuccess: (resp: CreateUserResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<CreateUserRequest, CreateUserResponse>(
        '/api/useraccounts/create_user_1',
        req,
        onSuccess,
        onError,
    );
}

