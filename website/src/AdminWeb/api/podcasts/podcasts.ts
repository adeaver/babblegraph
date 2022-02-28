import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type SearchOptions = {
    supportedLanguages: Array<string>;
    supportedRegions: Array<SupportedRegion>;
    genres: Array<SupportedGenre>;
}

export type SupportedRegion = {
    displayName: string;
    apiValue: string;
}

export type SupportedGenre = {
    displayName: string;
    apiValue: number;
}

export type GetPodcastSearchOptionsRequest = {}

export type GetPodcastSearchOptionsResponse = {
    options: SearchOptions;
}

export function getPodcastSearchOptions(
    req: GetPodcastSearchOptionsRequest,
    onSuccess: (resp: GetPodcastSearchOptionsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetPodcastSearchOptionsRequest, GetPodcastSearchOptionsResponse>(
        '/ops/api/podcasts/get_podcast_search_options_1',
        req,
        onSuccess,
        onError,
    );
}
