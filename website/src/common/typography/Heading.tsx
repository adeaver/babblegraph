import './Heading.scss';
import './common.scss';

import React from 'react';
import classNames from 'classnames';

import { Alignment, TypographyProps } from './common.ts';

type HeadingProps = TypographyProps

export class Heading1 extends React.Component<HeadingProps> {
    render() {
        const className = classNames('Heading1__root', 'Heading__root', this.props.className, {
            'Typography__centered': this.props.align == null || this.props.align === Alignment.Center,
            'Typography__right-aligned': this.props.align === Alignment.Right,
            'Typography__left-aligned': this.props.align == Alignment.Left,
        });
        return (
            <h1 className={className}>{this.props.children}</h1>
        );
    }
}

export class Heading2 extends React.Component<HeadingProps> {
    render() {
        const className = classNames('Heading2__root', 'Heading__root', this.props.className, {
            'Typography__centered': this.props.align == null || this.props.align === Alignment.Center,
            'Typography__right-aligned': this.props.align === Alignment.Right,
            'Typography__left-aligned': this.props.align == Alignment.Left,
        });
        return (
            <h2 className={className}>{this.props.children}</h2>
        );
    }
}

export class Heading3 extends React.Component<HeadingProps> {
    render() {
        const className = classNames('Heading3__root', 'Heading__root', this.props.className, {
            'Typography__centered': this.props.align == null || this.props.align === Alignment.Center,
            'Typography__right-aligned': this.props.align === Alignment.Right,
            'Typography__left-aligned': this.props.align == Alignment.Left,
        });
        return (
            <h3 className={className}>{this.props.children}</h3>
        );
    }
}

export class Heading4 extends React.Component<HeadingProps> {
    render() {
        const className = classNames('Heading4__root', 'Heading__root', this.props.className, {
            'Typography__centered': this.props.align == null || this.props.align === Alignment.Center,
            'Typography__right-aligned': this.props.align === Alignment.Right,
            'Typography__left-aligned': this.props.align == Alignment.Left,
        });
        return (
            <h4 className={className}>{this.props.children}</h4>
        );
    }
}

export class Heading5 extends React.Component<HeadingProps> {
    render() {
        const className = classNames('Heading5__root', 'Heading__root', this.props.className, {
            'Typography__centered': this.props.align == null || this.props.align === Alignment.Center,
            'Typography__right-aligned': this.props.align === Alignment.Right,
            'Typography__left-aligned': this.props.align == Alignment.Left,
        });
        return (
            <h5 className={className}>{this.props.children}</h5>
        );
    }
}

export class Heading6 extends React.Component<HeadingProps> {
    render() {
        const className = classNames('Heading6__root', 'Heading__root', this.props.className, {
            'Typography__centered': this.props.align == null || this.props.align === Alignment.Center,
            'Typography__right-aligned': this.props.align === Alignment.Right,
            'Typography__left-aligned': this.props.align == Alignment.Left,
        });
        return (
            <h6 className={className}>{this.props.children}</h6>
        );
    }
}
