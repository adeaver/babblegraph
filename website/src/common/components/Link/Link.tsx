import React from 'react';

import { makeStyles } from '@material-ui/core/styles';

import Color from 'common/styles/colors';
import Paragraph from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';

const styleClasses = makeStyles({
    linkAnchor: {
        textDecoration: 'none',
        color: 'inherit',
        '&:hover': {
            textDecoration: 'underline',
        }
    }
});

export enum LinkTarget {
    Blank = "_blank",
    Self = "_self",
}

type LinkProps = {
    href: string;
    children: React.ReactNode,
    target?: LinkTarget;
}

const Link = (props: LinkProps) => {
    const target = props.target ? props.target : "_blank";
    const classes = styleClasses();
    return (
        <Paragraph color={TypographyColor.LinkBlue} align={Alignment.Center}>
            <a className={classes.linkAnchor} href={props.href} target={target}>
                {props.children}
            </a>
        </Paragraph>
    )
}

export default Link;
