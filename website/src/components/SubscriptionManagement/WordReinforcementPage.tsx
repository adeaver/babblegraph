import React, { useState, useEffect } from 'react';

import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';

import Page from 'common/components/Page/Page';
import { Heading1, Heading3 } from 'common/typography/Heading';
import Paragraph from 'common/typography/Paragraph';
import { TypographyColor } from 'common/typography/common';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { RouteComponentProps } from 'react-router-dom';
import {
    GetLemmasMatchingTextResponse,
    Lemma,
    PartOfSpeech,
    Definition,
    getLemmasMatchingText
} from 'api/language/search';

type Params = {
    token: string
}

type WordReinforcementPageProps = RouteComponentProps<Params>;

const WordReinforcementPage = (props: WordReinforcementPageProps) => {
    const { token } = props.match.params;

    const [ searchTerm, setSearchTerm ] = useState<string>('');
    const [ lemmas, setLemmas ] = useState<Lemma[] | null>(null);

    const handleSubmit = () => {
        getLemmasMatchingText({
            languageCode: "es",
            token: token,
            text: searchTerm,
        },
        (resp: GetLemmasMatchingTextResponse) => {
            setLemmas(resp.lemmas);
        },
        (err: Error) => {
            console.log(err)
        });
    }
    const handleTrackLemma = (id: string) => {
        console.log(id);
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
                    handleSearchTermChange={setSearchTerm}
                    handleTrackLemma={handleTrackLemma}
                    handleSubmit={handleSubmit} />
            </Grid>
        </Page>
    );
}

type SearchBoxProps = {
    searchTerm: string;
    lemmas: Lemma[] | null;

    handleSearchTermChange: (searchTerm: string) => void;
    handleSubmit: () => void;
    handleTrackLemma: (id: string) => void;
}

const SearchBox = (props: SearchBoxProps) => {
    const handleSearchTermChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.handleSearchTermChange((event.target as HTMLInputElement).value);
    };
    return (
        <Grid item xs={12} md={6}>
            <Card>
                <Heading1 color={TypographyColor.Primary}>
                    Search
                </Heading1>
                <form>
                    <Grid container>
                        <Grid item xs={9} md={10}>
                            <PrimaryTextField
                                id="searchTerm"
                                label="Search for a Word"
                                defaultValue={props.searchTerm}
                                variant="outlined"
                                onChange={handleSearchTermChange} />
                        </Grid>
                        <Grid item xs={3} md={2}>
                            <PrimaryButton onClick={props.handleSubmit} disabled={!props.searchTerm}>
                                Search
                            </PrimaryButton>
                        </Grid>
                    </Grid>
                </form>
                {
                    props.lemmas != null && (
                        <div>
                            {
                                props.lemmas.map((lemma: Lemma) => (
                                    <LemmaDisplay key={lemma.id} lemma={lemma} handleTrack={props.handleTrackLemma} />
                                ))
                            }
                        </div>
                    )
                }
            </Card>
        </Grid>
    );
}

type LemmaDisplayProps = {
    lemma: Lemma;

    handleTrack: (id: string) => void;
}

const LemmaDisplay = (props: LemmaDisplayProps) => {
    const definitionText = (props.lemma.definitions || []).map((d: Definition) => (
        !!d.extraInfo ? `${d.text} ${d.extraInfo}` : d.text
    )).join('; ');
    const handleTrack = () => {
        props.handleTrack(props.lemma.id);
    }
    return (
        <div>
            <Heading3 color={TypographyColor.Primary}>
                {props.lemma.text} ({props.lemma.partOfSpeech.name})
            </Heading3>
            <Paragraph>
                { !!definitionText ? definitionText : 'No definition available' }
            </Paragraph>
            <PrimaryButton onClick={handleTrack}>
                Track this word
            </PrimaryButton>
        </div>
    )
}

export default WordReinforcementPage;
