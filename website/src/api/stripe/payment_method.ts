import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type GetSetupIntentForUserRequest = {}

export type GetSetupIntentForUserResponse = {
    setupIntentId: string;
    clientSecret: string;
}

export function getSetupIntentForUser(
    req: GetSetupIntentForUserRequest,
    onSuccess: (resp: GetSetupIntentForUserResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetSetupIntentForUserRequest, GetSetupIntentForUserResponse>(
        '/api/stripe/get_setup_intent_for_user_1',
        req,
        onSuccess,
        onError,
    );
}
