import './Paragraph.scss';
import './common.scss';

import React from 'react';
import classNames from 'classnames';

import { asTypography, TypographyProps } from './common';

export enum Size {
    Small,
    Medium,
    Large,
    ExtraLarge,
}

type ParagraphProps = {
    size?: Size;
} & TypographyProps;

const Paragraph = asTypography((props: ParagraphProps) => {
    const className = classNames('Paragraph__root', props.className, {
        'Paragraph__size-small': props.size === Size.Small,
        'Paragraph__size-medium': props.size == null || props.size === Size.Medium,
        'Paragraph__size-large': props.size === Size.Large,
        'Paragraph__size-extra-large': props.size === Size.ExtraLarge,
    });
    return (
        <p className={className}>{props.children}</p>
    );
});

export default Paragraph;
