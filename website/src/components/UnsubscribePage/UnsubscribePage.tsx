import './UnsubscribePage.scss';

import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Header from 'common/components/Header/Header';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';

type Params = {
    userID: string
}

type UnsubscribePageProps = RouteComponentProps<Params>

export default class UnsubscribePage extends React.Component<UnsubscribePageProps> {
    render() {
        return (
            <div className="UnsubscribePage__root">
                <Header className="UnsubscribePage__header">
                        <Heading1 className="UnsubscribePage__heading">Unsubscribe from Babblegraph</Heading1>
                        <Heading3 className="UnsubscribePage__subheading">
                            We’re sorry to see you go
                        </Heading3>
                </Header>
                <Paragraph className="UnsubscribePage__explanation">By unsubscribing, you will no longer receive daily emails from Babblegraph.</Paragraph>
            </div>
        );
    }
}
