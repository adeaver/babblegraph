import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type GetUserNewsletterPreferencesRequest = {
    languageCode: string;
    emailAddress: string;
    subscriptionManagementToken: string;
}

export type GetUserNewsletterPreferencesResponse = {
    languageCode: string;
    preferences: UserNewsletterPreferences;
}

export type UserNewsletterPreferences = {
    isLemmaReinforcementSpotlightActive: boolean;
}

export function getUserNewsletterPreferences(
    req: GetUserNewsletterPreferencesRequest,
    onSuccess: (resp: GetUserNewsletterPreferencesResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserNewsletterPreferencesRequest, GetUserNewsletterPreferencesResponse>(
        '/api/user/get_user_newsletter_preferences_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateUserNewsletterPreferencesRequest = {
    languageCode: string;
    emailAddress: string;
    subscriptionManagementToken: string;
    preferences: UserNewsletterPreferences;
}

export type UpdateUserNewsletterPreferencesResponse = {
    languageCode: string;
    success: boolean;
}

export function updateUserNewsletterPreferences(
    req: UpdateUserNewsletterPreferencesRequest,
    onSuccess: (resp: UpdateUserNewsletterPreferencesResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateUserNewsletterPreferencesRequest, UpdateUserNewsletterPreferencesResponse>(
        '/api/user/update_user_newsletter_preferences_1',
        req,
        onSuccess,
        onError,
    );
}
