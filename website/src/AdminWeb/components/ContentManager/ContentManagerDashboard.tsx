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

const ContentManagerDashboard = () => {
    const classes = styleClasses();
    return (
        <Page>
            <Heading1 color={TypographyColor.Primary}>
                Content Manager
            </Heading1>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <ActionCard className={classes.navigationCard} onClick={() => setLocation("/ops/content-manager/topics")}>
                        <Heading2 color={TypographyColor.Primary}>
                            Edit Topics
                        </Heading2>
                        <Paragraph>
                            Edit available content topics
                        </Paragraph>
                    </ActionCard>
                    <ActionCard className={classes.navigationCard} onClick={() => setLocation("/ops/content-manager/sources")}>
                        <Heading2 color={TypographyColor.Primary}>
                            Edit sources
                        </Heading2>
                        <Paragraph>
                            Add or deactivate content sources
                        </Paragraph>
                    </ActionCard>
                    <ActionCard className={classes.navigationCard} onClick={() => setLocation("/ops/content-manager/podcasts")}>
                        <Heading2 color={TypographyColor.Primary}>
                            Search Podcasts
                        </Heading2>
                        <Paragraph>
                            Search and add podcasts
                        </Paragraph>
                    </ActionCard>
                </Grid>
            </Grid>
        </Page>
    );
}

export default ContentManagerDashboard;
