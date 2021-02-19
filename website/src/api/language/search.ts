import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

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

export type Lemma = {
    text: string;
    id: string;
    partOfSpeech: PartOfSpeech;
    definitions: Definition[];
}

export type PartOfSpeech = {
    id: string;
    name: string;
}

export type Definition = {
    text: string;
    extraInfo: string | undefined;
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
