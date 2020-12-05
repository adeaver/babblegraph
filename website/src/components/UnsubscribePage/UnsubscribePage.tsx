import './UnsubscribePage.scss';

import React, {useState} from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Header from 'common/components/Header/Header';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import Input, { InputType} from 'common/components/Input/Input';
import Button, { ButtonType } from 'common/components/Button/Button';

import { UnsubscribeUser, UnsubscribeRequest, UnsubscribeResponse } from 'api/user/unsubscribe';

type Params = {
    userID: string
}

type UnsubscribePageInitialProps = {
    userID: string;
    email: string | undefined;
    isLoading: boolean;

    handleEmailUpdate: (newEmail: string) => void;
    handleClick: () => void;
}

const UnsubscribePageInitial = (props: UnsubscribePageInitialProps) => {
    return (
        <div className="UnsubscribePage__root">
            <Header className="UnsubscribePage__header">
                <Heading1 className="UnsubscribePage__heading">Unsubscribe from Babblegraph</Heading1>
                <Heading3 className="UnsubscribePage__subheading">
                    We’re sorry to see you go
                </Heading3>
            </Header>
            <div className="UnsubscribePage__content-container">
                    <Paragraph className="UnsubscribePage__explanation">
                        By unsubscribing, you will no longer receive any emails from Babblegraph, including:<br />
                        • Marketing emails<br />
                        • Daily emails
                    </Paragraph>
                    <Input className="UnsubscribePage__email-input" type={InputType.EMAIL} value={props.email} onChange={props.handleEmailUpdate} placeholder="Email address" />
                    <Button
                        onClick={props.handleClick}
                        className="UnsubscribePage__submit-button"
                        isLoading={props.isLoading}
                        type={ButtonType.Primary}>
                        Submit
                    </Button>
            </div>
        </div>
    )
}

type UnsubscribePageRequestSuccessfulProps = {}

const UnsubscribePageRequestSuccessful = (props: UnsubscribePageRequestSuccessfulProps) => {
    return (
        <div className="UnsubscribePage__root">
            <Header className="UnsubscribePage__header">
                <Heading1 className="UnsubscribePage__heading">Unsubscribe from Babblegraph</Heading1>
                <Heading3 className="UnsubscribePage__subheading">
                    You’ve been successfully unsubscribed from Babblegraph.
                </Heading3>
            </Header>
            <div className="UnsubscribePage__content-container">
                    <Paragraph className="UnsubscribePage__explanation">
                        You will no longer receive any emails from us<br />
                        You can always rejoin by entering your email again on the home page.
                    </Paragraph>
            </div>
        </div>
    );
}

type UnsubscribePageRequestFailedProps = {
    userID: string;
    email: string | undefined;
    isLoading: boolean;

    handleEmailUpdate: (newEmail: string) => void;
    handleClick: () => void;
}

const UnsubscribePageRequestFailed = (props: UnsubscribePageRequestFailedProps) => {
    return (
        <div className="UnsubscribePage__root">
            <Header className="UnsubscribePage__header">
                <Heading1 className="UnsubscribePage__heading">Unsubscribe from Babblegraph</Heading1>
                <Heading3 className="UnsubscribePage__subheading">
                    Something went wrong processing your request
                </Heading3>
            </Header>
            <div className="UnsubscribePage__content-container">
                    <Paragraph className="UnsubscribePage__explanation">
                        We couldn’t process your request to unsubscribe. Here are some things to try:<br />
                        • Make sure that your email address is entered correctly<br />
                        • Click the unsubscribe link in your inbox again<br />
                        • Send us an email
                    </Paragraph>
                    <Input className="UnsubscribePage__email-input" type={InputType.EMAIL} value={props.email} onChange={props.handleEmailUpdate} placeholder="Email address" />
                    <Button
                        onClick={props.handleClick}
                        className="UnsubscribePage__submit-button"
                        isLoading={props.isLoading}
                        type={ButtonType.Primary}>
                        Submit
                    </Button>
            </div>
        </div>
    );
}

const handleSubmit = (
    userID: string,
    emailAddress: string,
    setIsLoading: (boolean) => void,
    onSuccess: (didUpdate: boolean | null) => void,
    onError: (e: Error) => void,
) => {
    return () => {
        setIsLoading(true);
        const req: UnsubscribeRequest = {
            UserID: userID,
            EmailAddress: emailAddress,
        }
        UnsubscribeUser(req,
            (resp: UnsubscribeResponse) => {
                setIsLoading(false);
                onSuccess(resp.Success);
            },
            (e: Error) => {
                setIsLoading(false);
                onError(e);
            },
        );
    }
}


type UnsubscribePageProps = RouteComponentProps<Params>

const UnsubscribePage = (props: UnsubscribePageProps) => {
    const [email, setEmail] = useState<string | undefined>(undefined);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [requestSuccessful, setSuccess] = useState<boolean | null>(null);
    const [requestErrored, setError] = useState<Error | null>(null);
    const { userID } = props.match.params;

    if (!!requestSuccessful) {
        return <UnsubscribePageRequestSuccessful />
    } else if ((requestSuccessful != null && !requestSuccessful) || requestErrored) {
        return (
            <UnsubscribePageRequestFailed
                userID={userID}
                email={email}
                isLoading={isLoading}
                handleEmailUpdate={setEmail}
                handleClick={handleSubmit(userID, email || '', setIsLoading, setSuccess, setError)} />
        );
    } else {
        return (
            <UnsubscribePageInitial
                userID={userID}
                email={email}
                isLoading={isLoading}
                handleEmailUpdate={setEmail}
                handleClick={handleSubmit(userID, email || '', setIsLoading, setSuccess, setError)} />
        );
    }
}

export default UnsubscribePage;
