import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import ClearIcon from '@material-ui/icons/Clear';
import Autocomplete from '@material-ui/lab/Autocomplete';

import Color from 'common/styles/colors';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton } from 'common/components/Button/Button';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import DisplayCardHeader from 'common/components/DisplayCard/DisplayCardHeader';

import {
    Advertisement,

    GetCampaignResponse,
    getCampaign,

    GetAllSourcesResponse,
    getAllSources,

    GetAllVendorsResponse,
    getAllVendors,

    GetCampaignTopicMappingsResponse,
    getCampaignTopicMappings,

    UpdateCampaignTopicMappingsResponse,
    updateCampaignTopicMappings,
} from 'AdminWeb/api/advertising/advertising';
import {
    Topic,

    GetAllContentTopicsResponse,
    getAllContentTopics,
} from 'AdminWeb/api/content/topic';

import { asBaseComponent, BaseComponentProps } from 'common/base/BaseComponent';

import EditCampaignForm from './EditCampaignForm';
import EditAdvertisementForm from './EditAdvertisementForm';

const styleClasses = makeStyles({
    formComponent: {
        width: '100%',
        margin: '10px 0',
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

type CampaignEditPageProps = RouteComponentProps<Params>;

type CampaignsEditPageAPIProps = GetCampaignResponse & GetAllSourcesResponse & GetAllVendorsResponse;


const CampaignEditPage = asBaseComponent<CampaignsEditPageAPIProps, CampaignEditPageProps>(
    (props: CampaignsEditPageAPIProps & CampaignEditPageProps & BaseComponentProps) => {
        const handleNewAdvertisement = (ad: Advertisement) => {

        }

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
                    <TopicMappingEditForm campaignID={props.campaign.id} />
                    <Heading3 color={TypographyColor.Primary}>
                        Add a new advertisement in this campaign
                    </Heading3>
                    <EditAdvertisementForm
                        campaignID={props.campaign.id}
                        handleNewAdvertisement={handleNewAdvertisement}
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

type TopicMappingEditFormProps = {
    campaignID: string;
}

type TopicMappingEditFormAPIProps = GetCampaignTopicMappingsResponse & GetAllContentTopicsResponse;

const TopicMappingEditForm = asBaseComponent<GetCampaignTopicMappingsResponse, TopicMappingEditFormProps>(
    (props: TopicMappingEditFormAPIProps & TopicMappingEditFormProps & BaseComponentProps) => {
        const [ activeTopicMappings, setActiveTopicMappings ] = useState<Topic[]>(
            !props.topicIds || !props.topics ? (
                []
            ) : (
                props.topics.filter((t: Topic) => props.topicIds.indexOf(t.id) !== -1)
            )
        );
        const handleTopicMappingsSelectorUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedTopic: Topic) => {
            setActiveTopicMappings(activeTopicMappings.concat(selectedTopic));
        }

        const [ isLoading, setIsLoading ] = useState<boolean>(false);

        const handleSubmit = () => {
            setIsLoading(true);
            updateCampaignTopicMappings({
                campaignId: props.campaignID,
                activeTopicMappings: activeTopicMappings.map((t: Topic) => t.id),
            },
            (resp: UpdateCampaignTopicMappingsResponse) => {
                setIsLoading(false);
            },
            (err: Error) => {
                setIsLoading(false);
                props.setError(err);
            });
        }

        const classes = styleClasses();
        return (
            <Grid container>
                <Grid item xs={12}>
                    <Heading3 color={TypographyColor.Primary}>
                        Topic Mappings for this Campaign
                    </Heading3>
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
                <Grid item xs={12}>
                    <Autocomplete
                        id="topic-mapping-selector"
                        onChange={handleTopicMappingsSelectorUpdate}
                        options={props.topics.filter((topic: Topic) => activeTopicMappings.indexOf(topic) === -1)}
                        getOptionLabel={(option: Topic) => option.label}
                        renderInput={(params) => <PrimaryTextField label="Add Topic" {...params} />} />
                </Grid>
                <Grid item xs={6}>
                    <PrimaryButton
                        className={classes.formComponent}
                        onClick={handleSubmit}
                        disabled={isLoading}
                        type="submit">
                        Update Topics
                    </PrimaryButton>
                </Grid>
            { isLoading && <LoadingSpinner /> }
            </Grid>
        );
    },
    (
        ownProps: TopicMappingEditFormProps,
        onSuccess: (resp: TopicMappingEditFormAPIProps) => void,
        onError: (err: Error) => void,
    ) => {
        getCampaignTopicMappings({
            campaignId: ownProps.campaignID
        },
        (resp: GetCampaignTopicMappingsResponse) => {
            getAllContentTopics({},
            (resp2: GetAllContentTopicsResponse) => {
                onSuccess({
                    ...resp,
                    ...resp2,
                });
            }, onError);
        }, onError);
    },
    false
);

export default CampaignEditPage;
