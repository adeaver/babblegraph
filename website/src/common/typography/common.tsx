import './common.scss';

import React from 'react';
import classNames from 'classnames';

export enum Alignment {
    Center,
    Left,
    Right,
}

export enum Color {
    Primary,
    Secondary,
    Black,
    Gray,
    White,
}

export type TypographyProps = {
    children: React.ReactNode,
    className?: string,
    align?: Alignment,
    color?: Color,
}

export function asTypography<P extends TypographyProps>(WrappedComponent: React.ComponentType<P>) {
    return class extends React.Component<P> {
        render() {
            const className = classNames(this.props.className, {
                'Typography__centered': this.props.align == null || this.props.align === Alignment.Center,
                'Typography__right-aligned': this.props.align === Alignment.Right,
                'Typography__left-aligned': this.props.align === Alignment.Left,
                'Typography__gray-color': this.props.color == null || this.props.color === Color.Gray,
                'Typography__black-color': this.props.color === Color.Black,
                'Typography__white-color': this.props.color === Color.White,
                'Typography__primary-color': this.props.color === Color.Primary,
                'Typography__secondary-color': this.props.color === Color.Secondary,
            })
            return <WrappedComponent className={className} {...this.props as P} />
        }
    }
}
