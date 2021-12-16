import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import { setLocation } from 'util/window/Location';

const styleClasses = makeStyles({
    displayCardHeaderBackArrow: {
        alignSelf: 'center',
        cursor: 'pointer',
    }
});

type DisplayCardHeaderProps = {
    backArrowDestination?: string;
    title: string;
}

const DisplayCardHeader = (props: DisplayCardHeaderProps) => {
    const classes = styleClasses();
    return (
        <Grid container>
            {
                !!props.backArrowDestination && (
                    <Grid className={classes.displayCardHeaderBackArrow} onClick={() => { setLocation(props.backArrowDestination) }} item xs={1}>
                        <ArrowBackIcon color='action' />
                    </Grid>
                )
            }
            <Grid item xs={!!props.backArrowDestination ? 11 : 12}>
                <Paragraph size={Size.Large} color={TypographyColor.Primary} align={Alignment.Left}>
                    {props.title}
                </Paragraph>
            </Grid>
        </Grid>
    );
}

export default DisplayCardHeader;
