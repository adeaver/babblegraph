import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type Topic = {
    id: string;
    label: string;
    isActive: boolean;
}

export type GetAllContentTopicsRequest = {}

export type GetAllContentTopicsResponse = {
    topics: Array<Topic>;
}

export function getAllContentTopics(
    req: GetAllContentTopicsRequest,
    onSuccess: (resp: GetAllContentTopicsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetAllContentTopicsRequest, GetAllContentTopicsResponse>(
        '/ops/api/content/get_all_topics_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetTopicByIDRequest = {
    id: string;
}

export type GetTopicByIDResponse = {
    topic: Topic;
}

export function getTopicByID(
    req: GetTopicByIDRequest,
    onSuccess: (resp: GetTopicByIDResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetTopicByIDRequest, GetTopicByIDResponse>(
        '/ops/api/content/get_topic_by_id_1',
        req,
        onSuccess,
        onError,
    );
}

export type AddTopicRequest = {
    label: string;
}

export type AddTopicResponse = {
    id: string;
}

export function addTopic(
    req: AddTopicRequest,
    onSuccess: (resp: AddTopicResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<AddTopicRequest, AddTopicResponse>(
        '/ops/api/content/add_topic_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateIsContentTopicActiveRequest = {
    id: string;
    isActive: boolean;
}

export type UpdateIsContentTopicActiveResponse = {
    success: boolean;
}

export function updateIsContentTopicActive(
    req: UpdateIsContentTopicActiveRequest,
    onSuccess: (resp: UpdateIsContentTopicActiveResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateIsContentTopicActiveRequest, UpdateIsContentTopicActiveResponse>(
        '/ops/api/content/update_is_topic_active_1',
        req,
        onSuccess,
        onError,
    );
}
