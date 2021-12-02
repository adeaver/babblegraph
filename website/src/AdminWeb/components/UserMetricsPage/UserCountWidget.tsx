import React, { useState, useEffect } from 'react';

import Grid from '@material-ui/core/Grid';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';

const UserCountWidget = () => {
    return (
        <DisplayCard>
            <Grid container>
                <NumberDisplay
                    value={100}
                    label="Verified Users"
                    color={TypographyColor.Confirmation} />
                <NumberDisplay
                    value={50}
                    label="Unsubscribed Users"
                    color={TypographyColor.Warning} />
                <NumberDisplay
                    value={10}
                    label="Unverified Users" />
                <NumberDisplay
                    value={10}
                    label="Blocklisted Users" />
            </Grid>
        </DisplayCard>
    );
}

type NumberDisplayProps = {
    value: number;
    label: string;
    color?: TypographyColor;
}

const NumberDisplay = (props: NumberDisplayProps) => {
    return (
        <Grid item xs={12} sm={6} md={3}>
            <Paragraph
                size={Size.Large}
                color={props.color ? props.color : TypographyColor.Gray}>
                {props.value} {props.label}
            </Paragraph>
        </Grid>
    );
}

export default UserCountWidget;
