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
                    <NavigationCard
                        location="/ops/permission-manager"
                        title="Manage Admin Permissions"
                        description="Activate or deactivate permissions for admin users" />
                    <NavigationCard
                        location="/ops/user-metrics"
                        title="View User Metrics"
                        description="View metrics such as user counts, email usage statistics, etc." />
                </Grid>
            </Grid>
        </Page>
    );
}

type NavigationCardProps = {
    location: string;
    title: string;
    description: string;
}

const NavigationCard = (props: NavigationCardProps) => {
    const classes = styleClasses();
    return (
        <ActionCard className={classes.navigationCard} onClick={() => setLocation(props.location)}>
            <Heading2 color={TypographyColor.Primary}>
                { props.title }
            </Heading2>
            <Paragraph>
                { props.description }
            </Paragraph>
        </ActionCard>
    );
}

export default Dashboard;
