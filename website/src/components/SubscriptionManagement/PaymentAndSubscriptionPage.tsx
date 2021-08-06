import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { Heading1 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';

import { ContentHeader } from './common';

type Params = {
    token: string;
}

type PaymentAndSubscriptionPageProps = RouteComponentProps<Params>

const PaymentAndSubscriptionPage = (props: PaymentAndSubscriptionPageProps) => {
    const { token } = props.match.params;

    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <ContentHeader
                            title="Subscription and Payment Settings"
                            token={token} />
                    </DisplayCard>
                </Grid>
            </Grid>
        </Page>
    );
}

export default PaymentAndSubscriptionPage;
