import React from 'react';
import classNames from 'classnames';

import { makeStyles } from '@material-ui/core/styles';

import Color from 'common/styles/colors';

export enum Alignment {
    Center = 'center',
    Left = 'left',
    Right = 'right',
}

export enum TypographyColor {
    Primary = Color.Primary,
    Secondary = Color.Secondary,
    Warning = Color.Warning,
    Black = Color.Black,
    Gray = Color.TextGray,
    White = Color.White,
}

const styleClasses = makeStyles({
    typographyStyle: (props: TypographyProps) => ({
        textAlign: props.align != null ? props.align : Alignment.Center,
        color: props.color != null ? props.color : TypographyColor.Gray,
    }),
});

export type TypographyProps = {
    children: React.ReactNode,
    className?: string,
    align?: Alignment,
    color?: TypographyColor,
}

export function asTypography<P extends TypographyProps>(WrappedComponent: React.ComponentType<P>) {
    return (props: P) => {
        const classes = styleClasses(props);
        const className = classNames(props.className, `${classes.typographyStyle}`);
        return <WrappedComponent className={className} {...props as P} />
    }
}
