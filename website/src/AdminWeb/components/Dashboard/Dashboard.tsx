import React from 'react';

import Grid from '@material-ui/core/Grid';

import Page from 'common/components/Page/Page';
import { TypographyColor } from 'common/typography/common';
import { Heading1, Heading2 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import ActionCard from 'common/components/ActionCard/ActionCard';
import { setLocation } from 'util/window/Location';

const Dashboard = () => {
    return (
        <Page>
            <Heading1 color={TypographyColor.Primary}>
                babblegraph
            </Heading1>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <ActionCard onClick={() => setLocation("/ops/user-metrics")}>
                        <Heading2 color={TypographyColor.Primary}>
                            View User Metrics
                        </Heading2>
                        <Paragraph>
                            View metrics such as user counts, email usage statistics, etc.
                        </Paragraph>
                    </ActionCard>
                </Grid>
            </Grid>
        </Page>
    );
}

export default Dashboard;
