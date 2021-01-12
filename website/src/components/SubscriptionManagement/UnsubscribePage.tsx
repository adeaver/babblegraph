import React from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import Divider from '@material-ui/core/Divider';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';

import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, Color } from 'common/typography/common';

const styleClasses = makeStyles({
    displayCard: {
        padding: '10px',
    },
    contentHeaderBackArrow: {
        alignSelf: 'center',
        cursor: 'pointer',
    },
});

type ContentHeaderProps = {
    token: string;
}

const ContentHeader = (props: ContentHeaderProps) => {
    const classes = styleClasses();
    const history = useHistory();
    return (
        <Grid container>
            <Grid className={classes.contentHeaderBackArrow} onClick={() => history.push(`/manage/${props.token}`)} item xs={1}>
                <ArrowBackIcon color='action' />
            </Grid>
            <Grid item xs={11}>
                <Paragraph size={Size.Large} color={Color.Primary} align={Alignment.Left}>
                    Unsubscribe
                </Paragraph>
            </Grid>
        </Grid>
    );
}

type Params = {
    token: string
}

type UnsubscribePageProps = RouteComponentProps<Params>

const UnsubscribePage = (props: UnsubscribePageProps) => {
    const classes = styleClasses();
    const { token } = props.match.params;
    return (
        <Page>
            <Grid container>
                <Grid item xs={0} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard} variant='outlined'>
                        <ContentHeader token={token} />
                        <Divider />
                        <Paragraph size={Size.Medium} align={Alignment.Left}>
                            We’re sorry to see you go! When you unsubscribe, you won’t receive any more daily emails or any other communication from us. But you can always come back by resubscribing on the homepage. To unsubscribe, just enter your email and click the unsubscribe button.
                        </Paragraph>
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

export default UnsubscribePage;
