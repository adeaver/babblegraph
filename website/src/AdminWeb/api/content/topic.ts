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
    topicId: string;
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

export type TopicDisplayName = {
	id: string;
	topicId: string;
	languageCode: string;
	label: string;
	isActive: boolean;
}

export type GetAllTopicDisplayNamesForTopicRequest = {
    topicId: string;
}

export type GetAllTopicDisplayNamesForTopicResponse = {
    topicDisplayNames: Array<TopicDisplayName>;
}

export function getAllTopicDisplayNamesForTopic(
    req: GetAllTopicDisplayNamesForTopicRequest,
    onSuccess: (resp: GetAllTopicDisplayNamesForTopicResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetAllTopicDisplayNamesForTopicRequest, GetAllTopicDisplayNamesForTopicResponse>(
        '/ops/api/content/get_all_topic_display_names_for_topic_1',
        req,
        onSuccess,
        onError,
    );
}

export type AddTopicDisplayNameForTopicRequest = {
    topicId: string;
    label: string;
    languageCode: string;
}

export type AddTopicDisplayNameForTopicResponse = {
    topicDisplayNameId: string;
}

export function addTopicDisplayNameForTopic(
    req: AddTopicDisplayNameForTopicRequest,
    onSuccess: (resp: AddTopicDisplayNameForTopicResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<AddTopicDisplayNameForTopicRequest, AddTopicDisplayNameForTopicResponse>(
        '/ops/api/content/add_topic_display_name_for_topic_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateTopicDisplayNameLabelRequest = {
    topicDisplayNameId: string;
    label: string;
}

export type UpdateTopicDisplayNameLabelResponse = {
    success: boolean;
}

export function updateTopicDisplayNameLabel(
    req: UpdateTopicDisplayNameLabelRequest,
    onSuccess: (resp: UpdateTopicDisplayNameLabelResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateTopicDisplayNameLabelRequest, UpdateTopicDisplayNameLabelResponse>(
        '/ops/api/content/update_topic_display_name_label_1',
        req,
        onSuccess,
        onError,
    );
}

export type ToggleTopicDisplayNameIsActiveRequest = {
    topicDisplayNameId: string;
    isActive: boolean;
}

export type ToggleTopicDisplayNameIsActiveResponse = {
    success: boolean;
}

export function toggleTopicDisplayNameIsActive(
    req: ToggleTopicDisplayNameIsActiveRequest,
    onSuccess: (resp: ToggleTopicDisplayNameIsActiveResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<ToggleTopicDisplayNameIsActiveRequest, ToggleTopicDisplayNameIsActiveResponse>(
        '/ops/api/content/toggle_topic_display_name_is_active_1',
        req,
        onSuccess,
        onError,
    );
}
