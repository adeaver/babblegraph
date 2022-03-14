import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Page from 'common/components/Page/Page';
import { TypographyColor } from 'common/typography/common';
import { Heading1, Heading2 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import ActionCard from 'common/components/ActionCard/ActionCard';
import { setLocation } from 'util/window/Location';

const styleClasses = makeStyles({
    navigationCard: {
        margin: '15px',
    }
});

const AdvertisingDashboard = () => {
    const classes = styleClasses();
    return (
        <Page>
            <Heading1 color={TypographyColor.Primary}>
                Advertising Manager
            </Heading1>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <ActionCard className={classes.navigationCard} onClick={() => setLocation("/ops/advertising-manager/vendors")}>
                        <Heading2 color={TypographyColor.Primary}>
                            Vendors
                        </Heading2>
                        <Paragraph>
                            Add and edit vendors
                        </Paragraph>
                    </ActionCard>
                    <ActionCard className={classes.navigationCard} onClick={() => setLocation("/ops/advertising-manager/sources")}>
                        <Heading2 color={TypographyColor.Primary}>
                            Sources
                        </Heading2>
                        <Paragraph>
                            Add and edit sources (i.e. affiliate programs, sponorships, etc.)
                        </Paragraph>
                    </ActionCard>
                    <ActionCard className={classes.navigationCard} onClick={() => setLocation("/ops/advertising-manager/campaigns")}>
                        <Heading2 color={TypographyColor.Primary}>
                            Campaigns
                        </Heading2>
                        <Paragraph>
                            Add or edit campaigns
                        </Paragraph>
                    </ActionCard>
                    <ActionCard className={classes.navigationCard} onClick={() => setLocation("/ops/advertising-manager/metrics")}>
                        <Heading2 color={TypographyColor.Primary}>
                            Metrics
                        </Heading2>
                        <Paragraph>
                            View metrics
                        </Paragraph>
                    </ActionCard>
                </Grid>
            </Grid>
        </Page>
    );
}

export default AdvertisingDashboard;
