import React, { useEffect, useState } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PrimaryButton } from 'common/components/Button/Button';
import Link, { LinkTarget } from 'common/components/Link/Link';
import { TypographyColor } from 'common/typography/common';
import Page from 'common/components/Page/Page';

import PremiumInformationPanel from 'components/PremiumInformation/PremiumInformationPanel';

const styleClasses = makeStyles({
    callToActionButton: {
        margin: '15px 0',
        width: '100%',
    },
});

type Params = {
    token: string;
}

type PremiumInformationPageProps = RouteComponentProps<Params>

const PremiumInformationPage = (props: PremiumInformationPageProps) => {
    const { token } = props.match.params;

    const classes = styleClasses();
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <PremiumInformationPanel />
                        <Divider />
                        <Grid container>
                            <Grid item xs={2}>
                                &nbsp;
                            </Grid>
                            <Grid item xs={8}>
                                <PrimaryButton className={classes.callToActionButton} size="large">
                                    Try Babblegraph Premium
                                </PrimaryButton>
                            </Grid>
                        </Grid>
                        <Link href={`/manage/${token}`} target={LinkTarget.Self}>
                            Go back to subscription management
                        </Link>
                    </DisplayCard>
                </Grid>
            </Grid>
        </Page>
    );
}

export default PremiumInformationPage;
