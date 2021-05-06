import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type GetUserContentTopicsForTokenRequest = {
    token: string;
}

export type GetUserContentTopicsForTokenResponse = {
    contentTopics: string[];
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
    token: string;
    emailAddress: string;
    contentTopics: string[];
}

export type UpdateUserContentTopicsForTokenResponse = {};

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

export type ContentTopicDisplayMapping = {
    displayText: string;
    apiValue: string[];
}

export const contentTopicDisplayMappings: Array<ContentTopicDisplayMapping> = [
    { displayText: "Art", apiValue: ["art"] },
    { displayText: "Architecture", apiValue: ["architecture"] },
    { displayText: "Automotive", apiValue: ["automotive"] },
	{ displayText: "Business & Finance", apiValue: ["business", "finance"] },
	{ displayText: "Celebrity News", apiValue: ["celebrity-news"] },
	{ displayText: "Cooking", apiValue: ["cooking"] },
	{ displayText: "Argentina", apiValue: ["current-events-argentina"] },
	{ displayText: "Chile", apiValue: ["current-events-chile"] },
	{ displayText: "Colombia", apiValue: ["current-events-colombia"] },
	{ displayText: "Costa Rica", apiValue: ["current-events-costa-rica"] },
	{ displayText: "El Salvador", apiValue: ["current-events-el-salvador"] },
	{ displayText: "Guatemala", apiValue: ["current-events-guatemala"] },
	{ displayText: "Honduras", apiValue: ["current-events-honduras"] },
	{ displayText: "Mexico", apiValue: ["current-events-mexico"] },
	{ displayText: "Nicaragua", apiValue: ["current-events-nicaragua"] },
	{ displayText: "Panama", apiValue: ["current-events-panama"] },
	{ displayText: "Peru", apiValue: ["current-events-peru"] },
	{ displayText: "Paraguay", apiValue: ["current-events-paraguay"] },
	{ displayText: "Spain", apiValue: ["current-events-spain"] },
	{ displayText: "United States", apiValue: ["current-events-united-states"] },
	{ displayText: "Venezuela", apiValue: ["current-events-venezuela"] },
	{ displayText: "Uruguay", apiValue: ["current-events-uruguay"] },
	{ displayText: "Economy", apiValue: ["economy"] },
	{ displayText: "Entertainment", apiValue: ["entertainment", "culture"] },
	{ displayText: "Fashion", apiValue: ["fashion"] },
	{ displayText: "Film", apiValue: ["film"] },
	{ displayText: "Health", apiValue: ["health"] },
	{ displayText: "Home", apiValue: ["home"] },
	{ displayText: "Lifestyle", apiValue: ["lifestyle"] },
	{ displayText: "Literature", apiValue: ["literature"] },
	{ displayText: "Music", apiValue: ["music"] },
	{ displayText: "Opinion", apiValue: ["opinion"] },
	{ displayText: "Politics", apiValue: ["politics"] },
	{ displayText: "Sports", apiValue: ["sports"] },
	{ displayText: "Science & Technology", apiValue: ["science", "technology"] },
	{ displayText: "Theater", apiValue: ["theater"] },
	{ displayText: "Travel", apiValue: ["travel"] },
	{ displayText: "Video Games", apiValue: ["video-games"] },
	{ displayText: "World News", apiValue: ["world-news"] },
    { displayText: "Environment", apiValue: ["environment"] },
];
