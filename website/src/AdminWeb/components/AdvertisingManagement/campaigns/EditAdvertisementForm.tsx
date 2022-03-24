import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Autocomplete from '@material-ui/lab/Autocomplete';

import Form from 'common/components/Form/Form';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimarySwitch } from 'common/components/Switch/Switch';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import {
    Advertisement,

    InsertAdvertisementResponse,
    insertAdvertisement,

    UpdateAdvertisementResponse,
    updateAdvertisement,
} from 'AdminWeb/api/advertising/advertising';

import { WordsmithLanguageCode, getEnglishNameForLanguageCode } from 'common/model/language/language';

const styleClasses = makeStyles({
    formComponent: {
        width: '100%',
        margin: '10px 0',
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

type EditAdvertisementFormProps = {
    campaignID: string;
    advertisement?: Advertisement | undefined;

    handleNewAdvertisement: (a: Advertisement) => void;
    onError: (err: Error) => void;
}


const EditAdvertisementForm = (props: EditAdvertisementFormProps) => {
    const hasAdvertisement = !!props.advertisement;

    const [ isLoading, setIsLoading ] = useState<boolean>(false);

    const [ title, setTitle ] = useState<string>(props.advertisement ? props.advertisement.title : null);
    const handleTitleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setTitle((event.target as HTMLInputElement).value);
    }

    const [ description, setDescription ] = useState<string>(props.advertisement ? props.advertisement.description : null);
    const handleDescriptionChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setDescription((event.target as HTMLInputElement).value);
    }

    const [ imageURL, setImageURL ] = useState<string>(props.advertisement ? props.advertisement.imageUrl : null);
    const handleImageURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setImageURL((event.target as HTMLInputElement).value);
    }

    const [ additionalLinkURL, setAdditionalLinkURL ] = useState<string>(props.advertisement ? props.advertisement.additionalLinkUrl : null);
    const handleAdditionalLinkURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setAdditionalLinkURL((event.target as HTMLInputElement).value);
    }

    const [ additionalLinkText, setAdditionalLinkText ] = useState<string>(props.advertisement ? props.advertisement.additionalLinkText : null);
    const handleAdditionalLinkTextChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setAdditionalLinkText((event.target as HTMLInputElement).value);
    }

    const [ languageCode, setLanguageCode ] = useState<WordsmithLanguageCode>(props.advertisement ? props.advertisement.languageCode : null);
    const handleLanguageCodeUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedLanguageCode: WordsmithLanguageCode) => {
        setLanguageCode(selectedLanguageCode);
    }

    const [ isActive, setIsActive ] = useState<boolean>(props.advertisement ? props.advertisement.isActive : false);

    const handleSubmit = () => {
        setIsLoading(true);
        hasAdvertisement ? (
            updateAdvertisement({
                id: props.advertisement.id,
                isActive: isActive,
                title: title,
                languageCode: languageCode,
                description: description,
                imageUrl: imageURL,
                additionalLinkUrl: !!additionalLinkURL ? additionalLinkURL : undefined,
                additionalLinkText: !!additionalLinkText ? additionalLinkText : undefined,
            },
            (resp: UpdateAdvertisementResponse) => {
                setIsLoading(false);
            },
            (err: Error) => {
                setIsLoading(false);
                props.onError(err);
            })
        ) : (
            insertAdvertisement({
                title: title,
                languageCode: languageCode,
                description: description,
                imageUrl: imageURL,
                campaignId: props.campaignID,
                additionalLinkUrl: !!additionalLinkURL ? additionalLinkURL : undefined,
                additionalLinkText: !!additionalLinkText ? additionalLinkText : undefined,
            },
            (resp: InsertAdvertisementResponse) => {
                setIsLoading(false);
                props.handleNewAdvertisement({
                    id: resp.id,
                    campaignId: props.campaignID,
                    description: description,
                    title: title,
                    languageCode: languageCode,
                    imageUrl: imageURL,
                    isActive: false,
                    additionalLinkUrl: !!additionalLinkURL ? additionalLinkURL : undefined,
                    additionalLinkText: !!additionalLinkText ? additionalLinkText : undefined,
                });
            },
            (err: Error) => {
                setIsLoading(false);
                props.onError(err);
            })
        );
    }

    const classes = styleClasses();
    return (
        <Form handleSubmit={handleSubmit}>
            <Grid className={classes.formContainer} container>
                <Grid item xs={hasAdvertisement ? 9 : 12}>
                    <PrimaryTextField
                        id="advertisement-title"
                        className={classes.formComponent}
                        label="Title"
                        variant="outlined"
                        defaultValue={title}
                        onChange={handleTitleChange} />
                </Grid>
                {
                    hasAdvertisement && (
                        <Grid className={classes.isActiveToggleContainer} item xs={3}>
                            <PrimarySwitch checked={isActive} onClick={() => {setIsActive(!isActive)}} />
                        </Grid>
                    )
                }
                <Grid item xs={12}>
                    <PrimaryTextField
                        id="advertisement-description"
                        className={classes.formComponent}
                        label="Description"
                        variant="outlined"
                        defaultValue={description}
                        onChange={handleDescriptionChange}
                        multiline />
                </Grid>
                <Grid item xs={12}>
                    <PrimaryTextField
                        id="advertisement-image-url"
                        className={classes.formComponent}
                        label="Image URL"
                        variant="outlined"
                        defaultValue={imageURL}
                        onChange={handleImageURLChange} />
                </Grid>
                <Grid item xs={12}>
                    <PrimaryTextField
                        id="advertisement-additional-link-url"
                        className={classes.formComponent}
                        label="Additional Link URL"
                        variant="outlined"
                        defaultValue={additionalLinkURL}
                        onChange={handleAdditionalLinkURLChange} />
                </Grid>
                <Grid item xs={12}>
                    <PrimaryTextField
                        id="advertisement-additional-link-text"
                        className={classes.formComponent}
                        label="Additional Link Text"
                        variant="outlined"
                        defaultValue={additionalLinkText}
                        onChange={handleAdditionalLinkTextChange} />
                </Grid>
                <Grid item xs={12}>
                    <Autocomplete
                        id="language-selector"
                        className={classes.formComponent}
                        onChange={handleLanguageCodeUpdate}
                        options={Object.values(WordsmithLanguageCode)}
                        value={languageCode}
                        getOptionLabel={(option: WordsmithLanguageCode) => getEnglishNameForLanguageCode(option)}
                        getOptionSelected={(option: WordsmithLanguageCode) => option === languageCode}
                        renderInput={(params) => <PrimaryTextField label="Select Language Code" {...params} />} />
                </Grid>
                <Grid item xs={6}>
                    <PrimaryButton
                        className={classes.formComponent}
                        disabled={!imageURL || !title || isLoading || !languageCode || !description}
                        type="submit">
                        { hasAdvertisement ? "Update" : "Create" }
                    </PrimaryButton>
                </Grid>
            </Grid>
            { isLoading && <LoadingSpinner /> }
        </Form>
    );
}

export default EditAdvertisementForm;
