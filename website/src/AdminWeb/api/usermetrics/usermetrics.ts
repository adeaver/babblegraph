import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type GetUserStatusDataRequest = {}

export type GetUserStatusDataResponse = {
    verifiedUserCount: number;
    unsubscribedUserCount: number;
    unverifiedUserCount: number;
    blocklistedUserCount: number;
    verifiedUserCountNetChangeOverWeek: number;
    verifiedUserCountNetChangeOverMonth: number;
}

export function getUserStatusData(
    req: GetUserStatusDataRequest,
    onSuccess: (resp: GetUserStatusDataResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserStatusDataRequest, GetUserStatusDataResponse>(
        '/ops/api/usermetrics/get_user_status_data_1',
        req,
        onSuccess,
        onError,
    );
}
