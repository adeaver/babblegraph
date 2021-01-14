import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import Grid from '@material-ui/core/Grid';

import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';

const styleClasses = makeStyles({
    displayCard: {
        padding: '10px',
    },
});

const HomePage = () => {
    const classes = styleClasses();
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard} variant='outlined'>
                        <Paragraph size={Size.Large} color={TypographyColor.Primary}>
                            Still in development
                        </Paragraph>
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

export default HomePage;
