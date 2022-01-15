import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { asBaseComponent, BaseComponentProps } from 'AdminWeb/common/Base/BaseComponent';
import { PrimarySwitch } from 'common/components/Switch/Switch';

import {
    getTopicByID,
    GetTopicByIDResponse,
    updateIsContentTopicActive,
    UpdateIsContentTopicActiveResponse,
} from 'AdminWeb/api/content/topic';

const styleClasses = makeStyles({
    headerContainer: {
        display: 'flex',
        alignItems: 'center',
    },
});

type Params = {
    id: string;
}

type TopicManagementPageOwnProps = RouteComponentProps<Params>;

const TopicManagementPage = asBaseComponent<GetTopicByIDResponse, TopicManagementPageOwnProps>(
    (props: TopicManagementPageOwnProps & GetTopicByIDResponse & BaseComponentProps) => {
        if (!props.topic) {
            return <div />;
        }

        const [ isActive, setIsActive ] = useState<boolean>(props.topic && props.topic.isActive);

        const handleToggleTopic = () => {
            props.setIsLoading(true);
            updateIsContentTopicActive({
                id: props.topic.id,
                isActive: !isActive,
            },
            (resp: UpdateIsContentTopicActiveResponse) => {
                props.setIsLoading(false);
                setIsActive(!isActive);
            },
            (err: Error) => {
                props.setIsLoading(false);
                props.setError(err);
            });
        }

        const classes = styleClasses();
        return (
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <Grid className={classes.headerContainer} container>
                            <Grid item xs={8} md={10}>
                                <Heading1
                                    align={Alignment.Left}
                                    color={isActive ? TypographyColor.Primary : TypographyColor.Gray}>
                                    { props.topic.label }
                                </Heading1>
                            </Grid>
                            <Grid item xs={4} md={2}>
                                <PrimarySwitch checked={isActive} onClick={handleToggleTopic} />
                            </Grid>
                        </Grid>
                    </DisplayCard>
                </Grid>
            </Grid>
        );
    },
    (
        ownProps: TopicManagementPageOwnProps,
        onSuccess: (resp: GetTopicByIDResponse) => void,
        onError: (err: Error) => void,
    ) => {
        getTopicByID({
            id: ownProps.match.params.id,
        },
        onSuccess,
        onError)
    },
    true
);

export default TopicManagementPage;
