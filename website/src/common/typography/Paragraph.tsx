import './Paragraph.scss';
import './common.scss';

import React from 'react';
import classNames from 'classnames';

import { Alignment, TypographyProps } from './common.ts';

type ParagraphProps = TypographyProps;

export default class Paragraph extends React.Component<ParagraphProps> {
    render() {
        const className = classNames('Paragraph__root', this.props.className, {
            'Typography__centered': this.props.align == null || this.props.align === Alignment.Center,
            'Typography__right-aligned': this.props.align === Alignment.Right,
            'Typography__left-aligned': this.props.align === Alignment.Left,
        });
        return (
            <p className={className}>{this.props.children}</p>
        );
    }
}
