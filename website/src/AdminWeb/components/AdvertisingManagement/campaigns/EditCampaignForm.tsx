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
    AdvertisementSource,
    Campaign,
    Vendor,

    UpdateCampaignResponse,
    updateCampaign,
} from 'AdminWeb/api/advertising/advertising';

type EditCampaignFormProps = {
    campaign: Campaign;
    sources: Array<AdvertisementSource>;
    vendors: Array<Vendor>;

    onError: (err: Error) => void;
}

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

const EditCampaignForm = (props: EditCampaignFormProps) => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);

    const [ name, setName ] = useState<string>(props.campaign.name);
    const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName((event.target as HTMLInputElement).value);
    }

    const [ url, setURL ] = useState<string>(props.campaign.url);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setURL((event.target as HTMLInputElement).value);
    }

    const [ sourceID, setSourceID ] = useState<string>(props.campaign.sourceId);
    const handleSourceUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedSource: AdvertisementSource) => {
        setSourceID(selectedSource.id);
    }

    const [ vendorID, setVendorID ] = useState<string>(props.campaign.vendorId);
    const handleVendorUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedVendor: Vendor) => {
        setVendorID(selectedVendor.id);
    }

    const [ rolloutPercentage, setRolloutPercentage ] = useState<number>(props.campaign.rolloutPercentage);
    const handleRolloutPercentageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRolloutPercentage(parseInt((event.target as HTMLInputElement).value, 10));
    }

    const [ isActive, setIsActive ] = useState<boolean>(props.campaign.isActive);

    const handleSubmit = () => {
        setIsLoading(true);
        updateCampaign({
            campaignId: props.campaign.id,
            isActive: isActive,
            url: url,
            name: name,
            vendorId: vendorID,
            sourceId: sourceID,
            shouldApplyToAllUsers: false,
            rolloutPercentage: rolloutPercentage,
        },
        (resp: UpdateCampaignResponse) => {
            setIsLoading(false);
        },
        (err: Error) => {
            setIsLoading(false);
            props.onError(err);
        });
    }

    const classes = styleClasses();
    return (
        <Form handleSubmit={handleSubmit}>
            <Grid className={classes.formContainer} container>
                <Grid item xs={9}>
                    <PrimaryTextField
                        id="campaign-name"
                        className={classes.formComponent}
                        label="Campaign Name"
                        variant="outlined"
                        defaultValue={name}
                        onChange={handleNameChange} />
                </Grid>
                <Grid className={classes.isActiveToggleContainer} item xs={3}>
                    <PrimarySwitch checked={isActive} onClick={() => {setIsActive(!isActive)}} />
                </Grid>
                <Grid item xs={12}>
                    <PrimaryTextField
                        id="campaign-url"
                        className={classes.formComponent}
                        label="Campaign URL"
                        variant="outlined"
                        defaultValue={url}
                        onChange={handleURLChange} />
                </Grid>
                <Grid item xs={12}>
                    <Autocomplete
                        id="source-selector"
                        className={classes.formComponent}
                        onChange={handleSourceUpdate}
                        options={props.sources}
                        value={props.sources.filter((s: AdvertisementSource) => s.id === sourceID)[0]}
                        getOptionLabel={(option: AdvertisementSource) => `${option.name} (${option.url})`}
                        getOptionSelected={(option: AdvertisementSource) => option.id === sourceID}
                        renderInput={(params) => <PrimaryTextField label="Select Source" {...params} />} />
                </Grid>
                <Grid item xs={12}>
                    <Autocomplete
                        id="vendor-selector"
                        className={classes.formComponent}
                        onChange={handleVendorUpdate}
                        options={props.vendors}
                        value={props.vendors.filter((s: Vendor) => s.id === vendorID)[0]}
                        getOptionLabel={(option: Vendor) => `${option.name} (${option.websiteUrl})`}
                        getOptionSelected={(option: Vendor) => option.id === vendorID}
                        renderInput={(params) => <PrimaryTextField label="Select Vendor" {...params} />} />
                </Grid>
                <Grid item xs={12}>
                    <PrimaryTextField
                        id="campaign-rollout-percentage"
                        className={classes.formComponent}
                        label="Current Rollout Percentage"
                        variant="outlined"
                        defaultValue={rolloutPercentage}
                        onChange={handleRolloutPercentageChange} />
                </Grid>
                <Grid item xs={6}>
                    <PrimaryButton
                        className={classes.formComponent}
                        disabled={!url || !name || isLoading || !vendorID || !sourceID}
                        type="submit">
                        Update
                    </PrimaryButton>
                </Grid>
            </Grid>
            { isLoading && <LoadingSpinner /> }
        </Form>
    );
}

export default EditCampaignForm;
