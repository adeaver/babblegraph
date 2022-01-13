import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Page from 'common/components/Page/Page';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import ActionCard from 'common/components/ActionCard/ActionCard';
import { asBaseComponent, BaseComponentProps } from 'AdminWeb/common/Base/BaseComponent';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import Link, { LinkTarget } from 'common/components/Link/Link';

import {
    Topic,
    GetAllContentTopicsResponse,
    getAllContentTopics,

    AddTopicResponse,
    addTopic,
} from 'AdminWeb/api/content/topic';

const styleClasses = makeStyles({
    submitButtonContainer: {
        alignSelf: 'center',
        padding: '5px',
    },
    labelField: {
        width: '100%',
    },
    confirmationForm: {
        padding: '10px 0',
        width: '100%',
    },
    topicDisplayCard: {
        margin: '10px 0',
    },
});

const TopicListPage = asBaseComponent(
    (props: GetAllContentTopicsResponse & BaseComponentProps) => {
        const [ addedTopics, setAddedTopics ] = useState<Array<Topic>>([]);

        const handleAddNewTopic = (topic: Topic) => {
            setAddedTopics(addedTopics.concat(topic));
        }
        const topics = (props.topics || []).concat(addedTopics);
        return (
            <div>
                <AddTopicForm handleAddNewTopic={handleAddNewTopic} />
                <Grid container>
                    <Grid item xs={false} md={3}>
                        &nbsp;
                    </Grid>
                    <Grid item xs={12} md={6}>
                        {
                            topics.map((t: Topic, idx: number) => (
                                <TopicDisplay key={`topic-display-${idx}`} {...t} />
                            ))
                        }
                    </Grid>
                </Grid>
            </div>
        );
    },
    (
        onSuccess: (resp: GetAllContentTopicsResponse) => void,
        onError: (err: Error) => void,
    ) => getAllContentTopics({}, onSuccess, onError),
    true,
)

const TopicDisplay = (props: Topic) => {
    const classes = styleClasses();
    return (
        <DisplayCard className={classes.topicDisplayCard}>
            <Heading3 color={TypographyColor.Primary}>
                {props.label}
            </Heading3>
            <Link href={`/ops/content-manager/topics/${props.id}`} target={LinkTarget.Self}>
                Manage this topic
            </Link>
        </DisplayCard>
    );
}

type AddTopicFormProps = {
    handleAddNewTopic: (topic: Topic) => void;
}

const AddTopicForm = asBaseComponent<{}, AddTopicFormProps>(
    (props: AddTopicFormProps & BaseComponentProps) => {
        const [ label, setLabel ] = useState<string>(null);

        const handleSubmit = () => {
            props.setIsLoading(true);
            addTopic({
                label: label,
            },
            (resp: AddTopicResponse) => {
                props.setIsLoading(false);
                props.handleAddNewTopic({
                    id: resp.id,
                    label: label,
                    isActive: false,
                });
                setLabel(null);
            },
            (err: Error) => {
                props.setIsLoading(false);
                props.setError(err);
            });
        }
        const handleLabelChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setLabel((event.target as HTMLInputElement).value);
        };

        const classes = styleClasses();
        return (
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <Heading3 color={TypographyColor.Primary}>
                            Add a topic
                        </Heading3>
                        <Form className={classes.confirmationForm} handleSubmit={handleSubmit}>
                            <Grid container>
                                <Grid item xs={9} md={10}>
                                    <PrimaryTextField
                                        id="label"
                                        className={classes.labelField}
                                        label="Label"
                                        variant="outlined"
                                        defaultValue={label}
                                        onChange={handleLabelChange} />
                                </Grid>
                                <Grid item xs={3} md={2} className={classes.submitButtonContainer}>
                                    <PrimaryButton disabled={!label} type="submit">
                                        Submit
                                    </PrimaryButton>
                                </Grid>
                            </Grid>
                        </Form>
                    </DisplayCard>
                </Grid>
            </Grid>
        );
    },
    (
        onSuccess: (props: {}) => void,
        onError: (err: Error) => void,
    ) => onSuccess({}),
    false,
);

export default TopicListPage;
