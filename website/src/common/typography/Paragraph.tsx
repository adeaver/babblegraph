import './Paragraph.scss';
import './common.scss';

import React from 'react';
import classNames from 'classnames';

import { Alignment, Color, TypographyProps } from './common.ts';

export enum Size {
    Small,
    Medium,
    Large,
    ExtraLarge,
}

type ParagraphProps = {
    size?: Size;
} & TypographyProps;

export default class Paragraph extends React.Component<ParagraphProps> {
    render() {
        const className = classNames('Paragraph__root', this.props.className, {
            'Typography__centered': this.props.align == null || this.props.align === Alignment.Center,
            'Typography__right-aligned': this.props.align === Alignment.Right,
            'Typography__left-aligned': this.props.align === Alignment.Left,
            'Typography__gray-color': this.props.color == null || this.props.color === Color.Gray,
            'Typography__black-color': this.props.color === Color.Black,
            'Typography__white-color': this.props.color === Color.White,
            'Typography__primary-color': this.props.color === Color.Primary,
            'Typography__secondary-color': this.props.color === Color.Secondary,
            'Paragraph__size-small': this.props.size === Size.Small,
            'Paragraph__size-medium': this.props.size == null || this.props.size === Size.Medium,
            'Paragraph__size-large': this.props.size === Size.Large,
            'Paragraph__size-extra-large': this.props.size === Size.ExtraLarge,
        });
        return (
            <p className={className}>{this.props.children}</p>
        );
    }
}
