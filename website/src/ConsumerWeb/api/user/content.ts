import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { ClientError } from 'ConsumerWeb/api/clienterror';
import { WordsmithLanguageCode } from 'common/model/language/language';

export type GetUserContentTopicsRequest = {
    subscriptionManagementToken: string;
}

export type GetUserContentTopicsResponse = {
    error: ClientError | undefined;
    topicIds: Array<string>;
}

export function getUserContentTopics(
    req: GetUserContentTopicsRequest,
    onSuccess: (resp: GetUserContentTopicsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserContentTopicsRequest, GetUserContentTopicsResponse>(
        '/api/user/get_user_content_topics_for_token_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpsertUserContentTopicsRequest = {
    emailAddress: string | undefined;
    subscriptionManagementToken: string;
    topicIds: Array<string>;
    languageCode: WordsmithLanguageCode;
}

export type UpsertUserContentTopicsResponse = {
    success: boolean;
    error: ClientError | undefined;
}

export function upsertUserContentTopics(
    req: UpsertUserContentTopicsRequest,
    onSuccess: (resp: UpsertUserContentTopicsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpsertUserContentTopicsRequest, UpsertUserContentTopicsResponse>(
        '/api/user/upsert_user_content_topics_for_token_1',
        req,
        onSuccess,
        onError,
    );
}
