import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Color from 'common/styles/colors';
import { Photographer, Source } from 'common/data/photos/Photos';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';

const styleClasses = makeStyles({
    photoCreditText: {
        backgroundColor: Color.BackgroundGray,
    },
    photoCreditLink: {
        textDecoration: 'none',
        fontWeight: '800',
        color: Color.White,
    },
    photoCredit: {
        opacity: '0.75',
    },
});

type BackgroundPhotoCreditProps = {
    photographer: Photographer;
    source: Source;
}

const BackgroundPhotoCredit = (props: BackgroundPhotoCreditProps) => {
    const classes = styleClasses();
    return (
        <Grid className={classes.photoCredit} container>
            <Grid item xs={6} md={9}>
                &nbsp;
            </Grid>
            <Grid className={classes.photoCreditText} item xs={6} md={3}>
                <Paragraph size={Size.Small} color={TypographyColor.White}>
                    Photo by <a className={classes.photoCreditLink} target="_blank" href={props.photographer.url}>{props.photographer.name}</a> on <a className={classes.photoCreditLink} target="_blank" href={props.source.url}>{props.source.name}</a>
                </Paragraph>
            </Grid>
        </Grid>
    );
}

export default BackgroundPhotoCredit;
