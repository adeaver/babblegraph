import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import Divider from '@material-ui/core/Divider';

import Page from 'common/components/Page/Page';
import Color from 'common/styles/colors';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { RouteComponentProps } from 'react-router-dom';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { toTitleCase } from 'util/string/StringConvert';
import {
    GetLemmasMatchingTextResponse,
    Lemma,
    PartOfSpeech,
    Definition,
    getLemmasMatchingText
} from 'api/language/search';
import {
    AddUserLemmasForTokenResponse,
    addUserLemmasForToken
} from 'api/user/userlemma';

const styleClasses = makeStyles({
    searchCard: {
        padding: '25px',
    },
    searchHeaderDivider: {
        marginBottom: '10px',
    },
    searchTextInput: {
        width: '100%',
    },
    buttonContainer: {
        alignSelf: 'center',
    },
    button: {
        marginLeft: '15px',
    },
    lemmaDisplayRoot: {
        padding: '15px',
        borderStyle: 'solid',
        borderWidth: '1px',
        borderRadius: '5px',
        borderColor: Color.BorderGray,
        margin: '10px 0',
    },
});

type Params = {
    token: string
}

type WordReinforcementPageProps = RouteComponentProps<Params>;

const WordReinforcementPage = (props: WordReinforcementPageProps) => {
    const { token } = props.match.params;

    const [ searchTerm, setSearchTerm ] = useState<string>('');
    const [ lemmas, setLemmas ] = useState<Lemma[] | null>(null);
    const [ isLoadingLemmas, setIsLoadingLemmas ] = useState<boolean>(false);
    const [ lemmaSearchError, setLemmaSearchError ] = useState<Error>(null);

    const [ currentLoadingLemmaID, setCurrentLoadingLemmaID ] = useState<string | null>(null);
    const [ didAddLemma, setDidAddLemma ] = useState<boolean>(false);
    const [ addUserLemmaError, setAddUserLemmaError ] = useState<Error>(null);

    const handleSubmit = () => {
        setIsLoadingLemmas(true);
        getLemmasMatchingText({
            languageCode: "es",
            token: token,
            text: searchTerm,
        },
        (resp: GetLemmasMatchingTextResponse) => {
            setIsLoadingLemmas(false);
            setLemmas(resp.lemmas != null ? resp.lemmas : []);
        },
        (err: Error) => {
            setIsLoadingLemmas(false);
            setLemmaSearchError(err)
        });
    }
    const handleSelectLemma = (id: string) => {
        setCurrentLoadingLemmaID(id);
        addUserLemmasForToken({
            token: token,
            lemmaId: id,
        },
        (resp: AddUserLemmasForTokenResponse) => {
            setCurrentLoadingLemmaID(null);
            setDidAddLemma(resp.didUpdate);
        },
        (err: Error) => {
            setCurrentLoadingLemmaID(null);
            setAddUserLemmaError(err);
        });
    }

    return (
        <Page>
            <Grid container>
                <Grid item xs={0} md={3}>
                    &nbsp;
                </Grid>
                <SearchBox
                    searchTerm={searchTerm}
                    lemmas={lemmas}
                    isLoadingLemmas={isLoadingLemmas}
                    lemmaSearchError={lemmaSearchError}
                    handleSearchTermChange={setSearchTerm}
                    handleSelectLemma={handleSelectLemma}
                    handleSubmit={handleSubmit} />
            </Grid>
        </Page>
    );
}

type SearchBoxProps = {
    searchTerm: string;
    lemmas: Lemma[] | null;
    isLoadingLemmas: boolean;
    lemmaSearchError: Error;

    handleSearchTermChange: (searchTerm: string) => void;
    handleSubmit: () => void;
    handleSelectLemma: (id: string) => void;
}

const SearchBox = (props: SearchBoxProps) => {
    const handleSearchTermChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.handleSearchTermChange((event.target as HTMLInputElement).value);
    };
    const classes = styleClasses();
    return (
        <Grid item xs={12} md={6}>
            <Card className={classes.searchCard}>
                <Heading1 align={Alignment.Left} color={TypographyColor.Primary}>
                    Search for a word to track
                </Heading1>
                <Divider className={classes.searchHeaderDivider} />
                <form>
                    <Grid container>
                        <Grid item xs={9} md={10}>
                            <PrimaryTextField
                                className={classes.searchTextInput}
                                id="searchTerm"
                                label="Search for a Word"
                                defaultValue={props.searchTerm}
                                variant="outlined"
                                onChange={handleSearchTermChange} />
                        </Grid>
                        <Grid className={classes.buttonContainer} item xs={3} md={2}>
                            <PrimaryButton className={classes.button} onClick={props.handleSubmit} disabled={!props.searchTerm}>
                                Search
                            </PrimaryButton>
                        </Grid>
                    </Grid>
                </form>
                <LemmaSearchResultsDisplay
                    lemmas={props.lemmas}
                    isLoading={props.isLoadingLemmas}
                    lemmaSearchError={props.lemmaSearchError}
                    handleSelectLemma={props.handleSelectLemma} />
            </Card>
        </Grid>
    );
}

type LemmaSearchResultsDisplayProps = {
    lemmas: Lemma[] | null;
    isLoading: boolean;
    lemmaSearchError: Error;

    handleSelectLemma: (id: string) => void;
}

const LemmaSearchResultsDisplay = (props: LemmaSearchResultsDisplayProps) => {
    if (props.isLoading) {
        return (
            <LoadingSpinner />
        );
    } else if (props.lemmaSearchError != null) {
        return (
            <Paragraph>Something went wrong. Try again!</Paragraph>
        );
    } else if (props.lemmas === null) {
        return <div />;
    } else if (props.lemmas.length === 0) {
        return (
            <Paragraph>No words found.</Paragraph>
        );
    }
    return (
        <div>
        {
            props.lemmas.map((lemma: Lemma) => (
                <LemmaDisplay key={lemma.id} lemma={lemma} handleSelectLemma={props.handleSelectLemma} />
            ))
        }
        </div>
    );
}

type LemmaDisplayProps = {
    lemma: Lemma;

    handleSelectLemma: (id: string) => void;
}

const LemmaDisplay = (props: LemmaDisplayProps) => {
    const definitionText = (props.lemma.definitions || []).map((d: Definition) => (
        !!d.extraInfo ? `${d.text} ${d.extraInfo}` : d.text
    )).join('; ');
    const handleSelect = () => {
        props.handleSelectLemma(props.lemma.id);
    }
    const classes = styleClasses();
    return (
        <Grid className={classes.lemmaDisplayRoot} container>
            <Grid item xs={10}>
                <Heading3 align={Alignment.Left} color={TypographyColor.Primary}>
                    { toTitleCase(props.lemma.text) } ({props.lemma.partOfSpeech.name.toLowerCase()})
                </Heading3>
                <Paragraph align={Alignment.Left}>
                    { !!definitionText ? definitionText : 'No definition available' }
                </Paragraph>
            </Grid>
            <Grid className={classes.buttonContainer} item xs={2}>
                <PrimaryButton className={classes.button} onClick={handleSelect}>
                    Track this word
                </PrimaryButton>
            </Grid>
        </Grid>
    )
}

export default WordReinforcementPage;
