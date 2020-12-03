import './Paragraph.scss';

import React from 'react';
import classNames from 'classnames';

type ParagraphProps = {
    className?: string;
}

export class Paragraph from React.Component<ParagraphProps> {
    render() {
        const className = classNames('Paragraph__root', this.props.className);
        return (
            <p className={className}>{this.props.children}</p>
        )
    }
}
