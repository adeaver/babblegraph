import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { ClientError } from 'ConsumerWeb/api/clienterror';

export type GetUserContentTopicsForTokenRequest = {
    subscriptionManagementToken: string;
}

export type GetUserContentTopicsForTokenResponse = {
    topics: string[];
    error: ClientError | undefined;
}

export function getUserContentTopicsForToken(
    req: GetUserContentTopicsForTokenRequest,
    onSuccess: (resp: GetUserContentTopicsForTokenResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserContentTopicsForTokenRequest, GetUserContentTopicsForTokenResponse>(
        '/api/user/get_user_content_topics_for_token_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateUserContentTopicsForTokenRequest = {
    subscriptionManagementToken: string;
    emailAddress: string | undefined;
    activeTopicIds: string[];
}

export type UpdateUserContentTopicsForTokenResponse = {
    success: boolean;
    error: ClientError | undefined;
};

export function updateUserContentTopicsForToken(
    req: UpdateUserContentTopicsForTokenRequest,
    onSuccess: (resp: UpdateUserContentTopicsForTokenResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateUserContentTopicsForTokenRequest, UpdateUserContentTopicsForTokenResponse>(
        '/api/user/update_user_content_topics_for_token_1',
        req,
        onSuccess,
        onError,
    );
}
