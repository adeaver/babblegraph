import React from 'react';
import classnames from 'classnames';

import { makeStyles } from '@material-ui/core/styles';
import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';

import { Alignment } from 'common/typography/common';
import Color from 'common/styles/colors';
import Paragraph, { Size } from 'common/typography/Paragraph';

const styleClasses = makeStyles({
    loadingSpinner: {
        color: Color.Primary,
        display: 'block',
        margin: 'auto',
    },
});

type LoadingSpinnerProps = {
    className?: string;
}

const LoadingSpinner = (props: LoadingSpinnerProps) => {
    const classes = styleClasses();
    return (
        <Grid className={props.className} container>
            <Grid item xs={false} md={3}>
                &nbsp;
            </Grid>
            <Grid item xs={12} md={6}>
                <CircularProgress className={classes.loadingSpinner} />
                <Paragraph size={Size.Medium} align={Alignment.Center}>
                    Loading, please wait.
                </Paragraph>
            </Grid>
        </Grid>
    )
}

export default LoadingSpinner;
