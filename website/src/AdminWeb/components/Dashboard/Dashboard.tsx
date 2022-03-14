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
                <NavigationCard
                    location="/ops/permission-manager"
                    title="Manage Admin Permissions"
                    description="Activate or deactivate permissions for admin users" />
                <NavigationCard
                    location="/ops/user-metrics"
                    title="View User Metrics"
                    description="View metrics such as user counts, email usage statistics, etc." />
                <NavigationCard
                    location="/ops/advertising-manager"
                    title="Advertising Manager"
                    description="View and edit affiliate and sponsor advertisements" />
                <NavigationCard
                    location="/ops/content-manager"
                    title="Content Manager"
                    description="View and edit content topics and available content sources" />
                <NavigationCard
                    location="/ops/billing-manager"
                    title="Billing Manager"
                    description="Manage premium newsletter subscriptions" />
                <NavigationCard
                    location="/ops/blog-manager"
                    title="Blog Manager"
                    description="Write, edit, delete, or promote blog posts." />
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
        <Grid item xs={12} md={4}>
            <ActionCard className={classes.navigationCard} onClick={() => setLocation(props.location)}>
                <Heading2 color={TypographyColor.Primary}>
                    { props.title }
                </Heading2>
                <Paragraph>
                    { props.description }
                </Paragraph>
            </ActionCard>
        </Grid>
    );
}

export default Dashboard;
