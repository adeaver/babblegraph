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

export type SearchPodcastsParams = {
    language: string;
    region: string;
    genre: number;
    pageNumber: number | undefined;
}

export type PodcastMetadata = {
    title: string;
    country: string;
    description: string;
    website: string;
    language: string;
    type: string;
    totalNumberOfEpisodes: number;
    listenNotesUrl: string;
}

export type SearchPodcastsRequest = {
    params: SearchPodcastsParams;
}

export type SearchPodcastsResponse = {
    podcasts: Array<PodcastMetadata>;
    nextPageNumber: number | undefined;
}

export function searchPodcasts(
    req: SearchPodcastsRequest,
    onSuccess: (resp: SearchPodcastsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<SearchPodcastsRequest, SearchPodcastsResponse>(
        '/ops/api/podcasts/search_podcasts_1',
        req,
        onSuccess,
        onError,
    );
}
