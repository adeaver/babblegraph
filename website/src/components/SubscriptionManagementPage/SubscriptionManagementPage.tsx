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
    getUserPreferencesForToken,
    UpdateUserPreferencesForTokenRequest,
    UpdateUserPreferencesForTokenResponse,
    updateUserPreferencesForToken
} from 'api/user/management';

type Params = {
    token: string
}

type SubscriptionManagementPageBodyProps = {
    readingClassifications: Array<ReadingLevelClassificationForLanguage>;
    updateReadingClassifications: (classifications: Array<ReadingLevelClassificationForLanguage>) => void;

    emailAddress: string | null;
    handleEmailUpdate: (string) => void;

    handleSubmit: () => void;
}

const SubscriptionManagementPageBody = (props: SubscriptionManagementPageBodyProps) => {
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
        <div className="SubscriptionManagementPageBody__root">
            <Paragraph>You can change what reading level you receive emails at here.</Paragraph>
            <Selector
                className="SubscriptionManagementPageBody__selector"
                options={["Beginner", "Intermediate", "Advanced", "Professional"]}
                initialValue={props.readingClassifications[0].readingLevelClassification}
                onValueChange={handleSelectNewReadingLevel("es")} />
            <Paragraph>Confirm your email to submit changes</Paragraph>
            <Input className="SubscriptionManagementPageBody__input" type={InputType.EMAIL} value={props.emailAddress} onChange={props.handleEmailUpdate} placeholder="Email address" />
            <Button
                onClick={props.handleSubmit}
                className="SubscriptionManagementPageBody__submit-button"
                isLoading={false}
                type={ButtonType.Primary}
                isDisabled={props.emailAddress == null}>
                Submit
            </Button>
        </div>
    )
}

type SubscriptionManagementPageProps = RouteComponentProps<Params>

const SubscriptionManagementPage = (props: SubscriptionManagementPageProps) => {
    const [ emailAddress, setEmailAddress ] = useState<string | null>(null);
    const [ isLoading, setIsLoading ] = useState<boolean>(true);
    const [ readingLevelClassifications, setReadingLevelClassifications ] = useState<Array<ReadingLevelClassificationForLanguage>>([]);
    const [ error, setError ] = useState<Error | null>(null);
    const [ didUpdate, setDidUpdate ] = useState<boolean | null>(null);

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
                setError(e);
            });
        }
    });

    const handleSubmit = () => {
        setIsLoading(true);
        setError(null);
        setDidUpdate(null);
        updateUserPreferencesForToken({
            token: token,
            emailAddress: emailAddress || '',
            classificationsByLanguage: readingLevelClassifications,
        },
        (resp: UpdateUserPreferencesForTokenResponse) => {
            setIsLoading(false);
            setDidUpdate(resp.didUpdate);
        },
        (e: Error) => {
            setIsLoading(false);
            setError(e);
        });
    }

    let subheading = null;
    if (didUpdate) {
        subheading = "Your subscription has been updated successfully!"
    } else if (didUpdate != null || error) {
        subheading = "Something wasn't quite right. Make sure the email is correct."
    }
    return (
        <div className="SubscriptionManagementPage__root">
            <Header className="SubscriptionManagementPage__header">
                <Heading1 className="SubscriptionManagementPage__heading">Manage your subscription</Heading1>
                { subheading && <Paragraph className="SubscriptionManagementPage__heading">{subheading}</Paragraph> }
            </Header>
            <div className="SubscriptionManagementPage__content-container">
                {
                    isLoading ? (
                        <Spinner sizeInPx={200} innerRingColor={RingColor.PrimaryPurple} outerRingColor={RingColor.PrimaryPurple} />
                    ) : (
                        <SubscriptionManagementPageBody
                            handleSubmit={handleSubmit}
                            emailAddress={emailAddress}
                            handleEmailUpdate={setEmailAddress}
                            updateReadingClassifications={setReadingLevelClassifications}
                            readingClassifications={readingLevelClassifications} />
                    )
                }
            </div>
        </div>
    );
}

export default SubscriptionManagementPage;
