import React, { useEffect, useState } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';

import Page from 'common/components/Page/Page';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { SurveyKey, getSurveyForKey } from 'common/data/surveys/Surveys';
import SurveyDisplay from 'common/data/surveys/components/SurveyDisplay';

const styleClasses = makeStyles({
    displayCard: {
        padding: '20px',
        marginTop: '20px',
    },
})

type Params = {
    token: string
}

type FeedbackSurveyPageProps = RouteComponentProps<Params>

const FeedbackSurveyPage = (props: FeedbackSurveyPageProps) => {
    const { token } = props.match.params;

    const survey = getSurveyForKey(SurveyKey.HighOpen1);
    const classes = styleClasses();
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard}>
                        <SurveyDisplay {...survey} surveyToken={token} />
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

export default FeedbackSurveyPage;
