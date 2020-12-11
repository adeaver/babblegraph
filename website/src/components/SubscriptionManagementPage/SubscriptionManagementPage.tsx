import './SubscriptionManagementPage.scss';

import React, {useState} from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Header from 'common/components/Header/Header';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import Input, { InputType} from 'common/components/Input/Input';
import Button, { ButtonType } from 'common/components/Button/Button';

type Params = {
    token: string
}

type SubscriptionManagementPageInitialProps = {}

const SubscriptionManagementPageInitial = (props: SubscriptionManagementPageInitialProps) => {
    return (
        <div className="SubscriptionManagementPage__root">
            <Header className="SubscriptionManagementPage__header">
                <Heading1 className="SubscriptionManagementPage__heading">Manage your subscription</Heading1>
            </Header>
            <div className="SubscriptionManagementPage__content-container">
            </div>
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
    const { token } = props.match.params;

    return (
        <SubscriptionManagementPageInitial />
    );
}

export default SubscriptionManagementPage;
