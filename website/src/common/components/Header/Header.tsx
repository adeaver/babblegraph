import './Header.scss';

import React from 'react';
import classNames from 'classnames';

type HeaderProps = {
    className?: string;
}

export default class Header extends React.Component<HeaderProps> {
    render() {
        return (
            <div className={classNames('Header__root', this.props.className)}>
                {this.props.children}
            </div>
        );
    }
}
