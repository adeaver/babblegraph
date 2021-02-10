import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type SetPageLoadEventRequest = {};

export type SetPageLoadEventResponse = {};

export function setPageLoadEvent(
    req: SetPageLoadEventRequest,
    onSuccess: (resp: SetPageLoadEventResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<SetPageLoadEventRequest, SetPageLoadEventResponse>(
        '/api/utm/set_page_load_event_1',
        req,
        onSuccess,
        onError,
    );
}
