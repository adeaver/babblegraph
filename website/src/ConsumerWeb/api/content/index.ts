import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { ClientError } from 'ConsumerWeb/api/clienterror';
import {
    Topic,
    TopicDisplayName,
} from 'common/api/content';

import { WordsmithLanguageCode } from 'common/model/language/language';

export type TopicWithDisplay = {
    topic: Topic;
    displayName: TopicDisplayName;
}

export type GetActiveTopicsForLanguageCodeRequest = {
    languageCode: WordsmithLanguageCode;
}

export type GetActiveTopicsForLanguageCodeResponse = {
    topics: Array<TopicWithDisplay> | undefined;
    error: ClientError | undefined;
}

export function getActiveTopicsForLanguageCode(
    req: GetActiveTopicsForLanguageCodeRequest,
    onSuccess: (resp: GetActiveTopicsForLanguageCodeResponse) => void,
    onError: (error: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetActiveTopicsForLanguageCodeRequest, GetActiveTopicsForLanguageCodeResponse>(
        '/api/content/get_active_topics_for_language_code_1',
        req,
        onSuccess,
        onError,
    );
}
