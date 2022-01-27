import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Autocomplete from '@material-ui/lab/Autocomplete';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import ClearIcon from '@material-ui/icons/Clear';

import { asBaseComponent, BaseComponentProps } from 'AdminWeb/common/Base/BaseComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Color from 'common/styles/colors';
import { Heading1, Heading3, Heading5 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';
import Form from 'common/components/Form/Form';

import { WordsmithLanguageCode, getEnglishNameForLanguageCode } from 'common/model/language/language';
import { CountryCode, getEnglishNameForCountryCode } from 'common/model/geo/geo';
import {
    Source,
    SourceSeed,
    SourceFilter,
    SourceSeedTopicMapping,
    SourceType,
    IngestStrategy,

    getSourceByID,
    GetSourceByIDResponse,
    updateSource,
    UpdateSourceResponse,

    getAllSourceSeedsForSource,
    GetAllSourceSeedsForSourceResponse,
    addSourceSeed,
    AddSourceSeedResponse,

    upsertSourceSeedMappings,
    UpsertSourceSeedMappingsResponse,

    getSourceSourceSeedMappingsForSource,
    GetSourceSourceSeedMappingsForSourceResponse,

    getSourceFilterForSourceID,
    GetSourceFilterForSourceIDResponse,

    upsertSourceFilterForSource,
    UpsertSourceFilterForSourceResponse,
} from 'AdminWeb/api/content/sources';
import {
    Topic,

    getAllContentTopics,
    GetAllContentTopicsResponse,
} from 'AdminWeb/api/content/topic';

const styleClasses = makeStyles({
    headerDisplayCard: {
        padding: '5px',
    },
    headerWithCheckbox: {
        display: 'flex',
        alignItems: 'center',
    },
    updateSourceFormCell: {
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        padding: '5px',
    },
    updateSourceFormInput: {
        minWidth: '100%',
    },
    sourceDisplayRoot: {
        padding: '5px',
    },
    sourceDisplayHeader: {
        display: 'flex',
        alignItems: 'center',
    },
    sourceSeedDisplayContainer: {
        padding: '5px',
    },
    alignedContainer: {
        display: 'flex',
        alignItems: 'center',
    },
    removeContentTopicIcon: {
        color: Color.Warning,
    },
});

type Params = {
    id: string;
}

type SourceManagementPageOwnProps = RouteComponentProps<Params>;

type SourceManagementPageProps = {
    source: Source,
    allTopics: Array<Topic>;
}

const SourceManagementPage = asBaseComponent<SourceManagementPageProps, SourceManagementPageOwnProps>(
    (props: BaseComponentProps & SourceManagementPageProps & SourceManagementPageOwnProps) => {
        if (!props.source) {
            return <div />;
        }

        return (
            <div>
                <Heading1 color={TypographyColor.Primary}>
                    {props.source.title}
                </Heading1>
                <Grid container>
                    <UpdateSourceForm
                        setIsLoading={props.setIsLoading}
                        setError={props.setError}
                        source={props.source} />
                    <SourceFilterForm
                        sourceId={props.source.id} />
                </Grid>
                <SourceSeedsList
                    allTopics={props.allTopics}
                    sourceId={props.source.id} />
            </div>
        );
    },
    (
        ownProps: SourceManagementPageOwnProps,
        onSuccess: (resp: SourceManagementPageProps) => void,
        onError: (err: Error) => void,
    ) => {
        getSourceByID({
            id: ownProps.match.params.id,
        },
        (resp: GetSourceByIDResponse) => {
            const source = resp.source;
            getAllContentTopics({},
            (resp: GetAllContentTopicsResponse) => {
                onSuccess({
                    source: source,
                    allTopics: resp.topics,
                });
            },
            onError);
        },
        onError)
    },
    true
);

type UpdateSourceFormProps = {
    source: Source;

    setIsLoading: (isLoading: boolean) => void;
    setError: (err: Error) => void;
}

const UpdateSourceForm = (props: UpdateSourceFormProps) => {
    const [ isActive, setIsActive ] = useState<boolean>(props.source.isActive);
    const [ shouldUseURLAsSeedURL, setShouldUseURLAsSeedURL ] = useState<boolean>(props.source.shouldUseUrlAsSeedUrl);

    const [ title, setTitle ] = useState<string>(props.source.title);
    const handleTitleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setTitle((event.target as HTMLInputElement).value);
    }

    const [ url, setURL ] = useState<string>(props.source.url);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setURL((event.target as HTMLInputElement).value);
    }

    const [ countryCode, setCountryCode ] = useState<CountryCode>(props.source.country);
    const handleCountyCodeUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedCountryCode: CountryCode) => {
        setCountryCode(selectedCountryCode);
    }

    const [ languageCode, setLanguageCode ] = useState<WordsmithLanguageCode>(props.source.languageCode);
    const handleLanguageCodeUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedLanguageCode: WordsmithLanguageCode) => {
        setLanguageCode(selectedLanguageCode);
    }

    const [ ingestStrategy, setIngestStrategy ] = useState<IngestStrategy>(props.source.ingestStrategy);
    const handleIngestStrategySelectorUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedIngestStrategy: IngestStrategy) => {
        setIngestStrategy(selectedIngestStrategy);
    }

    const [ sourceType, setSourceType ] = useState<SourceType>(props.source.type);
    const handleSourceTypeSelectorUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedSourceType: SourceType) => {
        setSourceType(selectedSourceType);
    }

    const [ monthlyAccessLimit, setMonthlyAccessLimit ] = useState<number>(props.source.monthlyAccessLimit);
    const handleMonthlyAccessLimitChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const value = (event.target as HTMLInputElement).value;
        setMonthlyAccessLimit(!!value ? parseInt(value) : null);
    }

    const handleSubmit = () => {
        props.setIsLoading(true);
        updateSource({
            id: props.source.id,
            languageCode: languageCode,
            title: title,
            url: url,
            type: sourceType,
            ingestStrategy: ingestStrategy,
            monthlyAccessLimit: monthlyAccessLimit,
            country: countryCode,
            isActive: isActive,
            shouldUseUrlAsSeedUrl: shouldUseURLAsSeedURL,
        },
        (resp: UpdateSourceResponse) => {
            props.setIsLoading(false);
        },
        (err: Error) => {
            props.setIsLoading(false);
            props.setError(err);
        });
    }

    const classes = styleClasses();
    return (
        <Grid item xs={12} md={6} className={classes.headerDisplayCard}>
            <DisplayCard>
                <Form handleSubmit={handleSubmit}>
                    <Grid container>
                        <Grid className={classes.updateSourceFormCell} item xs={12} md={8}>
                            <PrimaryTextField
                                id="title"
                                className={classes.updateSourceFormInput}
                                label="Title"
                                variant="outlined"
                                defaultValue={title}
                                onChange={handleTitleChange} />
                        </Grid>
                        <Grid className={classes.updateSourceFormCell} item xs={6} md={4}>
                            <Autocomplete
                                id="language-selector"
                                onChange={handleLanguageCodeUpdate}
                                options={Object.values(WordsmithLanguageCode)}
                                value={languageCode}
                                getOptionLabel={(option: WordsmithLanguageCode) => getEnglishNameForLanguageCode(option)}
                                getOptionSelected={(option: WordsmithLanguageCode) => option === languageCode}
                                renderInput={(params) => <PrimaryTextField label="Select Language Code" {...params} />} />
                        </Grid>
                        <Grid className={classes.updateSourceFormCell} item xs={12} md={8}>
                            <PrimaryTextField
                                id="url"
                                className={classes.updateSourceFormInput}
                                label="URL"
                                variant="outlined"
                                defaultValue={url}
                                onChange={handleURLChange} />
                        </Grid>
                        <Grid className={classes.updateSourceFormCell} item xs={6} md={4}>
                            <FormControlLabel
                                control={
                                    <PrimaryCheckbox
                                        checked={shouldUseURLAsSeedURL}
                                        onChange={() => { setShouldUseURLAsSeedURL(!shouldUseURLAsSeedURL) }}
                                        name="checkbox-is-seed-url" />
                                }
                                label="Should use URL as seed URL??" />
                        </Grid>
                        <Grid className={classes.updateSourceFormCell} item xs={6} md={4}>
                            <Autocomplete
                                id="country-code-selector"
                                onChange={handleCountyCodeUpdate}
                                options={Object.values(CountryCode)}
                                value={countryCode}
                                getOptionLabel={(option: CountryCode) => getEnglishNameForCountryCode(option)}
                                getOptionSelected={(option: CountryCode) => option === countryCode}
                                renderInput={(params) => <PrimaryTextField label="Select Country" {...params} />} />
                        </Grid>
                        <Grid className={classes.updateSourceFormCell} item xs={6} md={4}>
                            <Autocomplete
                                id="source-type-selector"
                                onChange={handleSourceTypeSelectorUpdate}
                                options={Object.values(SourceType)}
                                value={sourceType}
                                getOptionLabel={(option: SourceType) => option}
                                getOptionSelected={(option: SourceType) => option === sourceType}
                                renderInput={(params) => <PrimaryTextField label="Select Source Type" {...params} />} />
                        </Grid>
                        <Grid className={classes.updateSourceFormCell} item xs={6} md={4}>
                            <Autocomplete
                                id="ingest-strategy-selector"
                                onChange={handleIngestStrategySelectorUpdate}
                                options={Object.values(IngestStrategy)}
                                value={ingestStrategy}
                                getOptionLabel={(option: IngestStrategy) => option}
                                getOptionSelected={(option: IngestStrategy) => option === ingestStrategy}
                                renderInput={(params) => <PrimaryTextField label="Select Ingest Strategy" {...params} />} />
                        </Grid>
                        <Grid className={classes.updateSourceFormCell} item xs={12} md={5}>
                            <PrimaryTextField
                                id="monthly-access-limit"
                                label="Monthly Access Limit"
                                className={classes.updateSourceFormInput}
                                variant="outlined"
                                type="number"
                                defaultValue={monthlyAccessLimit}
                                onChange={handleMonthlyAccessLimitChange} />
                        </Grid>
                        <Grid className={classes.updateSourceFormCell} item xs={3}>
                            <FormControlLabel
                                control={
                                    <PrimaryCheckbox
                                        checked={isActive}
                                        onChange={() => { setIsActive(!isActive) }}
                                        name="checkbox-is-active" />
                                }
                                label="Is Active?" />
                        </Grid>
                        <Grid className={classes.updateSourceFormCell} item xs={3} md={4}>
                            <PrimaryButton disabled={!url || !sourceType || !languageCode || !ingestStrategy || !countryCode} type="submit">
                                Update
                            </PrimaryButton>
                        </Grid>
                    </Grid>
                </Form>
            </DisplayCard>
        </Grid>
    );
}

type SourceFilterFormOwnProps = {
    sourceId: string;
}

type SourceFilterFormApiProps = {
    hasLoaded: boolean;
} & GetSourceFilterForSourceIDResponse;

const SourceFilterForm = asBaseComponent<SourceFilterFormApiProps, SourceFilterFormOwnProps>(
    (props: BaseComponentProps & SourceFilterFormApiProps & SourceFilterFormOwnProps) => {
        if (!props.hasLoaded) {
            return <Grid item xs={12} md={6} />;
        }

        const [ isActive, setIsActive ] = useState<boolean>(!!props.sourceFilter && props.sourceFilter.isActive);

        const [ useLDJSONValidation, setUseLDJSONValidation ] = useState<boolean>(!!props.sourceFilter && props.sourceFilter.useLdJsonValidation);

        const [ paywallClassInput, setPaywallClassInput ] = useState<string>(null);
        const handlePaywallClassInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setPaywallClassInput((event.target as HTMLInputElement).value);
        }
        const [ paywallClasses, setPaywallClasses ] = useState<Array<string>>(
            !props.sourceFilter || !props.sourceFilter.paywallClasses ? (
                []
            ) : (
                props.sourceFilter.paywallClasses
            )
        );
        const addPaywallClass = () => {
            if (paywallClasses.some((paywallClass: string) => paywallClass === paywallClassInput)) {
                return;
            }
            setPaywallClasses(paywallClasses.concat(paywallClassInput));
            setPaywallClassInput(null);
        }

        const [ paywallIDInput, setPaywallIDInput ] = useState<string>(null);
        const handlePaywallIDInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setPaywallIDInput((event.target as HTMLInputElement).value);
        }
        const [ paywallIDs, setPaywallIDs ] = useState<Array<string>>(
            !props.sourceFilter || !props.sourceFilter.paywallIds ? (
                []
            ) : (
                props.sourceFilter.paywallIds
            )
        );
        const addPaywallID = () => {
            if (paywallIDs.some((paywallID: string) => paywallID === paywallIDInput)) {
                return;
            }
            setPaywallIDs(paywallIDs.concat(paywallIDInput));
            setPaywallIDInput(null);
        }

        const handleUpdateSourceFilter = () => {
            props.setIsLoading(true);
            upsertSourceFilterForSource({
                sourceId: props.sourceId,
                isActive: isActive,
                useLdJsonValidation: useLDJSONValidation,
                paywallClasses: paywallClasses,
                paywallIds: paywallIDs
            },
            (resp: UpsertSourceFilterForSourceResponse) => {
                props.setIsLoading(false);
            },
            (err: Error) => {
                props.setIsLoading(false);
                props.setError(err);
            });
        }

        const classes = styleClasses();
        return (
            <Grid item xs={12} md={6} className={classes.headerDisplayCard}>
                <DisplayCard>
                    <Grid container className={classes.headerWithCheckbox}>
                        <Grid item xs={2} md={1}>
                            <PrimaryCheckbox
                                checked={isActive}
                                onChange={() => { setIsActive(!isActive) }}
                                name="checkbox-is-active" />
                        </Grid>
                        <Grid item xs={10} md={11}>
                            <Heading3 align={Alignment.Left} color={TypographyColor.Primary}>
                                Source Filter
                            </Heading3>
                        </Grid>
                    </Grid>
                    <Grid className={classes.updateSourceFormCell} item xs={12}>
                        <FormControlLabel
                            control={
                                <PrimaryCheckbox
                                    checked={useLDJSONValidation}
                                    onChange={() => { setUseLDJSONValidation(!useLDJSONValidation) }}
                                    name="checkbox-ld-json" />
                            }
                            label="Use LD+JSON Validation?" />
                    </Grid>
                    <Form handleSubmit={addPaywallClass}>
                        <Grid container>
                            <Grid item xs={8} md={10} className={classes.updateSourceFormCell}>
                                <PrimaryTextField
                                    id="paywall-classes"
                                    className={classes.updateSourceFormInput}
                                    label="Add Paywall Classes"
                                    variant="outlined"
                                    value={paywallClassInput}
                                    onChange={handlePaywallClassInputChange} />
                            </Grid>
                            <Grid item xs={4} md={2} className={classes.updateSourceFormCell}>
                                <PrimaryButton disabled={!paywallClassInput} type="submit">
                                    Add
                                </PrimaryButton>
                            </Grid>
                        </Grid>
                    </Form>
                    <Grid item xs={12}>
                            {
                                paywallClasses.map((paywallClass: string, idx: number) => (
                                    <Grid container
                                        className={classes.alignedContainer}>
                                        <Grid item xs={10} md={11}>
                                            <Paragraph align={Alignment.Left}>
                                                { paywallClass }
                                            </Paragraph>
                                        </Grid>
                                        <Grid item xs={2} md={1}>
                                            <ClearIcon
                                                className={classes.removeContentTopicIcon}
                                                onClick={() => {
                                                    setPaywallClasses(paywallClasses.filter((p: string) => p !== paywallClass))
                                                }}  />
                                        </Grid>
                                    </Grid>
                                ))
                            }
                    </Grid>
                    <Form handleSubmit={addPaywallID}>
                        <Grid container>
                            <Grid item xs={8} md={10} className={classes.updateSourceFormCell}>
                                <PrimaryTextField
                                    id="paywall-id"
                                    className={classes.updateSourceFormInput}
                                    label="Add Paywall IDs"
                                    variant="outlined"
                                    value={paywallIDInput}
                                    onChange={handlePaywallIDInputChange} />
                            </Grid>
                            <Grid item xs={4} md={2} className={classes.updateSourceFormCell}>
                                <PrimaryButton disabled={!paywallIDInput} type="submit">
                                    Add
                                </PrimaryButton>
                            </Grid>
                        </Grid>
                    </Form>
                    <Grid item xs={12}>
                            {
                                paywallIDs.map((paywallID: string, idx: number) => (
                                    <Grid container
                                        className={classes.alignedContainer}>
                                        <Grid item xs={10} md={11}>
                                            <Paragraph align={Alignment.Left}>
                                                { paywallID }
                                            </Paragraph>
                                        </Grid>
                                        <Grid item xs={2} md={1}>
                                            <ClearIcon
                                                className={classes.removeContentTopicIcon}
                                                onClick={() => {
                                                    setPaywallIDs(paywallIDs.filter((p: string) => p !== paywallID))
                                                }}  />
                                        </Grid>
                                    </Grid>
                                ))
                            }
                    </Grid>
                    <PrimaryButton onClick={handleUpdateSourceFilter} type="submit">
                        Update Source Filter
                    </PrimaryButton>
                    <Paragraph color={TypographyColor.Primary} size={Size.Small}>
                        Refresh after updating to see up-to-date topic mappings because there’s a bug and I'm moving too quickly to fix it
                    </Paragraph>
                </DisplayCard>
            </Grid>
        );
    },
    (
        ownProps: SourceFilterFormOwnProps,
        onSuccess: (resp: SourceFilterFormApiProps) => void,
        onError: (err: Error) => void
    ) => getSourceFilterForSourceID({
        sourceId: ownProps.sourceId,
    },
    (resp: GetSourceFilterForSourceIDResponse) => {
        onSuccess({
            ...resp,
            hasLoaded: true,
        });
    },
    onError),
    false,
);

type SourceSeedsListOwnProps = {
    sourceId: string;
    allTopics: Array<Topic>;
}

type SourceSeedListApiProps = GetAllSourceSeedsForSourceResponse &  GetSourceSourceSeedMappingsForSourceResponse;

const SourceSeedsList = asBaseComponent<SourceSeedListApiProps, SourceSeedsListOwnProps>(
    (props: BaseComponentProps & SourceSeedListApiProps & SourceSeedsListOwnProps) => {

        const sourceSeedMappingsBySourceSeedID = (props.sourceSeedMappings || []).reduce(
            (acc: { [key: string]: Array<SourceSeedTopicMapping> }, s: SourceSeedTopicMapping) => ({
                ...acc,
                [s.sourceSeedId]: (acc[s.sourceSeedId] || []).concat(s),
            }),
        {});

        const [ newSourceSeeds, setNewSourceSeeds ] = useState<Array<SourceSeed>>([]);
        const handleNewSourceSeed = (s: SourceSeed) => {
            setNewSourceSeeds(newSourceSeeds.concat(s));
        }

        const sourceSeeds = (props.sourceSeeds || []).concat(newSourceSeeds);
        return (
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <Heading3 color={TypographyColor.Primary}>
                            Add new seed
                        </Heading3>
                        <AddNewSourceSeedForm
                            sourceId={props.sourceId}
                            handleNewSourceSeed={handleNewSourceSeed}
                            setIsLoading={props.setIsLoading}
                            setError={props.setError} />
                    </DisplayCard>
                </Grid>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                {
                    sourceSeeds.map((s: SourceSeed, idx: number) => (
                        <SourceSeedDisplay
                            key={`source-seed-display-${idx}`}
                            sourceSeed={s}
                            topicMappings={sourceSeedMappingsBySourceSeedID[s.id] || []}
                            allTopics={props.allTopics}
                            setIsLoading={props.setIsLoading}
                            setError={props.setError} />
                    ))
                }
            </Grid>
        );
    },
    (
        ownProps: SourceSeedsListOwnProps,
        onSuccess: (resp: SourceSeedListApiProps) => void,
        onError: (err: Error) => void
    ) => getAllSourceSeedsForSource({
        sourceId: ownProps.sourceId
        },
        (resp: GetAllSourceSeedsForSourceResponse) => {
            const sourceSeeds = resp.sourceSeeds;
            getSourceSourceSeedMappingsForSource({
                sourceId: ownProps.sourceId,
            },
            (resp: GetSourceSourceSeedMappingsForSourceResponse) => {
                onSuccess({
                    ...resp,
                    sourceSeeds: sourceSeeds,
                })
            },
            onError);
        },
        onError),
    false
);

type AddNewSourceSeedFormProps = {
    sourceId: string;
    handleNewSourceSeed: (s: SourceSeed) => void;

    setIsLoading: (isLoading: boolean) => void;
    setError: (err: Error) => void;
}

const AddNewSourceSeedForm = (props: AddNewSourceSeedFormProps) => {
    const [ url, setURL ] = useState<string>(null);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setURL((event.target as HTMLInputElement).value);
    }

    const handleSubmit = () => {
        props.setIsLoading(true);
        addSourceSeed({
            sourceId: props.sourceId,
            url: url,
        },
        (resp: AddSourceSeedResponse) => {
            props.setIsLoading(false);
            props.handleNewSourceSeed({
                id: resp.id,
                rootId: props.sourceId,
                url: url,
                isActive: false,
            });
        },
        (err: Error) => {
            props.setIsLoading(false);
            props.setError(err);
        });
    }

    const classes = styleClasses();
    return (
        <Form handleSubmit={handleSubmit}>
            <Grid container>
                <Grid className={classes.updateSourceFormCell} item xs={12} md={8}>
                    <PrimaryTextField
                        id="url"
                        className={classes.updateSourceFormInput}
                        label="URL"
                        variant="outlined"
                        defaultValue={url}
                        onChange={handleURLChange} />
                </Grid>
                <Grid className={classes.updateSourceFormCell} item xs={3} md={4}>
                    <PrimaryButton disabled={!url} type="submit">
                        Add
                    </PrimaryButton>
                </Grid>
            </Grid>
        </Form>
    );
}

type SourceSeedDisplayProps = {
    sourceSeed: SourceSeed;
    allTopics: Array<Topic>;
    topicMappings: Array<SourceSeedTopicMapping>;

    setIsLoading: (isLoading: boolean) => void;
    setError: (err: Error) => void;
}

const SourceSeedDisplay = (props: SourceSeedDisplayProps) => {
    const [ url, setURL ] = useState<string>(props.sourceSeed.url);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setURL((event.target as HTMLInputElement).value);
    }
    const [ isActive, setIsActive ] = useState<boolean>(props.sourceSeed.isActive);

    const [ activeTopicMappings, setActiveTopicMappings ] = useState<Array<Topic>>(
        props.topicMappings
            .filter((s: SourceSeedTopicMapping) => s.isActive && props.allTopics.some((t: Topic) => t.id === s.topicId))
            .map((s: SourceSeedTopicMapping) => props.allTopics.reduce((acc: Topic | null, t: Topic) => t.id === s.topicId ? t : acc, null)!)
    );
    const [ inactiveTopicMappings, setInactiveTopicMappings ] = useState<Array<Topic>>(
        props.topicMappings
            .filter((s: SourceSeedTopicMapping) => !s.isActive && props.allTopics.some((t: Topic) => t.id === s.topicId))
            .map((s: SourceSeedTopicMapping) => props.allTopics.reduce((acc: Topic | null, t: Topic) => t.id === s.topicId ? t : acc, null)!)
);
    const handleTopicMappingsSelectorUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedTopic: Topic) => {
        setActiveTopicMappings(activeTopicMappings.concat(selectedTopic));
        setInactiveTopicMappings(inactiveTopicMappings.filter((t: Topic) => t !== selectedTopic));
    }

    const handleSubmit = () => {
        props.setIsLoading(true);
        const updates = [{
            sourceSeedId: props.sourceSeed.id,
            isActive: true,
            topicIds: activeTopicMappings.map((t: Topic) => t.id),
        }, {
            sourceSeedId: props.sourceSeed.id,
            isActive: false,
            topicIds: inactiveTopicMappings.map((t: Topic) => t.id),
        }];
        upsertSourceSeedMappings({
            updates: updates,
        },
        (resp: UpsertSourceSeedMappingsResponse) => {
            props.setIsLoading(false);
        },
        (err: Error) => {
            props.setIsLoading(false);
            props.setError(err);
        });
    }

    const classes = styleClasses();
    return (
        <Grid className={classes.sourceSeedDisplayContainer} item xs={4}>
            <DisplayCard>
                <Form handleSubmit={handleSubmit}>
                    <Grid container>
                        <Grid className={classes.updateSourceFormCell} item xs={12} md={8}>
                            <PrimaryTextField
                                id="url"
                                className={classes.updateSourceFormInput}
                                label="URL"
                                variant="outlined"
                                defaultValue={url}
                                onChange={handleURLChange} />
                        </Grid>
                        <Grid className={classes.updateSourceFormCell} item xs={6} md={4}>
                            <FormControlLabel
                                control={
                                    <PrimaryCheckbox
                                        checked={isActive}
                                        onChange={() => { setIsActive(!isActive) }}
                                        name="checkbox-is-active" />
                                }
                                label="Active?" />
                        </Grid>
                        <Heading5 color={TypographyColor.Primary}>
                            Topic Mappings
                        </Heading5>
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
                                                    setInactiveTopicMappings(inactiveTopicMappings.concat(topic));
                                                }}  />
                                        </Grid>
                                    </Grid>
                                ))
                            }
                        </Grid>
                        <Grid item xs={12}>
                            <Autocomplete
                                id="topic-mapping-selector"
                                onChange={handleTopicMappingsSelectorUpdate}
                                options={props.allTopics.filter((topic: Topic) => activeTopicMappings.indexOf(topic) === -1)}
                                getOptionLabel={(option: Topic) => option.label}
                                renderInput={(params) => <PrimaryTextField label="Add Topic" {...params} />} />
                        </Grid>
                        <Grid className={classes.updateSourceFormCell} item xs={3} md={4}>
                            <PrimaryButton disabled={!url} type="submit">
                                Update
                            </PrimaryButton>
                        </Grid>
                    </Grid>
                    <Paragraph color={TypographyColor.Primary} size={Size.Small}>
                        Refresh after updating to see up-to-date topic mappings because there’s a bug and I'm moving too quickly to fix it
                    </Paragraph>
                </Form>
            </DisplayCard>
        </Grid>
    );
}

export default SourceManagementPage;
