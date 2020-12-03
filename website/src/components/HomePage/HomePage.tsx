import './HomePage.scss';

import React from 'react';

import { Heading1, Heading3 } from 'common/typography/Heading.tsx';

export default class HomePage extends React.Component{
    render() {
        return (
            <div className="HomePage__root">
                <Heading1 className="HomePage__heading">Babblegraph</Heading1>
                <Heading3 className="HomePage__subheading">
                    Practice your Spanish and stay up-to-date <br />
                    with the latest news from the Spanish speaking world.<br />
                    With one email a day.
                </Heading3>
            </div>
        );
    }
}
