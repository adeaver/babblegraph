import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';
import { PrimarySwitch } from 'common/components/Switch/Switch';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import Form from 'common/components/Form/Form';
import Autocomplete from '@material-ui/lab/Autocomplete';

import {
    WordsmithLanguageCode,
    getEnglishNameForLanguageCode,
} from 'common/model/language/language';

import {
    getTopicByID,
    GetTopicByIDResponse,
    updateIsContentTopicActive,
    UpdateIsContentTopicActiveResponse,

    TopicDisplayName,
    getAllTopicDisplayNamesForTopic,
    GetAllTopicDisplayNamesForTopicResponse,
    addTopicDisplayNameForTopic,
    AddTopicDisplayNameForTopicResponse,
    toggleTopicDisplayNameIsActive,
    ToggleTopicDisplayNameIsActiveResponse,
    updateTopicDisplayNameLabel,
    UpdateTopicDisplayNameLabelResponse,
} from 'AdminWeb/api/content/topic';

const styleClasses = makeStyles({
    headerContainer: {
        display: 'flex',
        alignItems: 'center',
    },
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
    languageSelectorContainer: {
        margin: '10px 0',
    },
    displayNameHeader: {
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
                <Grid item xs={12}>
                    <TopicDisplayNameList topicId={props.topic.id} />
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

type TopicDisplayNameListOwnProps = {
    topicId: string;
}

const TopicDisplayNameList = asBaseComponent<GetAllTopicDisplayNamesForTopicResponse, TopicDisplayNameListOwnProps>(
    (props: GetAllTopicDisplayNamesForTopicResponse & TopicDisplayNameListOwnProps & BaseComponentProps) => {
        const [ topicDisplayNames, setTopicDisplayNames ] = useState<Array<TopicDisplayName>>([]);

        const [ newTopicDisplayNameLabel, setNewTopicDisplayNameLabel ] = useState<string>(null);
        const handleLabelChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setNewTopicDisplayNameLabel((event.target as HTMLInputElement).value);
        };

        const [ languageCode, setLanguageCode ] = useState<WordsmithLanguageCode>(WordsmithLanguageCode.Spanish);
        const handleLanguageCodeUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedLanguageCode: WordsmithLanguageCode) => {
            setLanguageCode(selectedLanguageCode);
        }

        const handleSubmit = () => {
            props.setIsLoading(true);
            addTopicDisplayNameForTopic({
                topicId: props.topicId,
                languageCode: languageCode,
                label: newTopicDisplayNameLabel,
            },
            (resp: AddTopicDisplayNameForTopicResponse) => {
                setTopicDisplayNames(topicDisplayNames.concat({
                    id: resp.topicDisplayNameId,
                    topicId: props.topicId,
                    languageCode: languageCode,
                    label: newTopicDisplayNameLabel,
                    isActive: false,
                }));
                setNewTopicDisplayNameLabel(null);
                props.setIsLoading(false);
            },
            (err: Error) => {
                props.setIsLoading(false);
                props.setError(err);
            });
        }

        const displayNames = (props.topicDisplayNames || []).concat(topicDisplayNames);

        const classes = styleClasses();
        return (
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <Heading3
                            align={Alignment.Left}
                            color={TypographyColor.Primary}>
                            Add a display name
                        </Heading3>
                        <Form className={classes.confirmationForm} handleSubmit={handleSubmit}>
                            <Grid container>
                                <Grid className={classes.languageSelectorContainer} item xs={12}>
                                    <Autocomplete
                                        id="language-selector"
                                        onChange={handleLanguageCodeUpdate}
                                        options={Object.values(WordsmithLanguageCode)}
                                        getOptionLabel={(option: WordsmithLanguageCode) => getEnglishNameForLanguageCode(option)}
                                        getOptionSelected={(option: WordsmithLanguageCode) => option === languageCode}
                                        renderInput={(params) => <PrimaryTextField label="Select Language Code" {...params} />} />
                                    </Grid>
                                <Grid item xs={9} md={10}>
                                    <PrimaryTextField
                                        id="label"
                                        className={classes.labelField}
                                        label="Label"
                                        variant="outlined"
                                        defaultValue={newTopicDisplayNameLabel}
                                        onChange={handleLabelChange} />
                                </Grid>
                                <Grid item xs={3} md={2} className={classes.submitButtonContainer}>
                                    <PrimaryButton disabled={!newTopicDisplayNameLabel} type="submit">
                                        Submit
                                    </PrimaryButton>
                                </Grid>
                            </Grid>
                        </Form>
                    </DisplayCard>
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                {
                    displayNames.map((t: TopicDisplayName, idx: number) => (
                        <TopicDisplayNameView
                            key={`topic-display-name-${idx}`}
                            displayName={t}
                            setIsLoading={props.setIsLoading}
                            setError={props.setError} />
                    ))
                }
            </Grid>
        );
    },
    (
        ownProps: TopicDisplayNameListOwnProps,
        onSuccess: (resp: GetAllTopicDisplayNamesForTopicResponse) => void,
        onError: (err: Error) => void,
    ) => {
        getAllTopicDisplayNamesForTopic({
            topicId: ownProps.topicId,
        },
        onSuccess,
        onError)
    },
    false
);

type TopicDisplayNameViewProps = {
    displayName: TopicDisplayName;

    setIsLoading: (isLoading: boolean) => void;
    setError: (err: Error) => void;
}

const TopicDisplayNameView = (props: TopicDisplayNameViewProps) => {
    const [ isActive, setIsActive ] = useState<boolean>(props.displayName.isActive);

    const [ label, setLabel ] = useState<string>(props.displayName.label);
    const handleLabelChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setLabel((event.target as HTMLInputElement).value);
    };

    const handleToggleTopicDisplayName = () => {
        props.setIsLoading(true);
        toggleTopicDisplayNameIsActive({
            topicDisplayNameId: props.displayName.id,
            isActive: !isActive,
        },
        (resp: ToggleTopicDisplayNameIsActiveResponse) => {
            props.setIsLoading(false);
            setIsActive(!isActive);
        },
        (err: Error) => {
            props.setIsLoading(false);
            props.setError(err);
        });
    }
    const handleSubmitLabelUpdate = () => {
        props.setIsLoading(true);
        updateTopicDisplayNameLabel({
            topicDisplayNameId: props.displayName.id,
            label: label,
        },
        (resp: UpdateTopicDisplayNameLabelResponse) => {
            props.setIsLoading(false);
        },
        (err: Error) => {
            props.setIsLoading(false);
            props.setError(err);
        });
    }

    const classes = styleClasses();
    return (
        <Grid item xs={12} md={4}>
            <DisplayCard>
                <Grid className={classes.displayNameHeader} container>
                    <Grid item xs={8}>
                        <Heading3
                            align={Alignment.Left}
                            color={isActive ? TypographyColor.Primary : TypographyColor.Gray}>
                            {getEnglishNameForLanguageCode(props.displayName.languageCode as WordsmithLanguageCode)}
                        </Heading3>
                    </Grid>
                    <Grid item xs={4}>
                        <PrimarySwitch checked={isActive} onClick={handleToggleTopicDisplayName} />
                    </Grid>
                </Grid>
                <Form className={classes.confirmationForm} handleSubmit={handleSubmitLabelUpdate}>
                    <Grid container>
                        <Grid item xs={9} md={10}>
                            <PrimaryTextField
                                id={`${props.displayName.id}-label`}
                                className={classes.labelField}
                                label="Label"
                                variant="outlined"
                                defaultValue={label}
                                onChange={handleLabelChange} />
                        </Grid>
                        <Grid item xs={3} md={2} className={classes.submitButtonContainer}>
                            <PrimaryButton disabled={!label} type="submit">
                                Update
                            </PrimaryButton>
                        </Grid>
                    </Grid>
                </Form>
            </DisplayCard>
        </Grid>
    );
}

export default TopicManagementPage;
