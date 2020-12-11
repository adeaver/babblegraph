import './SubscriptionManagementPage.scss';

import React, { useState, useEffect } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Header from 'common/components/Header/Header';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import Input, { InputType} from 'common/components/Input/Input';
import Button, { ButtonType } from 'common/components/Button/Button';
import Spinner, { RingColor } from 'common/ui/icons/Spinner/Spinner';
import Selector from 'common/components/Selector/Selector';

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
    updateReadingClassifications: (classifications: Array<ReadingLevelClassificationForLanguage>) => void;

    emailAddress: string | null;
    handleEmailUpdate: (string) => void;
}

const SubscriptionManagementPageInitial = (props: SubscriptionManagementPageInitialProps) => {
    const handleSelectNewReadingLevel = (language: string) => {
        return (newReadingLevel: string) => {
            props.updateReadingClassifications(
                props.readingClassifications.map((classification: ReadingLevelClassificationForLanguage) => {
                    const readingLevel = classification.languageCode ? newReadingLevel : classification.readingLevelClassification;
                    return {
                        ...classification,
                        readingLevelClassification: readingLevel,
                    };
                })
            );
        }
    }
    return (
        <div className="SubscriptionManagementPageInitial__root">
            <Paragraph>You can change what reading level you receive emails at here.</Paragraph>
            <Selector
                className="SubscriptionManagementPageInitial__selector"
                options={["Beginner", "Intermediate", "Advanced", "Professional"]}
                initialValue={props.readingClassifications[0].readingLevelClassification}
                onValueChange={handleSelectNewReadingLevel("es")} />
            <Paragraph>Confirm your email to submit changes</Paragraph>
            <Input className="SubscriptionManagementPageInitial__input" type={InputType.EMAIL} value={props.emailAddress} onChange={props.handleEmailUpdate} placeholder="Email address" />
            <Button
                onClick={() => { console.log("clicked") }}
                className="SubscriptionManagementPageInitial__submit-button"
                isLoading={false}
                type={ButtonType.Primary}
                isDisabled={props.emailAddress == null}>
                Submit
            </Button>
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
    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ readingLevelClassifications, setReadingLevelClassifications ] = useState<Array<ReadingLevelClassificationForLanguage>>([]);
    const { token } = props.match.params;

    useEffect(() => {
        if (!readingLevelClassifications.length) {
            getUserPreferencesForToken({
                token: token,
            },
            (resp: GetUserPreferencesForTokenResponse) => {
                setReadingLevelClassifications(resp.classificationsByLanguage);
                setIsLoading(false);
            },
            (e: Error) => {
                setIsLoading(false);
            });
        }
    });

    let body = (
        <SubscriptionManagementPageInitial
            emailAddress={emailAddress}
            handleEmailUpdate={setEmailAddress}
            updateReadingClassifications={setReadingLevelClassifications}
            readingClassifications={readingLevelClassifications} />
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
