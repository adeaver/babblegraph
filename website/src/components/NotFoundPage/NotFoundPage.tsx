import './NotFoundPage.scss';

import React from 'react';

import Header from 'common/components/Header/Header';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';

export default class NotFoundPage extends React.Component{
    render() {
        return (
            <div className="NotFoundPage__root">
                <Header />
                <Heading1 className="NotFoundPage__heading">Babblegraph</Heading1>
                <Heading3 className="NotFoundPage__subheading">
                    We couldn't find that page...
                </Heading3>
                <Paragraph className="NotFoundPage__paragraph">Whatever you’re looking for isn't here...</Paragraph>
            </div>
        );
    }
}
