import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { ClientError } from 'ConsumerWeb/api/clienterror';

export type GetPodcastMetadataRequest = {
	userPodcastId: string;
}

export type GetPodcastMetadataResponse = {
    error: ClientError | undefined;
    metadata: PodcastMetadata | undefined;
}

export type PodcastMetadata = {
    podcastTitle: string;
    episodeTitle: string;
    episodeDescription: string;
    podcastUrl: string;
    audioUrl: string;
    imageUrl: string | undefined;
}

export function getPodcastMetadata(
    req: GetPodcastMetadataRequest,
    onSuccess: (resp: GetPodcastMetadataResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetPodcastMetadataRequest, GetPodcastMetadataResponse>(
        '/api/podcasts/get_podcast_metadata_1',
        req,
        onSuccess,
        onError,
    );
}
