import React from 'react';
import classNames from 'classnames';

import { makeStyles } from '@material-ui/core/styles';

import { asTypography, Alignment, TypographyProps } from './common';

const styleClasses = makeStyles({
    headingRoot: {
        fontFamily: "'Quicksand', sans-serif",
    },
})

type HeadingProps = TypographyProps

export const Heading1 = asTypography((props: HeadingProps) => {
    const classes = styleClasses();
    const className = classNames(`${classes.headingRoot}`, props.className);
    return (
        <h1 className={className}>{props.children}</h1>
    );
});

export const Heading2 = asTypography((props: HeadingProps) => {
    const classes = styleClasses();
    const className = classNames(`${classes.headingRoot}`, props.className);
    return (
        <h2 className={className}>{props.children}</h2>
    );
});

export const Heading3 = asTypography((props: HeadingProps) => {
    const classes = styleClasses();
    const className = classNames(`${classes.headingRoot}`, props.className);
    return (
        <h3 className={className}>{props.children}</h3>
    );
});

export const Heading4 = asTypography((props: HeadingProps) => {
    const classes = styleClasses();
    const className = classNames(`${classes.headingRoot}`, props.className);
    return (
        <h4 className={className}>{props.children}</h4>
    );
});

export const Heading5 = asTypography((props: HeadingProps) => {
    const classes = styleClasses();
    const className = classNames(`${classes.headingRoot}`, props.className);
    return (
        <h5 className={className}>{props.children}</h5>
    );
});

export const Heading6 = asTypography((props: HeadingProps) => {
    const classes = styleClasses();
    const className = classNames(`${classes.headingRoot}`, props.className);
    return (
        <h6 className={className}>{props.children}</h6>
    );
});
