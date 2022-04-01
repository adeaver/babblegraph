import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { ClientError } from 'ConsumerWeb/api/clienterror';
import { WordsmithLanguageCode } from 'common/model/language/language';

export type Schedule = {
    ianaTimezone: string;
    hourIndex: number;
    quarterHourIndex: number;
    isActiveForDays: Array<boolean>;
}

export type UserNewsletterPreferences = {
    languageCode: WordsmithLanguageCode;
    isLemmaReinforcementSpotlightActive: boolean;
    arePodcastsEnabled: boolean;
    includeExplicitPodcasts: boolean;
    minimumPodcastDurationSeconds: number | undefined;
    maximumPodcastDurationSeconds: number | undefined;
    numberOfArticlesPerEmail: number;
    schedule: Schedule;
}

export type GetUserNewsletterPreferencesRequest = {
    languageCode: WordsmithLanguageCode;
    subscriptionManagementToken: string;
}

export type GetUserNewsletterPreferencesResponse = {
    preferences: UserNewsletterPreferences | undefined;
    error: ClientError | undefined;
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
    emailAddress: string | undefined;
    subscriptionManagementToken: string;
    preferences: UserNewsletterPreferences;
}

export type UpdateUserNewsletterPreferencesResponse = {
    success: boolean;
    error: ClientError | undefined;
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
