import {
    encodeToAPIObject,
    decodeAPIObject
} from 'api/bgfetch/apiencoding';

export function makePostRequest<T, U>(
    url: string,
    body: T,
    onSuccess: (resp: U) => void,
    onError: (e: Error) => void,
) {
    fetch(url, {
        method: 'POST',
        cache: 'no-cache',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    })
    .then(response => {
        if (!response.ok) {
            response.text().then(data => onError(new Error(data)));
            return
        }
        response.json().then(data => onSuccess(data as U));
    })
    .catch(onError);
}

export function makePostRequestWithStandardEncoding<T, U>(
    url: string,
    body: T,
    onSuccess: (resp: U) => void,
    onError: (e: Error) => void,
) {
    makePostRequest<object, object>(
        url,
        encodeToAPIObject<T>(body),
        (resp: object) => { onSuccess(decodeAPIObject<U>(resp)) },
        onError,
    );
}
