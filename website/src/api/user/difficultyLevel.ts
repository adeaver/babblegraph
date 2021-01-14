import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type GetUserPreferencesForTokenRequest = {
    token: string;
}

export type GetUserPreferencesForTokenResponse = {
    classificationsByLanguage: Array<ReadingLevelClassificationForLanguage>;
}

export type ReadingLevelClassificationForLanguage = {
    languageCode: string;
    readingLevelClassification: string;
}

export function getUserPreferencesForToken(
    req: GetUserPreferencesForTokenRequest,
    onSuccess: (resp: GetUserPreferencesForTokenResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserPreferencesForTokenRequest, GetUserPreferencesForTokenResponse>(
        '/api/user/get_user_preferences_for_token_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateUserPreferencesForTokenRequest = {
    token: string;
    emailAddress: string;
    classificationsByLanguage: Array<ReadingLevelClassificationForLanguage>;
}

export type UpdateUserPreferencesForTokenResponse = {
    didUpdate: boolean;
}

export function updateUserPreferencesForToken(
    req: UpdateUserPreferencesForTokenRequest,
    onSuccess: (resp: UpdateUserPreferencesForTokenResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateUserPreferencesForTokenRequest, UpdateUserPreferencesForTokenResponse>(
        '/api/user/update_user_preferences_for_token_1',
        req,
        onSuccess,
        onError,
    );
}
