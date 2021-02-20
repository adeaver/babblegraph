import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type AddUserLemmasForTokenRequest = {
    token: string;
    lemmaId: string;
}

export type AddUserLemmasForTokenResponse = {
    didUpdate: boolean;
}

export function addUserLemmasForToken(
    req: AddUserLemmasForTokenRequest,
    onSuccess: (resp: AddUserLemmasForTokenResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<AddUserLemmasForTokenRequest, AddUserLemmasForTokenResponse>(
        '/api/user/add_user_lemma_for_token_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetUserLemmasForTokenRequest = {
    token: string;
}

export type GetUserLemmasForTokenResponse = {
    lemmaMappingsByLanguageCode: Array<LemmaMappingsWithLanguageCode>;
}

export type LemmaMappingsWithLanguageCode = {
    languageCode: string;
    lemmaMappings: Array<LemmaMapping>;
}

export type LemmaMapping = {
    ID: string;
    LanguageCode: string;
    UserID: string;
    LemmaID: string;
    IsVisible: string;
    IsActive: string;
}

export function getUserLemmasForToken(
    req: GetUserLemmasForTokenRequest,
    onSuccess: (resp: GetUserLemmasForTokenResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserLemmasForTokenRequest, GetUserLemmasForTokenResponse>(
        '/api/user/get_user_lemmas_for_token_1',
        req,
        onSuccess,
        onError,
    );
}
