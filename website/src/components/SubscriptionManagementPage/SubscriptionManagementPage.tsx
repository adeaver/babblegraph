import './SubscriptionManagementPage.scss';

import React, { useState, useEffect } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Header from 'common/components/Header/Header';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import Input, { InputType} from 'common/components/Input/Input';
import Button, { ButtonType } from 'common/components/Button/Button';
import Spinner, { RingColor } from 'common/ui/icons/Spinner/Spinner';

import {
    GetUserPreferencesForTokenRequest,
    GetUserPreferencesForTokenResponse,
    ReadingLevelClassificationForLanguage,
    getUserPreferencesForToken
} from 'api/user/management';

type Params = {
    token: string
}

type SubscriptionManagementPageInitialProps = {
    readingClassifications: Array<ReadingLevelClassificationForLanguage>;
}

const SubscriptionManagementPageInitial = (props: SubscriptionManagementPageInitialProps) => {
    return (
        <div className="SubscriptionManagementPageInitial__root">
            <Paragraph>Your currently receiving emails at the following level: {props.readingClassifications[0].ReadingLevelClassification}</Paragraph>
        </div>
    )
}

type SubscriptionManagementPageRequestSuccessfulProps = {}

const SubscriptionManagementPageRequestSuccessful = (props: SubscriptionManagementPageRequestSuccessfulProps) => {
    return (
        <div className="SubscriptionManagementPage__root">
            <Header className="SubscriptionManagementPage__header">
                <Heading1 className="SubscriptionManagementPage__heading">Manage your subscription</Heading1>
            </Header>
            <div className="SubscriptionManagementPage__content-container">
            </div>
        </div>
    );
}

type SubscriptionManagementPageRequestFailedProps = {}

const SubscriptionManagementPageRequestFailed = (props: SubscriptionManagementPageRequestFailedProps) => {
    return (
        <div className="SubscriptionManagementPage__root">
            <Header className="SubscriptionManagementPage__header">
                <Heading1 className="SubscriptionManagementPage__heading">Manage your subscription</Heading1>
            </Header>
            <div className="SubscriptionManagementPage__content-container">
            </div>
        </div>
    );
}

type SubscriptionManagementPageProps = RouteComponentProps<Params>

const SubscriptionManagementPage = (props: SubscriptionManagementPageProps) => {
    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ readingLevelClassifications, setReadingLevelClassifications ] = useState<Array<ReadingLevelClassificationForLanguage>>([]);
    const { token } = props.match.params;

    useEffect(() => {
        getUserPreferencesForToken({
            Token: token,
        },
        (resp: GetUserPreferencesForTokenResponse) => {
            setIsLoading(false);
        },
        (e: Error) => {
            setIsLoading(false);
        });
    });

    let body = (
        <SubscriptionManagementPageInitial
            token={token}
            />
    );
    if (isLoading) {
        body = (
            <Spinner sizeInPx={200} innerRingColor={RingColor.PrimaryPurple} outerRingColor={RingColor.PrimaryPurple} />
        );
    }

    return (
        <div className="SubscriptionManagementPage__root">
            <Header className="SubscriptionManagementPage__header">
                <Heading1 className="SubscriptionManagementPage__heading">Manage your subscription</Heading1>
            </Header>
            <div className="SubscriptionManagementPage__content-container">
                { body }
            </div>
        </div>
    );
}

export default SubscriptionManagementPage;
