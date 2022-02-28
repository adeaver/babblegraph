import React from 'react';

import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';

import {
    SearchOptions,

    GetPodcastSearchOptionsResponse,
    getPodcastSearchOptions,
} from 'AdminWeb/api/podcasts/podcasts';

const PodcastSearchPage = asBaseComponent<GetPodcastSearchOptionsResponse, {}>(
    (props: GetPodcastSearchOptionsResponse & BaseComponentProps) => {
        return <div />;
    },
    (
        ownProps: {},
        onSuccess: (resp: GetPodcastSearchOptionsResponse) => void,
        onError: (err: Error) => void,
    ) => getPodcastSearchOptions(ownProps, onSuccess, onError),
    true,
);

export default PodcastSearchPage;
