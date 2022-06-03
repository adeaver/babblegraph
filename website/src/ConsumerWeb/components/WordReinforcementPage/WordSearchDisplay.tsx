import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

import { withCaptchaToken } from 'common/util/grecaptcha/grecaptcha';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph from 'common/typography/Paragraph';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Alert from 'common/components/Alert/Alert';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { PrimaryButton } from 'common/components/Button/Button';

import {
    IDType,
    LanguageLookupID,
    SearchResult,
    SearchTextResult,
    SearchTextResponse,
    searchText
} from 'ConsumerWeb/api/language/search';
import {
    UserVocabularyEntry,
    UserVocabularyType,

    UpsertUserVocabularyResponse,
    upsertUserVocabulary,
} from 'ConsumerWeb/api/user/userVocabulary';

import { WordsmithLanguageCode } from 'common/model/language/language';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

const styleClasses = makeStyles({
    container: {
        margin: '10px 0',
    },
    searchTextInput: {
        width: '100%',
    },
    submitButton: {
        margin: '5px 0',
        width: '100%',
    },
});

type WordSearchDisplayOwnProps = {
    searchTerms: string[];
    wordReinforcementToken: string;
    subscriptionManagementToken: string;
    userVocabularyEntries?: Array<UserVocabularyEntry>;

    handleAddNewUserVocabularyEntry: (newEntry: UserVocabularyEntry) => void;
}

type WordSearchDisplayAPIProps = SearchTextResponse;

const WordSearchDisplay = asBaseComponent(
    (props: BaseComponentProps & WordSearchDisplayOwnProps & WordSearchDisplayAPIProps) => {
        const uniqueHashes = (props.userVocabularyEntries || []).reduce(
            (acc: { [hash: string]: boolean }, next: UserVocabularyEntry) => ({
                ...acc,
                [next.uniqueHash]: true,
            }),
            {}
        );

        return (
            <Grid container>
                {
                    !props.result.results ? (
                        <Heading3 color={TypographyColor.Warning}>
                            No results found
                        </Heading3>
                    ) : (
                        props.result.results.map((r: SearchResult) => {
                            const key = r.lookupId.id.length ? `${r.lookupId.idType}-${r.lookupId.id.join("-")}` : `${r.lookupId.idType}-${r.displayText.replace(/ +/g, "-")}`;
                            return (
                                <SearchResultDisplay
                                    key={key}
                                    searchResult={r}
                                    isAdded={!!uniqueHashes[r.uniqueHash]}
                                    isOnlyDefinition={props.result.results.length <= 1}
                                    subscriptionManagementToken={props.subscriptionManagementToken}
                                    handleAddNewUserVocabularyEntry={props.handleAddNewUserVocabularyEntry} />
                        )
                        })
                    )
                }
            </Grid>
        );
    },
    (
        ownProps: WordSearchDisplayOwnProps,
        onSuccess: (resp: WordSearchDisplayAPIProps) => void,
        onError: (err: Error) => void,
    ) => {
        withCaptchaToken("searchtext", (token: string) => {
            searchText({
                wordReinforcementToken: ownProps.wordReinforcementToken,
                languageCode: WordsmithLanguageCode.Spanish,
                text: ownProps.searchTerms,
                captchaToken: token,
            },
            (resp: SearchTextResponse) => onSuccess(resp),
            onError)
        });
    },
    false
);

type SearchResultDisplayProps = {
    subscriptionManagementToken: string;
    searchResult: SearchResult;
    isOnlyDefinition: boolean;
    isAdded: boolean;

    handleAddNewUserVocabularyEntry: (newEntry: UserVocabularyEntry) => void;
}

const SearchResultDisplay = (props: SearchResultDisplayProps) => {
    const [ studyNote, setStudyNote ] = useState<string>(null);
    const handleStudyNoteChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setStudyNote((event.target as HTMLInputElement).value);
    };

    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ errorMessage, setErrorMessage ] = useState<string>(null);

    const handleSubmit = () => {
        setIsLoading(true);
        if (!!studyNote && studyNote.length > 250) {
            setErrorMessage("The study note must be less than 250 characters");
            setIsLoading(false);
            return;
        }
        const {
            idType,
            id,
        } = props.searchResult.lookupId
        let vocabularyType = UserVocabularyType.Lemma;
        if (idType === IDType.Phrase) {
             vocabularyType = UserVocabularyType.Phrase;
        } else if (idType !== IDType.Lemma) {
            throw new Error("Invalid ID Type")
        }
        upsertUserVocabulary({
            subscriptionManagementToken: props.subscriptionManagementToken,
            languageCode: WordsmithLanguageCode.Spanish,
            displayText: props.searchResult.displayText,
            definitionId: id.length === 1 ? id[0] : undefined,
            entryType: vocabularyType,
            studyNote: studyNote,
            isVisible: true,
            isActive: true,
        },
        (resp: UpsertUserVocabularyResponse) => {
            setIsLoading(false);
            if (!!resp.error) {
                setErrorMessage("There was an error processing your request. Try again later.");
                return;
            }
            props.handleAddNewUserVocabularyEntry({
                id: resp.id,
                vocabularyId: id.length === 1 ? id[0] : undefined,
                vocabularyType: vocabularyType,
                vocabularyDisplay: props.searchResult.displayText,
                definition: (props.searchResult.definitions || []).join("; "),
                studyNote: studyNote,
                isActive: true,
                isVisible: true,
                uniqueHash: props.searchResult.uniqueHash,
            });
        },
        (err: Error) => {
            setIsLoading(false);
            setErrorMessage("There was an error processing your request. Try again later.");
        });
    }

    const classes = styleClasses();
    return (
        <Grid className={classes.container} item xs={12}>
            <DisplayCard>
                <Grid container>
                    <Grid item xs={12}>
                        <Heading3 align={Alignment.Left} color={TypographyColor.Primary}>
                            {props.searchResult.displayText}
                        </Heading3>
                        <Paragraph align={Alignment.Left}>
                            {
                                !!props.searchResult.definitions ? (
                                    props.searchResult.definitions.join("; ")
                                ) : (
                                    props.isOnlyDefinition ? (
                                        "We couldn’t find a definition for the phrase, but you can still add a study note and keep it on your vocabulary list."
                                    ) : (
                                        "Don’t see what you’re looking for here? You can add this as a study note and keep it on your vocabulary list."
                                    )
                                )
                            }
                        </Paragraph>
                    </Grid>
                    <Grid item xs={12}>
                        <PrimaryTextField
                            className={classes.searchTextInput}
                            id="searchTerm"
                            label="Add a study note"
                            defaultValue={studyNote}
                            disabled={isLoading || props.isAdded}
                            error={!!studyNote && studyNote.length > 250 ? "Must be 250 characters or fewer" : null}
                            variant="outlined"
                            onChange={handleStudyNoteChange} />
                    </Grid>
                    {
                        isLoading ? (
                            <Grid item xs={12}>
                                <LoadingSpinner />
                            </Grid>
                        ) : (
                            <Grid item xs={12} md={6}>
                                <PrimaryButton
                                    className={classes.submitButton}
                                    disabled={isLoading || props.isAdded}
                                    onClick={handleSubmit}>
                                    { !props.isAdded ? "Add to your vocabulary list" : "Already added" }
                                </PrimaryButton>
                            </Grid>
                        )
                    }
                </Grid>
            </DisplayCard>
            <Snackbar open={!!errorMessage} autoHideDuration={6000} onClose={() => setErrorMessage(null)}>
                <Alert severity="error">{errorMessage}</Alert>
            </Snackbar>
        </Grid>
    );
}

export default WordSearchDisplay;
