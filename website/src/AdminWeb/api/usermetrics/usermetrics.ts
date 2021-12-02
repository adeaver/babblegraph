import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type GetUserAggregationByStatusRequest = {}

export type GetUserAggregationByStatusResponse = {
    verifiedUserCount: number;
    unsubscribedUserCount: number;
    unverifiedUserCount: number;
    blocklistedUserCount: number;
}

export function getUserAggregationByStatus(
    req: GetUserAggregationByStatusRequest,
    onSuccess: (resp: GetUserAggregationByStatusResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserAggregationByStatusRequest, GetUserAggregationByStatusResponse>(
        '/ops/api/usermetrics/get_user_aggregation_by_status_1',
        req,
        onSuccess,
        onError,
    );
}
