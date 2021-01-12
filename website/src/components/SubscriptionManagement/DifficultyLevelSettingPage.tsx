import React from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import Divider from '@material-ui/core/Divider';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';

import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, Color } from 'common/typography/common';

const styleClasses = makeStyles({
    displayCard: {
        padding: '10px',
    },
    contentHeaderBackArrow: {
        alignSelf: 'center',
        cursor: 'pointer',
    },
});

type ContentHeaderProps = {
    token: string;
}

const ContentHeader = (props: ContentHeaderProps) => {
    const classes = styleClasses();
    const history = useHistory();
    return (
        <Grid container>
            <Grid className={classes.contentHeaderBackArrow} onClick={() => history.push(`/manage/${props.token}`)} item xs={1}>
                <ArrowBackIcon color='action' />
            </Grid>
            <Grid item xs={11}>
                <Paragraph size={Size.Large} color={Color.Primary} align={Alignment.Left}>
                    Set your difficulty level
                </Paragraph>
            </Grid>
        </Grid>
    );
}

type Params = {
    token: string
}

type DifficultyLevelSettingPageProps = RouteComponentProps<Params>

const DifficultyLevelSettingPage = (props: DifficultyLevelSettingPageProps) => {
    const classes = styleClasses();
    const { token } = props.match.params;
    return (
        <Page>
            <Grid container>
                <Grid item xs={0} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard} variant='outlined'>
                        <ContentHeader token={token} />
                        <Divider />
                        <Paragraph size={Size.Medium} align={Alignment.Left}>
                            Select the difficulty you think is appropriate for your reading level. When you’re done, remember to enter your email on the bottom and click ‘Update’ to complete the process.
                        </Paragraph>
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

export default DifficultyLevelSettingPage;
