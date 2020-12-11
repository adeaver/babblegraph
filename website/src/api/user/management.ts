import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type GetUserPreferencesForTokenRequest = {
    Token: string;
}

export type GetUserPreferencesForTokenResponse = {
    ClassificationsByLanguage: Array<ReadingLevelClassificationForLanguage>;
}

export type ReadingLevelClassificationForLanguage = {
    LanguageCode: string;
    ReadingLevelClassification: string;
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
