import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import Divider from '@material-ui/core/Divider';
import ArrowForwardIcon from '@material-ui/icons/ArrowForward';

import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, Color } from 'common/typography/common';

type Params = {
    token: string
}

const styleClasses = makeStyles({
    actionCard: {
        padding: '5px',
        height: '100%',
    },
    headerArrow: {
        alignSelf: 'center',
    },
});

type ActionCardProps = {
    title: string;
    children: string;
}

const ActionCard = (props: ActionCardProps) => {
    const classes = styleClasses();
    return (
        <Grid item xs={12} md={6}>
            <Card className={classes.actionCard} variant='outlined'>
                <Grid container>
                    <Grid item xs={11}>
                        <Paragraph size={Size.Large} color={Color.Primary} align={Alignment.Left}>{props.title}</Paragraph>
                    </Grid>
                    <Grid className={classes.headerArrow} item xs={1}>
                        <ArrowForwardIcon color="action" />
                    </Grid>
                </Grid>
                <Divider />
                <Paragraph size={Size.Medium} color={Color.Gray} align={Alignment.Left}>{props.children}</Paragraph>
            </Card>
        </Grid>
    );
}

type SubscriptionManagementPageProps = RouteComponentProps<Params>

const SubscriptionManagementPage = (props: SubscriptionManagementPageProps) => {
    const classes = styleClasses();
    return (
        <Grid container spacing={2}>
            <ActionCard title='Manage Your Interests'>
                Select some topics you’re interested in reading more about or deselect some topics you’d like to read about less. This is a great way to make sure that the content you get is fun and engaging.
            </ActionCard>
            <ActionCard title='Set your difficulty level'>
                If your daily email is too hard or too easy, you can change the difficulty level here.
            </ActionCard>
            <ActionCard title='Unsubscribe'>
                If you’re no longer interested in receiving daily emails, you can unsubscribe here. By unsubscribing, we won’t send you any more emails about anything.
            </ActionCard>
        </Grid>
    );
}

export default SubscriptionManagementPage;
