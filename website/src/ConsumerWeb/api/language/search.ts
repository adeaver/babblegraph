import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { ClientError } from 'ConsumerWeb/api/clienterror';
import { WordsmithLanguageCode } from 'common/model/language/language';

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
    uniqueHash: string;
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
