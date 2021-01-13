import React, { useState, useEffect } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import Divider from '@material-ui/core/Divider';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';

import Page from 'common/components/Page/Page';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';
import { Alignment, TypographyColor } from 'common/typography/common';

import {
    contentTopicDisplayMappings,
    ContentTopicDisplayMapping,
} from 'api/user/contentTopics';

const styleClasses = makeStyles({
    displayCard: {
        padding: '10px',
    },
    contentHeaderBackArrow: {
        alignSelf: 'center',
        cursor: 'pointer',
    },
});

type Params = {
    token: string
}

type InterestSelectionPageProps = RouteComponentProps<Params>

const InterestSelectionPage = (props: InterestSelectionPageProps) => {
    const classes = styleClasses();
    const { token } = props.match.params;

    const [ selectedContentTopics, setSelectedContentTopics ] = useState<Object>({});

    const handleSelectContentTopicMapping = (apiValue: string[]) => {
        setSelectedContentTopics(apiValue.reduce((accumulator: Object, next: string) => {
            return {
                ...accumulator,
                [next]: !accumulator[next],
            }
        }, selectedContentTopics));
    };

    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <Card className={classes.displayCard} variant='outlined'>
                        <ContentHeader token={token} />
                        <Divider />
                        <Paragraph size={Size.Medium} align={Alignment.Left}>
                            Click on the topics that interest you to receive emails with content on that topic. You can select as many as you’d like. When you’re done, enter your email at the bottom and click ‘Update’ to complete the process. Not every email will contain content with all the topics you’ve picked.
                        </Paragraph>
                        <ContentTopicSelectionForm
                            contentTopicDisplayMappings={contentTopicDisplayMappings.map((m: ContentTopicDisplayMapping) => ({
                                ...m,
                                isChecked: m.apiValue.reduce((val: boolean, next: string) => val || !!selectedContentTopics[next], false),
                            }))}
                            handleSelectContentTopicMapping={handleSelectContentTopicMapping} />
                    </Card>
                </Grid>
            </Grid>
        </Page>
    );
}

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
                <Paragraph size={Size.Large} color={TypographyColor.Primary} align={Alignment.Left}>
                    Manage Your Interests
                </Paragraph>
            </Grid>
        </Grid>
    );
}

type ContentTopicDisplayMappingWithChecked = {
    isChecked: boolean;
} & ContentTopicDisplayMapping;

type ContentTopicSelectionFormProps = {
    contentTopicDisplayMappings: Array<ContentTopicDisplayMappingWithChecked>;
    handleSelectContentTopicMapping: (v: string[]) => void;
}

const ContentTopicSelectionForm = (props: ContentTopicSelectionFormProps) => {
    return (
        <FormGroup row>
            <Grid container>
                {
                    props.contentTopicDisplayMappings.map((mapping: ContentTopicDisplayMappingWithChecked, idx: number) => (
                        <Grid item key={`contentTopicGridItem-${idx}`} xs={6} md={4}>
                            <FormControlLabel
                                control={
                                    <PrimaryCheckbox
                                        checked={mapping.isChecked}
                                        onChange={() => { props.handleSelectContentTopicMapping(mapping.apiValue) }}
                                        name={`checkbox-${mapping.displayText}`} />
                                }
                                label={mapping.displayText} />
                        </Grid>
                    ))
                }
            </Grid>
        </FormGroup>
    );
}

export default InterestSelectionPage;
