import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Autocomplete from '@material-ui/lab/Autocomplete';
import FormControlLabel from '@material-ui/core/FormControlLabel';

import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { Heading1, Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import { PrimarySwitch } from 'common/components/Switch/Switch';
import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';
import Link, { LinkTarget } from 'common/components/Link/Link';
import Paragraph from 'common/typography/Paragraph';

import { WordsmithLanguageCode, getEnglishNameForLanguageCode } from 'common/model/language/language';
import { CountryCode, getEnglishNameForCountryCode } from 'common/model/geo/geo';
import {
    Source,
    SourceType,
    IngestStrategy,

    getAllSources,
    GetAllSourcesResponse,
    addSource,
    AddSourceResponse,
    updateSource,
    UpdateSourceResponse,
} from 'AdminWeb/api/content/sources';

const styleClasses = makeStyles({
    addSourceFormCell: {
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        padding: '5px',
    },
    addSourceFormInput: {
        minWidth: '100%',
    },
    sourceDisplayRoot: {
        padding: '5px',
    },
    sourceDisplayHeader: {
        display: 'flex',
        alignItems: 'center',
    },
});

const SourcesListPage = asBaseComponent(
    (props: GetAllSourcesResponse & BaseComponentProps) => {
        const [ newSources, setNewSources ] = useState<Array<Source>>([]);

        const handleNewSource = (newSource: Source) => {
            setNewSources(newSources.concat(newSource));
        }

        const sources = (props.sources || []).concat(newSources);
        return (
            <div>
                <AddSourceForm
                    setIsLoading={props.setIsLoading}
                    setError={props.setError}
                    handleNewSource={handleNewSource} />
                <Grid container>
                {
                    sources.map((s: Source, idx: number) => ((
                        <SourceDisplay
                            key={`sources-list-${idx}`}
                            source={s}
                            setIsLoading={props.setIsLoading}
                            setError={props.setError} />
                    )))
                }
                </Grid>
            </div>
        );
    },
    (
        ownProps: {},
        onSuccess: (resp: GetAllSourcesResponse) => void,
        onError: (err: Error) => void,
    ) => getAllSources({}, onSuccess, onError),
    true,
);

type AddSourceFormProps = {
    handleNewSource: (s: Source) => void;

    setIsLoading: (isLoading: boolean) => void;
    setError: (err: Error) => void;
}

const AddSourceForm = (props: AddSourceFormProps) => {
    const [ title, setTitle ] = useState<string>(null);
    const handleTitleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setTitle((event.target as HTMLInputElement).value);
    }

    const [ url, setURL ] = useState<string>(null);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setURL((event.target as HTMLInputElement).value);
    }

    const [ countryCode, setCountryCode ] = useState<CountryCode>(null);
    const handleCountyCodeUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedCountryCode: CountryCode) => {
        setCountryCode(selectedCountryCode);
    }

    const [ languageCode, setLanguageCode ] = useState<WordsmithLanguageCode>(null);
    const handleLanguageCodeUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedLanguageCode: WordsmithLanguageCode) => {
        setLanguageCode(selectedLanguageCode);
    }

    const [ ingestStrategy, setIngestStrategy ] = useState<IngestStrategy>(null);
    const handleIngestStrategySelectorUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedIngestStrategy: IngestStrategy) => {
        setIngestStrategy(selectedIngestStrategy);
    }

    const [ sourceType, setSourceType ] = useState<SourceType>(null);
    const handleSourceTypeSelectorUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedSourceType: SourceType) => {
        setSourceType(selectedSourceType);
    }

    const [ monthlyAccessLimit, setMonthlyAccessLimit ] = useState<number>(null);
    const handleMonthlyAccessLimitChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const value = (event.target as HTMLInputElement).value;
        setMonthlyAccessLimit(!!value ? parseInt(value) : null);
    }

    const [ shouldUseURLAsSeedURL, setShouldUseURLAsSeedURL ] = useState<boolean>(false);

    const handleSubmit = () => {
        props.setIsLoading(true);
        addSource({
            title: title,
	        url: url,
	        type: sourceType,
	        ingestStrategy: ingestStrategy,
	        languageCode: languageCode,
            monthlyAccessLimit: monthlyAccessLimit,
	        country: countryCode,
            shouldUseUrlAsSeedUrl: shouldUseURLAsSeedURL,
        },
        (resp: AddSourceResponse) => {
            props.setIsLoading(false);
            props.handleNewSource({
	            id: resp.id,
                title: resp.title,
	            url: url,
	            type: sourceType,
	            country: countryCode,
	            ingestStrategy: ingestStrategy,
	            languageCode: languageCode,
	            isActive: false,
	            monthlyAccessLimit: monthlyAccessLimit,
                shouldUseUrlAsSeedUrl: shouldUseURLAsSeedURL,
            });
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
                    <Heading3 color={TypographyColor.Primary}>
                        Add new source
                    </Heading3>
                    <Form handleSubmit={handleSubmit}>
                        <Grid container>
                            <Grid className={classes.addSourceFormCell} item xs={12} md={8}>
                                <PrimaryTextField
                                    id="title"
                                    className={classes.addSourceFormInput}
                                    label="Title (optional)"
                                    variant="outlined"
                                    defaultValue={title}
                                    onChange={handleTitleChange} />
                            </Grid>
                            <Grid className={classes.addSourceFormCell} item xs={6} md={4}>
                                <Autocomplete
                                    id="language-selector"
                                    onChange={handleLanguageCodeUpdate}
                                    options={Object.values(WordsmithLanguageCode)}
                                    getOptionLabel={(option: WordsmithLanguageCode) => getEnglishNameForLanguageCode(option)}
                                    getOptionSelected={(option: WordsmithLanguageCode) => option === languageCode}
                                    renderInput={(params) => <PrimaryTextField label="Select Language Code" {...params} />} />
                            </Grid>
                            <Grid className={classes.addSourceFormCell} item xs={12} md={8}>
                                <PrimaryTextField
                                    id="url"
                                    className={classes.addSourceFormInput}
                                    label="URL"
                                    variant="outlined"
                                    defaultValue={url}
                                    onChange={handleURLChange} />
                            </Grid>
                            <Grid className={classes.addSourceFormCell} item xs={6} md={4}>
                                <FormControlLabel
                                    control={
                                        <PrimaryCheckbox
                                            checked={shouldUseURLAsSeedURL}
                                            onChange={() => { setShouldUseURLAsSeedURL(!shouldUseURLAsSeedURL) }}
                                            name="checkbox-is-seed-url" />
                                    }
                                    label="Should use URL as seed URL??" />
                            </Grid>
                            <Grid className={classes.addSourceFormCell} item xs={6} md={4}>
                                <Autocomplete
                                    id="country-code-selector"
                                    onChange={handleCountyCodeUpdate}
                                    options={Object.values(CountryCode)}
                                    getOptionLabel={(option: CountryCode) => getEnglishNameForCountryCode(option)}
                                    getOptionSelected={(option: CountryCode) => option === countryCode}
                                    renderInput={(params) => <PrimaryTextField label="Select Country" {...params} />} />
                            </Grid>
                            <Grid className={classes.addSourceFormCell} item xs={6} md={4}>
                                <Autocomplete
                                    id="source-type-selector"
                                    onChange={handleSourceTypeSelectorUpdate}
                                    options={Object.values(SourceType)}
                                    getOptionLabel={(option: SourceType) => option}
                                    getOptionSelected={(option: SourceType) => option === sourceType}
                                    renderInput={(params) => <PrimaryTextField label="Select Source Type" {...params} />} />
                            </Grid>
                            <Grid className={classes.addSourceFormCell} item xs={6} md={4}>
                                <Autocomplete
                                    id="ingest-strategy-selector"
                                    onChange={handleIngestStrategySelectorUpdate}
                                    options={Object.values(IngestStrategy)}
                                    getOptionLabel={(option: IngestStrategy) => option}
                                    getOptionSelected={(option: IngestStrategy) => option === ingestStrategy}
                                    renderInput={(params) => <PrimaryTextField label="Select Ingest Strategy" {...params} />} />
                            </Grid>
                            <Grid className={classes.addSourceFormCell} item xs={12} md={8}>
                                <PrimaryTextField
                                    id="monthly-access-limit"
                                    label="Monthly Access Limit"
                                    className={classes.addSourceFormInput}
                                    variant="outlined"
                                    type="number"
                                    defaultValue={monthlyAccessLimit}
                                    onChange={handleMonthlyAccessLimitChange} />
                            </Grid>
                            <Grid className={classes.addSourceFormCell} item xs={3} md={4}>
                                <PrimaryButton disabled={!url || !sourceType || !languageCode || !ingestStrategy || !countryCode} type="submit">
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
        </Grid>
    );
}

type SourceDisplayProps = {
    source: Source,

    setIsLoading: (isLoading: boolean) => void;
    setError: (err: Error) => void;
}

const SourceDisplay = (props: SourceDisplayProps) => {
    const [ isActive, setIsActive ] = useState<boolean>(props.source.isActive);

    const handleToggleSource = () => {
        props.setIsLoading(true);
        updateSource({
            id: props.source.id,
            title: props.source.title,
            url: props.source.url,
            type: props.source.type,
            ingestStrategy: props.source.ingestStrategy,
            monthlyAccessLimit: props.source.monthlyAccessLimit,
            country: props.source.country,
            shouldUseUrlAsSeedUrl: props.source.shouldUseUrlAsSeedUrl,
            languageCode: props.source.languageCode,
            isActive: !isActive,
        },
        (resp: UpdateSourceResponse) => {
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
        <Grid className={classes.sourceDisplayRoot} xs={12} md={4} item>
            <DisplayCard>
                <Grid className={classes.sourceDisplayHeader} container>
                    <Grid item xs={8}>
                        <Heading3 color={isActive ? TypographyColor.Primary : TypographyColor.Gray}>
                            {props.source.title}
                        </Heading3>
                    </Grid>
                    <Grid item xs={4}>
                        <PrimarySwitch checked={isActive} onClick={handleToggleSource} />
                    </Grid>
                </Grid>
                <Paragraph>
                    {getEnglishNameForLanguageCode(props.source.languageCode)}
                </Paragraph>
                <Link href={`/ops/content-manager/sources/${props.source.id}`} target={LinkTarget.Self}>
                    Manage this source
                </Link>
            </DisplayCard>
        </Grid>
    )
}

export default SourcesListPage;
