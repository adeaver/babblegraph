import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type UnsubscribeRequest = {
    token: string;
    unsubscribeReason: string | null;
    emailAddress: string;
}


export type UnsubscribeResponse = {
    success: boolean;
}

export const unsubscribeUser = (
    req: UnsubscribeRequest,
    onSuccess: (UnsubscribeResponse) => void,
    onError: (Error) => void
) => {
    makePostRequestWithStandardEncoding<UnsubscribeRequest, UnsubscribeResponse>(
        '/api/user/unsubscribe_user_1',
        req,
        onSuccess,
        onError,
    );
}
