import { makePostRequest } from 'api/bgfetch/bgfetch';

export type UnsubscribeRequest = {
    UserID: string;
    EmailAddress: string;
}

type apiEncodedUnsubscribeRequest = {
    user_id: string;
    email_address: string;
}

const makeAPIEncodedUnsubscribeRequest = (req: UnsubscribeRequest) => ({
    user_id: req.UserID,
    email_address: req.EmailAddress,
});

export type UnsubscribeResponse = {
    Success: boolean;
}

type apiEncodedUnsubscribeResponse = {
    success: boolean;
}

const unencodeUnsubscribeResponse = (apiEncoded: apiEncodedUnsubscribeResponse) => {
    return {
        Success: apiEncoded.success,
    };
}

export const UnsubscribeUser = (
    req: UnsubscribeRequest,
    onSuccess: (UnsubscribeResponse) => void,
    onError: (Error) => void
) => {
    makePostRequest<apiEncodedUnsubscribeRequest, apiEncodedUnsubscribeResponse>(
        '/api/user/unsubscribe_user_1',
        makeAPIEncodedUnsubscribeRequest(req),
        (resp: apiEncodedUnsubscribeResponse) => {
            onSuccess(unencodeUnsubscribeResponse(resp));
        },
        onError,
    );
}
