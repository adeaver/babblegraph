import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { Lemma } from 'ConsumerWeb/api/model/language';

export type GetLemmasMatchingTextRequest = {
    languageCode: string;
    token: string;
    text: string;
}

export type GetLemmasMatchingTextResponse = {
    languageCode: string;
    text: string;
    lemmas: Lemma[];
}

export function getLemmasMatchingText(
    req: GetLemmasMatchingTextRequest,
    onSuccess: (resp: GetLemmasMatchingTextResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetLemmasMatchingTextRequest, GetLemmasMatchingTextResponse>(
        '/api/language/get_lemmas_matching_text_1',
        req,
        onSuccess,
        onError,
    );
}
