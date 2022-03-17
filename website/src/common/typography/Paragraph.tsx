import React from 'react';
import classNames from 'classnames';

import { makeStyles } from '@material-ui/core/styles';

import { asTypography, TypographyProps } from './common';

const styleClasses = makeStyles({
    paragraphStyle: (props: ParagraphProps) => ({
        fontFamily: "'Roboto', sans-serif",
        lineHeight: 1.5,
        fontSize: props.size != null ? props.size : Size.Medium,
        fontWeight: !!props.isBold ? 700 : 400,
    }),
})

export enum Size {
    Small = '12px',
    Medium = '16px',
    Large = '20px',
    ExtraLarge = '24px',
}

type ParagraphProps = {
    size?: Size;
} & TypographyProps;

const Paragraph = asTypography((props: ParagraphProps) => {
    const classes = styleClasses(props);
    const className = classNames(
        'Paragraph__root',
        props.className,
        `${classes.paragraphStyle}`
    );
    return (
        <p className={className}>{props.children}</p>
    );
});

export default Paragraph;
