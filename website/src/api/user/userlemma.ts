import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';
import { Lemma } from 'api/model/language';

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
    lemmaMappings: Array<LemmaMapping>;
}

export type LemmaMapping = {
    isActive: boolean;
    lemma: Lemma;
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
