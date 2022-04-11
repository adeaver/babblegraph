import React, { useEffect, useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Divider from '@material-ui/core/Divider';
import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import Color from 'common/styles/colors';
import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import { Heading3 } from 'common/typography/Heading';
import { PrimarySwitch } from 'common/components/Switch/Switch';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import {
    PrimaryButton,
    WarningButton,
    ConfirmationButton,
} from 'common/components/Button/Button';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';
import Link from 'common/components/Link/Link';
import Form from 'common/components/Form/Form';
import { withCaptchaToken, loadCaptchaScript } from 'common/util/grecaptcha/grecaptcha';
import { WordsmithLanguageCode } from 'common/model/language/language';

import {
    RouteEncryptionKey,
    LoginRedirectKey,
} from 'ConsumerWeb/api/routes/consts';
import {
    withUserProfileInformation,
    UserProfileComponentProps,
} from 'ConsumerWeb/base/UserProfile/withUserProfile';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';
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

    GetUserVocabularyResponse,
    getUserVocabulary,
} from 'ConsumerWeb/api/user/userVocabulary';

const styleClasses = makeStyles({
    paddedButtonContainer: {
        padding: '5px',
    },
    container: {
        margin: '10px 0',
    },
    premiumSubscriptionBox: {
        border: `solid 2px ${Color.Primary}`,
        margin: '10px 0',
        padding: '10px',
        borderRadius: '5px',
    },
    searchTextInput: {
        width: '100%',
    },
    submitButton: {
        margin: '5px 0',
        width: '100%',
    },
    buttonContainer: {
        alignSelf: 'center',
    },
    button: {
        display: 'block',
        margin: 'auto',
    },
    loadingSpinner: {
        color: Color.Primary,
        display: 'block',
        margin: 'auto',
    },
});

type Params = {
    token: string;
}

type WordReinforcementPageProps = RouteComponentProps<Params>;

const WordReinforcementPage = withUserProfileInformation<WordReinforcementPageProps>(
    RouteEncryptionKey.WordReinforcement,
    [ RouteEncryptionKey.SubscriptionManagement ],
    (ownProps: WordReinforcementPageProps) => {
        return ownProps.match.params.token;
    },
    LoginRedirectKey.Vocabulary,
    (props: WordReinforcementPageProps & UserProfileComponentProps) => {
        const { token } = props.match.params;
        const [ subscriptionManagementToken ] = props.userProfile.nextTokens;

        const [ hasLoadedCaptcha, setHasLoadedCaptcha ] = useState<boolean>(false);

        useEffect(() => {
            loadCaptchaScript();
            setHasLoadedCaptcha(true);
        }, []);

        if (!hasLoadedCaptcha) {
            return <LoadingSpinner />;
        }

        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Add vocabulary words"
                        backArrowDestination={`/manage/${subscriptionManagementToken}`} />
                    <UserVocabularyDisplay
                        wordReinforcementToken={token}
                        subscriptionManagementToken={subscriptionManagementToken}
                        hasSubscription={!!props.userProfile.subscriptionLevel} />
                </DisplayCard>
            </CenteredComponent>
        );
    }
);

type UserVocabularyDisplayProps = {
    wordReinforcementToken: string;
    subscriptionManagementToken: string;

    hasSubscription: boolean;
}

const UserVocabularyDisplay = asBaseComponent(
    (props: UserVocabularyDisplayProps & GetUserVocabularyResponse & BaseComponentProps) => {
        const [ userVocabularyEntries, setUserVocabularyEntries ] = useState<Array<UserVocabularyEntry>>(props.entries || []);
        const handleAddNewVocabularyEntry = (newEntry: UserVocabularyEntry) => {
            setUserVocabularyEntries(userVocabularyEntries.concat(newEntry));
        }
        const handleRemoveVocabularyEntry = (deletedEntry: UserVocabularyEntry) => {
            setUserVocabularyEntries(userVocabularyEntries.filter((e: UserVocabularyEntry) => e.uniqueHash !== deletedEntry.uniqueHash));
        }

        const uniqueHashes = userVocabularyEntries.map((e: UserVocabularyEntry) => e.uniqueHash);

        return (
            <Grid container>
                <Grid item xs={12}>
                    <WordSearchForm
                        uniqueHashes={uniqueHashes}
                        hasSubscription={props.hasSubscription}
                        wordReinforcementToken={props.wordReinforcementToken}
                        subscriptionManagementToken={props.subscriptionManagementToken}
                        handleAddNewUserVocabularyEntry={handleAddNewVocabularyEntry} />
                    <Divider />
                    <Heading3 color={TypographyColor.Primary}>
                        Your vocabulary list
                    </Heading3>
                    {
                        userVocabularyEntries.map((e: UserVocabularyEntry) => (
                            <UserVocabularyEntryDisplay key={e.uniqueHash}
                                subscriptionManagementToken={props.subscriptionManagementToken}
                                entry={e}
                                handleRemoveVocabularyEntry={handleRemoveVocabularyEntry} />
                        ))
                    }
                </Grid>
            </Grid>
        )
    },
    (
        ownProps: UserVocabularyDisplayProps,
        onSuccess: (resp: GetUserVocabularyResponse) => void,
        onError: (err: Error) => void,
    ) => getUserVocabulary({
            subscriptionManagementToken: ownProps.subscriptionManagementToken,
            languageCode: WordsmithLanguageCode.Spanish,
        }, onSuccess, onError),
    false,
);

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


type WordSearchFormProps = {
    uniqueHashes: Array<string>;
    wordReinforcementToken: string;
    subscriptionManagementToken: string;
    hasSubscription: boolean;

    handleAddNewUserVocabularyEntry: (newEntry: UserVocabularyEntry) => void;
}

const wordSearchErrorMessages = {
    "default": "Something went wrong processing your request",
}

const WordSearchForm = (props: WordSearchFormProps) => {
    const [ searchTerm, setSearchTerm ] = useState<string>(null);
    const handleSearchTermChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setSearchTerm((event.target as HTMLInputElement).value);
    };

    const [ isLoading, setIsLoading ] = useState<boolean>(false);
    const [ errorMessage, setErrorMessage ] = useState<string>(null);

    const [ searchResults, setSearchResults ] = useState<SearchResult[]>(null);

    const handleSubmit = () => {
        setIsLoading(true);
        const terms = searchTerm.trim().split(/ +/g);
        if (terms.length > 1 && !props.hasSubscription) {
            setIsLoading(false);
            setErrorMessage("You need to upgrade to Babblegraph Premium to lookup phrases.");
        } else if (terms.length == 0) {
            setIsLoading(false);
        } else {
            withCaptchaToken("searchtext", (token: string) => {
                searchText({
                    wordReinforcementToken: props.wordReinforcementToken,
                    languageCode: WordsmithLanguageCode.Spanish,
                    text: terms,
                    captchaToken: token,
                },
                (resp: SearchTextResponse) => {
                    setIsLoading(false);
                    if (!!resp.error) {
                        setErrorMessage(wordSearchErrorMessages["default"]);
                        setSearchResults(null);
                        return;
                    }
                    setErrorMessage(null);
                    setSearchResults(resp.result.results || []);
                },
                (err: Error) => {
                    setIsLoading(false);
                    setErrorMessage(wordSearchErrorMessages["default"]);
                });
            });
        }
    }

    let body;
    if (isLoading) {
        body = (
            <LoadingSpinner />
        );
    } else if (searchResults != null && !!searchResults.length) {
        body = (
            <Grid item xs={12}>
                <Grid container>
                {
                    searchResults.map((r: SearchResult) => {
                        const key = r.lookupId.id.length ? `${r.lookupId.idType}-${r.lookupId.id.join("-")}` : `${r.lookupId.idType}-${r.displayText.replace(/ +/g, "-")}`;
                        return (
                            <SearchResultDisplay
                                key={key}
                                searchResult={r}
                                isAdded={props.uniqueHashes.indexOf(r.uniqueHash) !== -1}
                                isOnlyDefinition={searchResults.length <= 1}
                                hasSubscription={props.hasSubscription}
                                subscriptionManagementToken={props.subscriptionManagementToken}
                                handleAddNewUserVocabularyEntry={props.handleAddNewUserVocabularyEntry} />
                        )
                    })
                }
                </Grid>
            </Grid>
        )
    } else if (searchResults != null && !searchResults.length) {
        body = (
            <Grid item xs={12}>
                <Heading3 align={Alignment.Left}>
                    No results found
                </Heading3>
            </Grid>
        )
    }

    const classes = styleClasses();
    return (
        <Form handleSubmit={handleSubmit}>
            <Grid container>
                {
                    !props.hasSubscription && (
                        <Grid className={classes.premiumSubscriptionBox} item xs={12}>
                            <Paragraph color={TypographyColor.Primary}>
                                ¡No tengas celos de los suscriptores de Babblegraph Premium!
                            </Paragraph>
                            <Paragraph>
                                With Babblegraph Premium, you can search for phrases like "tener celos" instead of just words.
                            </Paragraph>
                            <Link href={`/manage/${props.subscriptionManagementToken}/premium`}>
                                Learn more about Babblegraph Premium here.
                            </Link>
                        </Grid>
                    )
                }
                <Grid item xs={12}>
                    <PrimaryTextField
                        className={classes.searchTextInput}
                        id="searchTerm"
                        label={`Search for a word${props.hasSubscription ? " or phrase" : ""}`}
                        defaultValue={searchTerm}
                        variant="outlined"
                        onChange={handleSearchTermChange} />
                </Grid>
                <Grid item xs={12} md={6}>
                    <PrimaryButton
                        className={classes.submitButton}
                        disabled={!searchTerm || isLoading}
                        type="submit">
                        Search
                    </PrimaryButton>
                </Grid>
                {body}
            </Grid>
            <Snackbar open={!!errorMessage} autoHideDuration={6000} onClose={() => setErrorMessage(null)}>
                <Alert severity="error">{errorMessage}</Alert>
            </Snackbar>
        </Form>
    );
}

type SearchResultDisplayProps = {
    subscriptionManagementToken: string;
    searchResult: SearchResult;
    isOnlyDefinition: boolean;
    isAdded: boolean;
    hasSubscription: boolean;

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
                definition: props.searchResult.definitions.join("; "),
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
                    {
                        props.hasSubscription && (
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
                        )
                    }
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

export default WordReinforcementPage;
