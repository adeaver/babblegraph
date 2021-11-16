import React from 'react';
import { useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';

const styleClasses = makeStyles({
    contentHeaderBackArrow: {
        alignSelf: 'center',
        cursor: 'pointer',
    }
});

type ContentHeaderProps = {
    token: string;
    title: string;
}

export const ContentHeader = (props: ContentHeaderProps) => {
    const classes = styleClasses();
    const history = useHistory();
    return (
        <Grid container>
            <Grid className={classes.contentHeaderBackArrow} onClick={() => history.push(`/manage/${props.token}`)} item xs={1}>
                <ArrowBackIcon color='action' />
            </Grid>
            <Grid item xs={11}>
                <Paragraph size={Size.Large} color={TypographyColor.Primary} align={Alignment.Left}>
                    {props.title}
                </Paragraph>
            </Grid>
        </Grid>
    );
}

