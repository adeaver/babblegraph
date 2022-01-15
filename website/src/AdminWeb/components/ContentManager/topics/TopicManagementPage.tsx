import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { asBaseComponent, BaseComponentProps } from 'AdminWeb/common/Base/BaseComponent';

import {
    getTopicByID,
    GetTopicByIDResponse,
} from 'AdminWeb/api/content/topic';

type Params = {
    id: string;
}

type TopicManagementPageOwnProps = RouteComponentProps<Params>;

const TopicManagementPage = asBaseComponent<GetTopicByIDResponse, TopicManagementPageOwnProps>(
    (props: TopicManagementPageOwnProps & GetTopicByIDResponse & BaseComponentProps) => {
        return <h1>{props.match.params.id}</h1>
    },
    (
        ownProps: TopicManagementPageOwnProps,
        onSuccess: (resp: GetTopicByIDResponse) => void,
        onError: (err: Error) => void,
    ) => {
        getTopicByID({
            id: ownProps.match.params.id,
        },
        onSuccess,
        onError)
    },
    true
);

export default TopicManagementPage;
