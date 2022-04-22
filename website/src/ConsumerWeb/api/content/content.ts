import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { ClientError } from 'ConsumerWeb/api/clienterror';
import { WordsmithLanguageCode } from 'common/model/language/language';

export type GetTopicsForLanguageRequest = {
    languageCode: WordsmithLanguageCode,
}

export type GetTopicsForLanguageResponse = {
    error: ClientError | undefined;
    results: Topic[];
}

export type Topic = {
    topicId: string;
    displayName: string;
    englishLabel: string;
}

export function getTopicsForLanguage(
    req: GetTopicsForLanguageRequest,
    onSuccess: (resp: GetTopicsForLanguageResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetTopicsForLanguageRequest, GetTopicsForLanguageResponse>(
        '/api/content/get_topics_for_language_1',
        req,
        onSuccess,
        onError,
    );
}
