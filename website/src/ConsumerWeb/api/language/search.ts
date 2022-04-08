import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { Lemma } from 'ConsumerWeb/api/model/language';
import { ClientError } from 'ConsumerWeb/api/clienterror';
import { WordsmithLanguageCode } from 'common/model/language/language';

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

export type PartOfSpeech = {
	id: string;
	name: string;
}

export type SearchTextRequest = {
    languageCode: WordsmithLanguageCode;
    captchaToken: string;
    wordReinforcementToken;
    text: string[];
}

export type SearchTextResponse = {
    error: ClientError | undefined;
    result: SearchTextResult | undefined;
}

export type SearchTextResult = {
    results: SearchResult[];
    languageCode: WordsmithLanguageCode;
}

export type SearchResult = {
    displayText: string;
    definitions: string[];
    partOfSpeech: PartOfSpeech | undefined;
    lookupId: LanguageLookupID;
}

export type LanguageLookupID = {
    idType: IDType;
    id: string[];
}

export enum IDType {
    Lemma = 'lemma',
    Phrase = 'phrase',
}

export function searchText(
    req: SearchTextRequest,
    onSuccess: (resp: SearchTextResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<SearchTextRequest, SearchTextResponse>(
        '/api/language/search_text_1',
        req,
        onSuccess,
        onError,
    );
}
