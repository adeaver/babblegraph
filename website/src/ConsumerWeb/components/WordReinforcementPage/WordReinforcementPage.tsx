import React, { useEffect, useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Divider from '@material-ui/core/Divider';
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
import Paragraph from 'common/typography/Paragraph';
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
    PartOfSpeech,
    LanguageLookupID,
    SearchResult,
    SearchTextResult,
    SearchTextResponse,
    searchText
} from 'ConsumerWeb/api/language/search';
import {
    LemmaMapping,

    AddUserLemmasForTokenResponse,
    addUserLemmasForToken,

    GetUserLemmasForTokenResponse,
    getUserLemmasForToken,

    updateUserLemmaActiveStateForToken,
    UpdateUserLemmaActiveStateForTokenResponse,

    removeUserLemmaForToken,
    RemoveUserLemmaForTokenResponse,
} from 'ConsumerWeb/api/user/userlemma';

const styleClasses = makeStyles({
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
                    <LemmaMappingDisplay
                        wordReinforcementToken={token}
                        subscriptionManagementToken={subscriptionManagementToken}
                        hasSubscription={!!props.userProfile.subscriptionLevel} />
                </DisplayCard>
            </CenteredComponent>
        );
    }
);

type LemmaMappingDisplayProps = {
    wordReinforcementToken: string;
    subscriptionManagementToken: string;

    hasSubscription: boolean;
}

const LemmaMappingDisplay = asBaseComponent(
    (props: LemmaMappingDisplayProps & GetUserLemmasForTokenResponse & BaseComponentProps) => {
        const [ lemmaMappings, setLemmaMappings ] = useState<Array<LemmaMapping>>(props.lemmaMappings);
        const handleAddNewLemma = (newLemma: LemmaMapping) => {
            setLemmaMappings(lemmaMappings.concat(newLemma));
        }

        return (
            <Grid container>
                <Grid item xs={12}>
                    <WordSearchForm
                        hasSubscription={props.hasSubscription}
                        wordReinforcementToken={props.wordReinforcementToken}
                        subscriptionManagementToken={props.subscriptionManagementToken}
                        handleAddNewLemma={handleAddNewLemma} />
                    <Divider />
                    <Heading3 color={TypographyColor.Primary}>
                        Your vocabulary list
                    </Heading3>
                </Grid>
            </Grid>
        )
    },
    (
        ownProps: LemmaMappingDisplayProps,
        onSuccess: (resp: GetUserLemmasForTokenResponse) => void,
        onError: (err: Error) => void,
    ) => getUserLemmasForToken({
            token: ownProps.wordReinforcementToken,
        }, onSuccess, onError),
    false,
);


type WordSearchFormProps = {
    wordReinforcementToken: string;
    subscriptionManagementToken: string;
    hasSubscription: boolean;

    handleAddNewLemma: (newLemma: LemmaMapping) => void;
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
                                isOnlyDefinition={searchResults.length <= 1}/>
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
                <Grid item xs={12} md={4}>
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
    searchResult: SearchResult;
    isOnlyDefinition: boolean;
}

const SearchResultDisplay = (props: SearchResultDisplayProps) => {
    const handleSubmit = () => {
        // TODO: this
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
                                        "We couldn’t find a definition for the phrase, but you can add a study note and keep it on your vocabulary list."
                                    ) : (
                                        "Don’t see what you’re looking for here? You can add this as a study note and keep it on your vocabulary list."
                                    )
                                )
                            }
                        </Paragraph>
                    </Grid>
                    <Grid item xs={12} md={4}>
                        <PrimaryButton
                            className={classes.submitButton}
                            onClick={handleSubmit}>
                            Add to your vocabulary list
                        </PrimaryButton>
                    </Grid>
                </Grid>
            </DisplayCard>
        </Grid>
    );
}

export default WordReinforcementPage;
