import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';
import ClearIcon from '@material-ui/icons/Clear';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import Color from 'common/styles/colors';
import { Alignment, TypographyColor } from 'common/typography/common';
import { Heading3 } from 'common/typography/Heading';
import Paragraph, { Size } from 'common/typography/Paragraph';
import Link from 'common/components/Link/Link';

import Form from 'common/components/Form/Form';
import Autocomplete from '@material-ui/lab/Autocomplete';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';

import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';

import { WordsmithLanguageCode, getEnglishNameForLanguageCode } from 'common/model/language/language';
import { CountryCode, getEnglishNameForCountryCode } from 'common/model/geo/geo';
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
import {
    Topic,

    getAllContentTopics,
    GetAllContentTopicsResponse,
} from 'AdminWeb/api/content/topic';

const styleClasses = makeStyles({
    formComponent: {
        padding: '10px',
    },
    addPodcastFormComponent: {
        width: '100%',
    },
    addPodcastForm: {
        display: 'flex',
        justifyContent: 'center',
    },
    alignedContainer: {
        display: 'flex',
        alignItems: 'center',
    },
    removeContentTopicIcon: {
        color: Color.Warning,
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
                        backArrowDestination="/ops/content-manager" />
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
    const [ rssFeedURL, setRSSFeedURL ] = useState<string>(null);
    const handleRSSFeedURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRSSFeedURL((event.target as HTMLInputElement).value);
    }

    const [ shouldShowCaptureForm, setShouldShowCaptureForm ] = useState<boolean>(false);


    const handlePreparePodcast = () => {
        if (!rssFeedURL) {
            return;
        }
        setShouldShowCaptureForm(true);
    }

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
            <Form handleSubmit={handlePreparePodcast}>
                <Grid className={classes.addPodcastForm} container>
                    <Grid className={classes.formComponent} item xs={8}>
                        <PrimaryTextField
                            id="rss-feed-url"
                            label="RSS Feed URL"
                            variant="outlined"
                            className={classes.addPodcastFormComponent}
                            defaultValue={rssFeedURL}
                            disabled={shouldShowCaptureForm}
                            onChange={handleRSSFeedURLChange} />
                    </Grid>
                    <Grid className={classes.formComponent} item xs={4}>
                        <PrimaryButton className={classes.addPodcastFormComponent} type="submit" disabled={shouldShowCaptureForm}>
                            Add Podcast
                        </PrimaryButton>
                    </Grid>
                </Grid>
            </Form>
            {
                shouldShowCaptureForm && (
                    <PodcastCaptureForm
                        rssFeedURL={rssFeedURL}
                        website={props.website}
                        title={props.title} />
                )
            }
            <Divider />
        </Grid>
    );
}

type PodcastCaptureFormOwnProps = {
    rssFeedURL: string;
    website: string;
    title: string;
}

const PodcastCaptureForm = asBaseComponent<GetAllContentTopicsResponse, PodcastCaptureFormOwnProps>(
    (props: GetAllContentTopicsResponse & PodcastCaptureFormOwnProps & BaseComponentProps) => {
        const [ languageCode, setLanguageCode ] = useState<WordsmithLanguageCode>(null);
        const handleUpdateLanguageCode = (_: React.ChangeEvent<HTMLSelectElement>, selectedLanguageCode: WordsmithLanguageCode) => {
            setLanguageCode(selectedLanguageCode);
        }

        const [ countryCode, setCountryCode ] = useState<CountryCode>(null);
        const handleCountyCodeUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedCountryCode: CountryCode) => {
            setCountryCode(selectedCountryCode);
        }

        const [ activeTopicMappings, setActiveTopicMappings ] = useState<Array<Topic>>([]);
        const handleTopicMappingsSelectorUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedTopic: Topic) => {
            setActiveTopicMappings(activeTopicMappings.concat(selectedTopic));
        }

        const handleSubmit = () => {

        }

        const classes = styleClasses();
        return (
            <Form handleSubmit={handleSubmit}>
                <Grid container>
                    <Grid className={classes.formComponent} item xs={12} md={6}>
                        <Autocomplete
                            id="language-code-selector"
                            onChange={handleUpdateLanguageCode}
                            options={Object.values(WordsmithLanguageCode)}
                            value={languageCode}
                            getOptionLabel={(option: WordsmithLanguageCode) => getEnglishNameForLanguageCode(option)}
                            getOptionSelected={(option: WordsmithLanguageCode) => option === languageCode}
                            renderInput={(params) => <PrimaryTextField label="Confirm language" {...params} />} />
                    </Grid>
                    <Grid className={classes.formComponent} item xs={12} md={6}>
                        <Autocomplete
                            id="country-code-selector"
                            onChange={handleCountyCodeUpdate}
                            options={Object.values(CountryCode)}
                            value={countryCode}
                            getOptionLabel={(option: CountryCode) => getEnglishNameForCountryCode(option)}
                            getOptionSelected={(option: CountryCode) => option === countryCode}
                            renderInput={(params) => <PrimaryTextField label="Confirm Country" {...params} />} />
                    </Grid>
                    <Grid className={classes.formComponent} item xs={12}>
                        <Autocomplete
                            id="topic-mapping-selector"
                            onChange={handleTopicMappingsSelectorUpdate}
                            options={props.topics.filter((topic: Topic) => activeTopicMappings.indexOf(topic) === -1)}
                            getOptionLabel={(option: Topic) => option.label}
                            renderInput={(params) => <PrimaryTextField label="Add Topic" {...params} />} />
                    </Grid>
                    <Grid item xs={12}>
                        {
                            activeTopicMappings.map((topic: Topic, idx: number) => (
                                <Grid container
                                    className={classes.alignedContainer}>
                                    <Grid item xs={10} md={11}>
                                        <Paragraph align={Alignment.Left}>
                                            { topic.label }
                                        </Paragraph>
                                    </Grid>
                                    <Grid item xs={2} md={1}>
                                        <ClearIcon
                                            className={classes.removeContentTopicIcon}
                                            onClick={() => {
                                                setActiveTopicMappings(activeTopicMappings.filter((t: Topic) => t !== topic))
                                            }}  />
                                    </Grid>
                                </Grid>
                            ))
                        }
                    </Grid>
                    <CenteredComponent className={classes.formComponent}>
                        <PrimaryButton className={classes.addPodcastFormComponent} type="submit" disabled={!countryCode || !languageCode}>
                            Add Podcast
                        </PrimaryButton>
                    </CenteredComponent>
                </Grid>
            </Form>
        )
    },
    (
        ownProps: PodcastCaptureFormOwnProps,
        onSuccess: (resp: GetAllContentTopicsResponse) => void,
        onError: (err: Error) => void,
    ) => getAllContentTopics({}, onSuccess, onError),
    false,
);

export default PodcastSearchPage;
