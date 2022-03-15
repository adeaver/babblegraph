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
import Link, { LinkTarget } from 'common/components/Link/Link';

import {
    AdvertisementSource,
    Campaign,
    Vendor,

    GetAllCampaignsResponse,
    getAllCampaigns,

    InsertCampaignResponse,
    insertCampaign,

    UpdateCampaignResponse,
    updateCampaign,

    GetAllSourcesResponse,
    getAllSources,

    GetAllVendorsResponse,
    getAllVendors,
} from 'AdminWeb/api/advertising/advertising';

import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';

import EditCampaignForm from './EditCampaignForm';

const styleClasses = makeStyles({
    formComponent: {
        width: '100%',
        margin: '10px 0',
    },
    campaignContainer: {
        padding: '10px',
    },
});

type CampaignsListPageAPIProps = GetAllCampaignsResponse & GetAllSourcesResponse & GetAllVendorsResponse;

const CampaignsListPage = asBaseComponent(
    (props: CampaignsListPageAPIProps & BaseComponentProps) => {
        const [ addedCampaigns, setAddedCampaigns ] = useState<Campaign[]>([]);
        const handleAddNewCampaign = (campaign: Campaign) => {
            setAddedCampaigns(addedCampaigns.concat(campaign));
        }

        const classes = styleClasses();
        return (
            <div>
                <AddNewCampaignForm
                    handleAddNewCampaign={handleAddNewCampaign}
                    vendors={props.vendors}
                    sources={props.sources}
                    onError={props.setError} />
                <Grid container>
                {
                    addedCampaigns.concat(props.campaigns || []).map((v: Campaign, idx: number) => (
                        <Grid
                            key={`campaign-display-${idx}`}
                            className={classes.campaignContainer}
                            xs={12}
                            md={4}
                            item>
                            <DisplayCard>
                                <EditCampaignForm
                                    campaign={v}
                                    vendors={props.vendors}
                                    sources={props.sources}
                                    onError={props.setError} />
                                <Link href={`/ops/advertising-manager/campaigns/${v.id}`} target={LinkTarget.Self}>
                                    Manage this campaign
                                </Link>
                            </DisplayCard>
                        </Grid>
                    ))
                }
                </Grid>
            </div>
        );
    },
    (
        ownProps: {},
        onSuccess: (resp: CampaignsListPageAPIProps) => void,
        onError: (err: Error) => void,
    ) => {
        getAllCampaigns({},
        (resp: GetAllCampaignsResponse) => {
            getAllSources({},
                (resp2: GetAllSourcesResponse) => {
                    getAllVendors({},
                        (resp3: GetAllVendorsResponse) => {
                            onSuccess({
                                ...resp,
                                ...resp2,
                                ...resp3,
                            });
                        },
                        onError);
                },
                onError);
        },
        onError);
    },
    true,
);

type AddNewCampaignFormProps = {
    sources: Array<AdvertisementSource>;
    vendors: Array<Vendor>;
    handleAddNewCampaign: (campaign: Campaign) => void;
    onError: (err: Error) => void;
}

const AddNewCampaignForm = (props: AddNewCampaignFormProps) => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);

    const [ name, setName ] = useState<string>(null);
    const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName((event.target as HTMLInputElement).value);
    }

    const [ url, setURL ] = useState<string>(null);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setURL((event.target as HTMLInputElement).value);
    }

    const [ sourceID, setSourceID ] = useState<string>(null);
    const handleSourceUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedSource: AdvertisementSource) => {
        setSourceID(selectedSource.id);
    }

    const [ vendorID, setVendorID ] = useState<string>(null);
    const handleVendorUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedVendor: Vendor) => {
        setVendorID(selectedVendor.id);
    }

    const handleSubmit = () => {
        setIsLoading(true);
        insertCampaign({
            vendorId: vendorID,
            sourceId: sourceID,
            shouldApplyToAllUsers: false,
            name: name,
            url: url,
        },
        (resp: InsertCampaignResponse) => {
            setIsLoading(false);
            props.handleAddNewCampaign({
                id: resp.id,
                name: name,
                url: url,
                vendorId: vendorID,
                sourceId: sourceID,
                shouldApplyToAllUsers: false,
                isActive: false,
                expiresAt: undefined,
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
                    title="Add a new campaign"
                    backArrowDestination="/ops/advertising-manager" />
                    <Form handleSubmit={handleSubmit}>
                        <Grid container>
                            <Grid item xs={12}>
                                <PrimaryTextField
                                    id="campaign-name"
                                    className={classes.formComponent}
                                    label="Campaign Name"
                                    variant="outlined"
                                    defaultValue={name}
                                    onChange={handleNameChange} />
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
                                    getOptionLabel={(option: Vendor) => `${option.name} (${option.websiteUrl})`}
                                    getOptionSelected={(option: Vendor) => option.id === vendorID}
                                    renderInput={(params) => <PrimaryTextField label="Select Vendor" {...params} />} />
                            </Grid>
                            <Grid item xs={6}>
                                <PrimaryButton
                                    className={classes.formComponent}
                                    disabled={!url || !name || isLoading || !vendorID || !sourceID}
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

export default CampaignsListPage;
