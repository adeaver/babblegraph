import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import { asBasePage, BasePageProps } from 'AdminWeb/common/BasePage/BasePage';
import { Heading1 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';

import {
    Topic,
    GetAllContentTopicsResponse,
    getAllContentTopics,

    AddTopicResponse,
    addTopic,
} from 'AdminWeb/api/content/topic';

type TopicListPageProps = GetAllContentTopicsResponse & BasePageProps;

const TopicListPage = (props: TopicListPageProps) => {
    return (
        <Heading1>{props.topics}</Heading1>
    );
}

export default asBasePage(
    TopicListPage,
    (
        setData: (data: GetAllContentTopicsResponse) => void,
        setError: (err: Error) => void,
    ) => {
        getAllContentTopics(
            {},
            setData,
            setError,
        );
    },
);
