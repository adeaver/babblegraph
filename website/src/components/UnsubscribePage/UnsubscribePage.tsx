import './UnsubscribePage.scss';

import React, {useState} from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Header from 'common/components/Header/Header';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import Input, { InputType} from 'common/components/Input/Input';
import Button, { ButtonType } from 'common/components/Button/Button';

type Params = {
    userID: string
}

type UnsubscribePageProps = RouteComponentProps<Params>

export default (props: UnsubscribePageProps) => {
    const [email, setEmail] = useState<string | undefined>(undefined);

    return (
        <div className="UnsubscribePage__root">
            <Header className="UnsubscribePage__header">
                <Heading1 className="UnsubscribePage__heading">Unsubscribe from Babblegraph</Heading1>
                <Heading3 className="UnsubscribePage__subheading">
                    We’re sorry to see you go
                </Heading3>
            </Header>
            <div className="UnsubscribePage__content-container">
                    <Paragraph className="UnsubscribePage__explanation">By unsubscribing, you will no longer receive daily emails from Babblegraph.</Paragraph>
                    <Input type={InputType.EMAIL} value={email} onChange={setEmail} placeholder="Email address" />
                    <Button type={ButtonType.Primary}>Submit</Button>
            </div>
        </div>
    );
}
