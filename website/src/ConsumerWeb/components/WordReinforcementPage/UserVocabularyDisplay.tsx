import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import { PrimarySwitch } from 'common/components/Switch/Switch';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Heading3 } from 'common/typography/Heading';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Form from 'common/components/Form/Form';
import {
    PrimaryButton,
    WarningButton,
    ConfirmationButton,
} from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';

import {
    UserVocabularyEntry,
    UserVocabularyType,

    UpsertUserVocabularyResponse,
    upsertUserVocabulary,
} from 'ConsumerWeb/api/user/userVocabulary';

import { WordsmithLanguageCode } from 'common/model/language/language';

const styleClasses = makeStyles({
    buttonContainer: {
        alignSelf: 'center',
    },
    button: {
        display: 'block',
        margin: 'auto',
    },
    searchTextInput: {
        width: '100%',
    },
    paddedButtonContainer: {
        padding: '5px',
    },
    submitButton: {
        margin: '5px 0',
        width: '100%',
    },
});

type UserVocabularyDisplayProps = {
    subscriptionManagementToken: string;
    userVocabularyEntries: Array<UserVocabularyEntry>;

    handleRemoveVocabularyEntry: (e: UserVocabularyEntry) => void;
}

const UserVocabularyDisplay = (props: UserVocabularyDisplayProps) => {
    const activeEntries = props.userVocabularyEntries.filter(e => e.isActive);
    const inactiveEntries = props.userVocabularyEntries.filter(e => !e.isActive);
    return (
        <Grid container>
            <Grid item xs={12}>
                {
                    activeEntries.map((e: UserVocabularyEntry) => (
                        <UserVocabularyEntryDisplay
                            key={e.uniqueHash}
                            subscriptionManagementToken={props.subscriptionManagementToken}
                            entry={e}
                            handleRemoveVocabularyEntry={props.handleRemoveVocabularyEntry} />
                    ))
                }
            </Grid>
            <Divider />
            <Grid item xs={12}>
                {
                    inactiveEntries.map((e: UserVocabularyEntry) => (
                        <UserVocabularyEntryDisplay
                            key={e.uniqueHash}
                            subscriptionManagementToken={props.subscriptionManagementToken}
                            entry={e}
                            handleRemoveVocabularyEntry={props.handleRemoveVocabularyEntry} />
                    ))
                }
            </Grid>
        </Grid>
    )
}

type UserVocabularyEntryDisplayProps = {
    subscriptionManagementToken: string;
    entry: UserVocabularyEntry;
    handleRemoveVocabularyEntry: (e: UserVocabularyEntry) => void;
}

const UserVocabularyEntryDisplay = (props: UserVocabularyEntryDisplayProps) => {
    const [ errorMessage, setErrorMessage ] = useState<string>(null);

    const [ isLoading, setIsLoading ] = useState<boolean>(false);

    const [ isActive, setIsActive ] = useState<boolean>(props.entry.isActive);

    const [ studyNote, setStudyNote ] = useState<string>(props.entry.studyNote);
    const handleStudyNoteChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setStudyNote((event.target as HTMLInputElement).value);
    };

    const handleUpdateUserVocabularyEntry = () => {
        setIsLoading(true);
        upsertUserVocabulary({
            subscriptionManagementToken: props.subscriptionManagementToken,
            languageCode: WordsmithLanguageCode.Spanish,
            displayText: props.entry.vocabularyDisplay,
            definitionId: props.entry.vocabularyId,
            entryType: props.entry.vocabularyType,
            studyNote: studyNote,
            isVisible: true,
            isActive: isActive,
        },
        (resp: UpsertUserVocabularyResponse) => {
            setIsLoading(false);
            if (!!resp.error) {
                setErrorMessage("There was an error processing your request. Try again later.");
                return;
            }
        },
        (err: Error) => {
            setIsLoading(false);
            setErrorMessage("There was an error processing your request. Try again later.");
        });
    }

    const [ shouldShowDeleteConfirmation, setShouldShowDeleteConfirmation ] = useState<boolean>(false);
    const handleDelete = () => {
        setIsLoading(true);
        upsertUserVocabulary({
            subscriptionManagementToken: props.subscriptionManagementToken,
            languageCode: WordsmithLanguageCode.Spanish,
            displayText: props.entry.vocabularyDisplay,
            definitionId: props.entry.vocabularyId,
            entryType: props.entry.vocabularyType,
            studyNote: null,
            isVisible: false,
            isActive: true,
        },
        (resp: UpsertUserVocabularyResponse) => {
            setIsLoading(false);
            if (!!resp.error) {
                setErrorMessage("There was an error processing your request. Try again later.");
                return;
            }
            props.handleRemoveVocabularyEntry(props.entry);
        },
        (err: Error) => {
            setIsLoading(false);
            setErrorMessage("There was an error processing your request. Try again later.");
        });
    }

    const classes = styleClasses();
    return (
        <Grid item xs={12}>
            <DisplayCard>
                <Form handleSubmit={handleUpdateUserVocabularyEntry}>
                    <Grid container>
                        <Grid item xs={8} md={9}>
                            <Heading3 color={TypographyColor.Primary} align={Alignment.Left}>
                                {props.entry.vocabularyDisplay}
                            </Heading3>
                        </Grid>
                        <Grid item className={classes.buttonContainer} xs={4} md={3}>
                            <PrimarySwitch
                                className={classes.button}
                                checked={isActive} onClick={() => { setIsActive(!isActive) }}
                                disabled={isLoading} />
                        </Grid>
                        <Divider />
                        {
                            !!props.entry.definition && (
                                <div>
                                    <Paragraph align={Alignment.Left}>
                                        {props.entry.definition}
                                    </Paragraph>
                                </div>
                            )
                        }
                        <Grid item xs={12}>
                            <PrimaryTextField
                                className={classes.searchTextInput}
                                id="searchTerm"
                                label="Your study note"
                                defaultValue={studyNote}
                                disabled={isLoading}
                                error={!!studyNote && studyNote.length > 250 ? "Must be 250 characters or fewer" : null}
                                variant="outlined"
                                onChange={handleStudyNoteChange} />
                        </Grid>
                        <Grid item className={classes.paddedButtonContainer} xs={12} md={4}>
                            <PrimaryButton
                                className={classes.submitButton}
                                disabled={isLoading}
                                type="submit">
                                Update
                            </PrimaryButton>
                        </Grid>
                        {
                            !shouldShowDeleteConfirmation && (
                                <Grid item className={classes.paddedButtonContainer} xs={12} md={4}>
                                    <WarningButton
                                        className={classes.submitButton}
                                        disabled={isLoading}
                                        onClick={() => {setShouldShowDeleteConfirmation(true)}}>
                                        Remove from List
                                    </WarningButton>
                                </Grid>
                            )
                        }
                        {
                            shouldShowDeleteConfirmation && (
                                <Grid item className={classes.paddedButtonContainer} xs={12} md={4}>
                                    <ConfirmationButton
                                        className={classes.submitButton}
                                        disabled={isLoading}
                                        onClick={handleDelete}>
                                        Delete from List
                                    </ConfirmationButton>
                                </Grid>
                            )
                        }
                        {
                            shouldShowDeleteConfirmation && (
                                <Grid item className={classes.paddedButtonContainer} xs={12} md={4}>
                                    <WarningButton
                                        className={classes.submitButton}
                                        onClick={() => {setShouldShowDeleteConfirmation(false)}}
                                        disabled={isLoading}>
                                        Cancel
                                    </WarningButton>
                                </Grid>
                            )
                    }
                    </Grid>
                </Form>
            </DisplayCard>
        </Grid>
    );
}

export default UserVocabularyDisplay;
