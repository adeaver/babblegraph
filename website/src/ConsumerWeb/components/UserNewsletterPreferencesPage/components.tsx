import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import Grid from '@material-ui/core/Grid';

import Color from 'common/styles/colors';
import { Heading2 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

const styleClasses = makeStyles({
    selectCardContainer: {
        padding: '5px',
    },
    selectCard: (isSelected: boolean) => ({
        borderColor: isSelected ? Color.Primary : Color.BorderGray,
        borderStyle: 'solid',
        borderWidth: '0.5px',
        cursor: 'pointer',
        padding: '5px',
    }),
})

type DateSelectorProps = {}

export const DateTimeSelector = (props: DateSelectorProps) => {
    const classes = styleClasses();
    return (
        <Grid container>
            <Grid item xs={false} md={2}>
                &nbsp;
            </Grid>
            {
                ["Sunday", "Monday", "Tuesday", "Wednesday"].map((day: string) => (
                    <Grid className={classes.selectCardContainer} item xs={6} md={2}>
                        <DateSelect key={day} label={day} />
                    </Grid>
                ))
            }
            <Grid item xs={false} md={2}>
                &nbsp;
            </Grid>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            {
                ["Thursday", "Friday", "Saturday"].map((day: string) => (
                    <Grid className={classes.selectCardContainer} item xs={6} md={2}>
                        <DateSelect key={day} label={day} />
                    </Grid>
                ))
            }
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
        </Grid>
    );
}

type DateSelectProps = {
    label: string;
}

const DateSelect = (props: DateSelectProps) => {
    const [ isSelected, setIsSelected ] = useState<boolean>(false);

    const classes = styleClasses(isSelected);
    return (
        <Card className={classes.selectCard} onClick={() => {setIsSelected(!isSelected)}}>
            <Paragraph align={Alignment.Center} color={isSelected ? TypographyColor.Primary : TypographyColor.Gray}>
                { props.label }
            </Paragraph>
        </Card>
    );
}
