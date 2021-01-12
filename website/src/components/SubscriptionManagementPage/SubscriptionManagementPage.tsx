import React from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';

import Paragraph from 'common/typography/Paragraph';
import { Alignment } from 'common/typography/common';

type Params = {
    token: string
}


type SubscriptionManagementPageProps = RouteComponentProps<Params>

const SubscriptionManagementPage = (props: SubscriptionManagementPageProps) => {
    return (
        <Grid container spacing={3}>
            <Grid item xs={12} sm={12} md={6} lg={6} xl={6}>
                <Card variant='outlined'>
                    <Paragraph align={Alignment.Left}>Manage Your Interests</Paragraph>
                </Card>
            </Grid>
            <Grid item xs={12} sm={12} md={6} lg={6} xl={6}>
                <Card variant='outlined'>
                    <Paragraph align={Alignment.Left}>Set your difficulty level</Paragraph>
                </Card>
            </Grid>
            <Grid item xs={12} sm={12} md={6} lg={6} xl={6}>
                <Card variant='outlined'>
                    <Paragraph align={Alignment.Left}>Unsubscribe</Paragraph>
                </Card>
            </Grid>
        </Grid>
    );
}

export default SubscriptionManagementPage;
