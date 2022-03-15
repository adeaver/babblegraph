import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';

import {
    GetCampaignResponse,
    getCampaign,

    GetAllSourcesResponse,
    getAllSources,

    GetAllVendorsResponse,
    getAllVendors,
} from 'AdminWeb/api/advertising/advertising';

import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';

import EditCampaignForm from './EditCampaignForm';

type Params = {
    id: string;
}

type CampaignEditPageProps = RouteComponentProps<Params>;

type CampaignsEditPageAPIProps = GetCampaignResponse & GetAllSourcesResponse & GetAllVendorsResponse;


const CampaignEditPage = asBaseComponent<CampaignsEditPageAPIProps, CampaignEditPageProps>(
    (props: CampaignsEditPageAPIProps & CampaignEditPageProps & BaseComponentProps) => {
        return (
            <CenteredComponent>
                <DisplayCard>
                    <DisplayCardHeader
                        title="Edit Campaign"
                        backArrowDestination="/ops/advertising-manager/campaigns" />
                    <EditCampaignForm
                        campaign={props.campaign}
                        sources={props.sources}
                        vendors={props.vendors}
                        onError={props.setError} />
                </DisplayCard>
            </CenteredComponent>
        )
    },
    (
        ownProps: CampaignEditPageProps,
        onSuccess: (resp: CampaignsEditPageAPIProps) => void,
        onError: (err: Error) => void,
    ) => {
        getCampaign({
            id: ownProps.match.params.id,
        },
        (resp: GetCampaignResponse) => {
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
                (err: Error) => onError(err));
            },
            (err: Error) => onError(err));
        },
        (err: Error) => onError(err));
    },
    true,
);

export default CampaignEditPage;
