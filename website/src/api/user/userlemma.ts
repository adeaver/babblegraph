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
