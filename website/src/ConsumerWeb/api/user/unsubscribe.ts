import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type UnsubscribeRequest = {
    token: string;
    unsubscribeReason: string | null;
    emailAddress: string | undefined;
}


export type UnsubscribeResponse = {
    success: boolean;
    error: UnsubscribeError | undefined;
}

export enum UnsubscribeError {
    MissingEmail = 'missing-email',
    IncorrectEmail = 'incorrect-email',
    NoAuth = 'no-auth',
    IncorrectKey = 'incorrect-key',
    InvalidToken = 'invalid-token',
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
