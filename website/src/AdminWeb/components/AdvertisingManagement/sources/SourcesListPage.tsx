import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Autocomplete from '@material-ui/lab/Autocomplete';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';
import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { Heading3 } from 'common/typography/Heading';
import { TypographyColor } from 'common/typography/common';
import { PrimarySwitch } from 'common/components/Switch/Switch';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import {
    AdvertisementSource,
    AdvertisementSourceType,

    GetAllSourcesResponse,
    getAllSources,

    InsertSourceResponse,
    insertSource,

    UpdateSourceResponse,
    updateSource,
} from 'AdminWeb/api/advertising/advertising';

import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';

const styleClasses = makeStyles({
    formComponent: {
        width: '100%',
        margin: '10px 0',
    },
    sourceContainer: {
        padding: '10px',
    },
    isActiveToggleContainer: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
    formContainer: {
        alignItems: 'center',
    },
});

const SourcesListPage = asBaseComponent(
    (props: GetAllSourcesResponse & BaseComponentProps) => {
        const [ addedSources, setAddedSources ] = useState<AdvertisementSource[]>([]);
        const handleAddNewSource = (source: AdvertisementSource) => {
            setAddedSources(addedSources.concat(source));
        }

        return (
            <div>
                <AddNewSourceForm
                    handleAddNewSource={handleAddNewSource}
                    onError={props.setError} />
                <Grid container>
                {
                    addedSources.concat(props.sources || []).map((v: AdvertisementSource, idx: number) => (
                        <SourceDisplay
                            key={`source-display-${idx}`}
                            source={v}
                            onError={props.setError} />
                    ))
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

type AddNewSourceFormProps = {
    handleAddNewSource: (source: AdvertisementSource) => void;
    onError: (err: Error) => void;
}

const AddNewSourceForm = (props: AddNewSourceFormProps) => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);

    const [ name, setName ] = useState<string>(null);
    const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName((event.target as HTMLInputElement).value);
    }

    const [ url, setURL ] = useState<string>(null);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setURL((event.target as HTMLInputElement).value);
    }

    const [ sourceType, setSourceType ] = useState<AdvertisementSourceType>(null);
    const handleSourceTypeUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedSourceType: AdvertisementSourceType) => {
        setSourceType(selectedSourceType);
    }

    const handleSubmit = () => {
        setIsLoading(true);
        insertSource({
            name: name,
            websiteUrl: url,
            sourceType: sourceType,
        },
        (resp: InsertSourceResponse) => {
            setIsLoading(false);
            props.handleAddNewSource({
                name: name,
                url: url,
                id: resp.id,
                type: sourceType,
                isActive: false,
            });
        },
        (err: Error) => {
            setIsLoading(false);
            props.onError(err);
        });
    }

    const classes = styleClasses();
    return (
        <CenteredComponent>
            <DisplayCard>
                <DisplayCardHeader
                    title="Add a new source"
                    backArrowDestination="/ops/advertising-manager" />
                    <Form handleSubmit={handleSubmit}>
                        <Grid container>
                            <Grid item xs={12}>
                                <PrimaryTextField
                                    id="source-name"
                                    className={classes.formComponent}
                                    label="Source Name"
                                    variant="outlined"
                                    defaultValue={name}
                                    onChange={handleNameChange} />
                            </Grid>
                            <Grid item xs={12}>
                                <PrimaryTextField
                                    id="source-url"
                                    className={classes.formComponent}
                                    label="Source URL"
                                    variant="outlined"
                                    defaultValue={url}
                                    onChange={handleURLChange} />
                            </Grid>
                            <Grid item xs={12}>
                                <Autocomplete
                                    id="source-type-selector"
                                    className={classes.formComponent}
                                    onChange={handleSourceTypeUpdate}
                                    options={Object.values(AdvertisementSourceType)}
                                    getOptionSelected={(option: AdvertisementSourceType) => option === sourceType}
                                    renderInput={(params) => <PrimaryTextField label="Select Source Type" {...params} />} />
                            </Grid>
                            <Grid item xs={6}>
                                <PrimaryButton
                                    className={classes.formComponent}
                                    disabled={!url || !name || isLoading}
                                    type="submit">
                                    Submit
                                </PrimaryButton>
                            </Grid>
                        </Grid>
                    </Form>
                    { isLoading && <LoadingSpinner /> }
            </DisplayCard>
        </CenteredComponent>
    );
}

type SourceDisplayProps = {
    source: AdvertisementSource;

    onError: (err: Error) => void;
}

const SourceDisplay = (props: SourceDisplayProps) => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);

    const [ name, setName ] = useState<string>(props.source.name);
    const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName((event.target as HTMLInputElement).value);
    }

    const [ url, setURL ] = useState<string>(props.source.url);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setURL((event.target as HTMLInputElement).value);
    }

    const [ sourceType, setSourceType ] = useState<AdvertisementSourceType>(props.source.type);
    const handleSourceTypeUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedSourceType: AdvertisementSourceType) => {
        setSourceType(selectedSourceType);
    }

    const [ isActive, setIsActive ] = useState<boolean>(props.source.isActive);

    const handleSubmit = () => {
        setIsLoading(true);
        updateSource({
            id: props.source.id,
            isActive: isActive,
            websiteUrl: url,
            name: name,
            sourceType: props.source.type,
        },
        (resp: UpdateSourceResponse) => {
            setIsLoading(false);
        },
        (err: Error) => {
            setIsLoading(false);
            props.onError(err);
        });
    }

    const classes = styleClasses();
    return (
        <Grid className={classes.sourceContainer} item xs={12} md={4}>
            <DisplayCard>
                <Form handleSubmit={handleSubmit}>
                    <Grid className={classes.formContainer} container>
                        <Grid item xs={9}>
                            <PrimaryTextField
                                id="source-name"
                                className={classes.formComponent}
                                label="Source Name"
                                variant="outlined"
                                defaultValue={name}
                                onChange={handleNameChange} />
                        </Grid>
                        <Grid className={classes.isActiveToggleContainer} item xs={3}>
                            <PrimarySwitch checked={isActive} onClick={() => {setIsActive(!isActive)}} />
                        </Grid>
                        <Grid item xs={12}>
                            <PrimaryTextField
                                id="source-url"
                                className={classes.formComponent}
                                label="Source URL"
                                variant="outlined"
                                defaultValue={url}
                                onChange={handleURLChange} />
                        </Grid>
                        <Grid item xs={12}>
                            <Autocomplete
                                id="source-type-selector"
                                className={classes.formComponent}
                                onChange={handleSourceTypeUpdate}
                                options={Object.values(AdvertisementSourceType)}
                                getOptionSelected={(option: AdvertisementSourceType) => option === sourceType}
                                renderInput={(params) => <PrimaryTextField label="Select Source Type" {...params} />} />
                        </Grid>
                        <Grid item xs={6}>
                            <PrimaryButton
                                className={classes.formComponent}
                                disabled={!url || !name}
                                type="submit">
                                Update
                            </PrimaryButton>
                        </Grid>
                    </Grid>
                </Form>
                { isLoading && <LoadingSpinner /> }
            </DisplayCard>
        </Grid>
    );
}

export default SourcesListPage;
