import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import Grid from '@material-ui/core/Grid';

import Page from 'common/components/Page/Page';
import { Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import { PhotoKey } from 'common/data/photos/Photos';

const styleClasses = makeStyles({
    infoCardGridItem: {
        padding: '20px',
        boxSizing: 'border-box',
    },
    infoCardCardContainer: {
        padding: '20px'
    },
});

const AboutPage = () => {
    return (
        <Page withBackground={PhotoKey.Seville}>
            <Grid container>
                <InfoCard
                    title="Receive an email every day from a trusted Spanish-language news source"
                    body="Keep up with your Spanish by practicing every day with articles delivered straight to your inbox. Babblegraph only uses real-world, trusted news sources from Spain and Latin America." />
                <InfoCard
                    title="Select topics that you’re interested in to keep your articles fun and engaging."
                    body="Stay engaged with your Spanish practice by reading articles that interest you. Babblegraph makes it easy by categorizing articles from our sources." />
                <InfoCard
                    title="Control the difficulty of the articles that you read."
                    body="Babblegraph rates the difficulty of articles, allowing you to select articles that are easier or more difficult. However, since Babblegraph uses real articles from real news sources, it may not be suitable for absolute beginners." />
            </Grid>
        </Page>
    );
}

type InfoCardProps = {
    title: string;
    body: string;
}

const InfoCard = (props: InfoCardProps) => {
    const classes = styleClasses();
    return (
        <Grid className={classes.infoCardGridItem}
            item xs={12} md={4}>
            <Card className={classes.infoCardCardContainer}>
                <Heading3>
                    { props.title}
                </Heading3>
                <Paragraph>
                    {props.body}
                </Paragraph>
            </Card>
        </Grid>
    );
}

export default AboutPage;
