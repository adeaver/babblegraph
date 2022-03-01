import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { TypographyColor } from 'common/typography/common';
import { Heading3 } from 'common/typography/Heading';
import Paragraph, { Size } from 'common/typography/Paragraph';
import Link from 'common/components/Link/Link';

import Form from 'common/components/Form/Form';
import Autocomplete from '@material-ui/lab/Autocomplete';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';

import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';

import {
    SupportedGenre,
    SupportedRegion,
    SearchOptions,
    PodcastMetadata,

    GetPodcastSearchOptionsResponse,
    getPodcastSearchOptions,

    SearchPodcastsResponse,
    searchPodcasts,
} from 'AdminWeb/api/podcasts/podcasts';

const styleClasses = makeStyles({
    formComponent: {
        padding: '10px',
    },
});

const PodcastSearchPage = asBaseComponent<GetPodcastSearchOptionsResponse, {}>(
    (props: GetPodcastSearchOptionsResponse & BaseComponentProps) => {
        const [ isLoading, setIsLoading ] = useState<boolean>(false);

        const [ language, setLanguage ] = useState<string>(null);
        const handleLanguageUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedLanguage: string) => {
            setLanguage(selectedLanguage);
        }

        const [ region, setRegion ] = useState<string>(null);
        const handleRegionUpdate =  (_: React.ChangeEvent<HTMLSelectElement>, selectedRegion: SupportedRegion) => {
            setRegion(selectedRegion.apiValue);
        }

        const [ genre, setGenre ] = useState<number>(null);
        const handleGenreUpdate =  (_: React.ChangeEvent<HTMLSelectElement>, selectedGenre: SupportedGenre) => {
            setGenre(selectedGenre.apiValue);
        }

        const [ podcasts, setPodcasts ] = useState<Array<PodcastMetadata>>(null);
        const [ pageNumber, setPageNumber ] = useState<number>(undefined);

        const handleSubmit = () => {
            setIsLoading(true);
            searchPodcasts({
                params: {
                    genre: genre,
                    region: region,
                    language: language,
                    pageNumber: pageNumber,
                },
            },
            (resp: SearchPodcastsResponse) => {
                setIsLoading(false);
                setPodcasts(resp.podcasts);
                setPageNumber(resp.nextPageNumber);
            },
            (err: Error) => {
                setIsLoading(false);
                props.setError(err);
            });
        }

        const classes = styleClasses();
        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Podcast Search"
                        backArrowDestination="/ops/content" />
                    <Form handleSubmit={handleSubmit}>
                        <Grid container>
                            <Grid className={classes.formComponent} item xs={12} md={6}>
                                <Autocomplete
                                    id="language-selector"
                                    onChange={handleLanguageUpdate}
                                    options={props.options.supportedLanguages}
                                    getOptionSelected={(option: string) => option === language}
                                    renderInput={(params) => <PrimaryTextField label="Select Language" {...params} />} />
                            </Grid>
                            <Grid className={classes.formComponent} item xs={12} md={6}>
                                <Autocomplete
                                    id="region-selector"
                                    onChange={handleRegionUpdate}
                                    options={props.options.supportedRegions}
                                    getOptionLabel={(option: SupportedRegion) => option.displayName}
                                    getOptionSelected={(option: SupportedRegion) => option.apiValue === region}
                                    renderInput={(params) => <PrimaryTextField label="Select Region" {...params} />} />
                            </Grid>
                            <Grid className={classes.formComponent} item xs={12} md={6}>
                                <Autocomplete
                                    id="genre-selector"
                                    onChange={handleGenreUpdate}
                                    options={props.options.genres}
                                    getOptionLabel={(option: SupportedGenre) => option.displayName}
                                    getOptionSelected={(option: SupportedGenre) => option.apiValue === genre}
                                    renderInput={(params) => <PrimaryTextField label="Select Genre" {...params} />} />
                            </Grid>
                            <Grid className={classes.formComponent} item xs={12}>
                                <PrimaryButton
                                    type="submit"
                                    disabled={!genre || !region || !language}>
                                    Submit
                                </PrimaryButton>
                            </Grid>
                        </Grid>
                    </Form>
                </DisplayCard>
                {
                    isLoading && <LoadingSpinner />
                }
                {
                    (!!podcasts && !!podcasts.length) && (
                        <PodcastResultsDisplay
                            podcasts={podcasts}
                            hasNextPage={pageNumber != null}
                            handleNextPage={handleSubmit} />
                    )
                }
            </CenteredComponent>
        );
    },
    (
        ownProps: {},
        onSuccess: (resp: GetPodcastSearchOptionsResponse) => void,
        onError: (err: Error) => void,
    ) => getPodcastSearchOptions(ownProps, onSuccess, onError),
    true,
);

type PodcastResultsDisplayProps = {
    hasNextPage: boolean;
    podcasts: Array<PodcastMetadata>;

    handleNextPage: () => void;
}

const PodcastResultsDisplay = (props: PodcastResultsDisplayProps) => {
    const classes = styleClasses();
    return (
        <DisplayCard>
            <Grid container>
                <Grid className={classes.formComponent} item xs={12}>
                    <PrimaryButton
                        type="submit"
                        disabled={props.hasNextPage}
                        onClick={props.handleNextPage}>
                        Next Page
                    </PrimaryButton>
                </Grid>
                {
                    props.podcasts.map((p: PodcastMetadata, idx: number) => (
                        <PodcastDisplay key={`podcast-${idx}`} {...p} />
                    ))
                }
            </Grid>
        </DisplayCard>
    )
}

const PodcastDisplay = (props: PodcastMetadata) => {
    const classes = styleClasses();
    return (
        <Grid className={classes.podcastDisplayRoot} item xs={12}>
            <Heading3 color={TypographyColor.Primary}>
                {props.title}
            </Heading3>
            <Paragraph>
                {props.description}
            </Paragraph>
            <Link href={props.website}>
                View wesbite
            </Link>
            <Link href={props.listenNotesUrl}>
                View on third party
            </Link>
            <Paragraph size={Size.Small}>
                Type: {props.type}, {props.totalNumberOfEpisodes} episodes
            </Paragraph>
            <Paragraph size={Size.Small}>
                Country: {props.country}, in {props.language}
            </Paragraph>
            <Divider />
        </Grid>
    );
}

export default PodcastSearchPage;
