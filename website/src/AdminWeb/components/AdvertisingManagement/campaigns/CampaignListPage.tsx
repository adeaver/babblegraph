import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

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

const styleClasses = makeStyles({
    formComponent: {
        width: '100%',
        margin: '10px 0',
    },
    campaignContainer: {
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

type CampaignsListPageAPIProps = GetAllCampaignsResponse & GetAllSourcesResponse & GetAllVendorsResponse;

const CampaignsListPage = asBaseComponent(
    (props: CampaignsListPageAPIProps & BaseComponentProps) => {
        const [ addedCampaigns, setAddedCampaigns ] = useState<Campaign[]>([]);
        const handleAddNewCampaign = (campaign: Campaign) => {
            setAddedCampaigns(addedCampaigns.concat(campaign));
        }

        return (
            <div>
                <AddNewCampaignForm
                    handleAddNewCampaign={handleAddNewCampaign}
                    onError={props.setError} />
                <Grid container>
                {
                    addedCampaigns.concat(props.campaigns || []).map((v: Campaign, idx: number) => (
                        <CampaignDisplay
                            key={`campaign-display-${idx}`}
                            campaign={v}
                            onError={props.setError} />
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

    const handleSubmit = () => {
        setIsLoading(true);
        insertCampaign({
            vendorId: "",
            sourceId: "",
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
                vendorId: "",
                sourceId: "",
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

type CampaignDisplayProps = {
    campaign: Campaign;

    onError: (err: Error) => void;
}

const CampaignDisplay = (props: CampaignDisplayProps) => {
    const [ isLoading, setIsLoading ] = useState<boolean>(false);

    const [ name, setName ] = useState<string>(props.campaign.name);
    const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName((event.target as HTMLInputElement).value);
    }

    const [ url, setURL ] = useState<string>(props.campaign.url);
    const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setURL((event.target as HTMLInputElement).value);
    }

    const [ isActive, setIsActive ] = useState<boolean>(props.campaign.isActive);

    const handleSubmit = () => {
        setIsLoading(true);
        updateCampaign({
            campaignId: props.campaign.id,
            isActive: isActive,
            url: url,
            name: name,
            vendorId: "",
            sourceId: "",
            shouldApplyToAllUsers: false,
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
        <Grid className={classes.campaignContainer} item xs={12} md={4}>
            <DisplayCard>
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

export default CampaignsListPage;
