import React, { useEffect, useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimarySlider } from 'common/components/Slider/Slider';
import { Heading1, Heading3, Heading4 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';

import {
    PodcastMetadata,

    GetPodcastMetadataResponse,
    getPodcastMetadata,
} from 'ConsumerWeb/api/podcasts/podcasts';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

const seekStepGranularity = 1000;
const styleClasses = makeStyles({
    toggleButton: {
        width: '100%',
    },
});

type Params = {
    userPodcastID: string;
}

type PodcastPlayerPageOwnProps = RouteComponentProps<Params>;

const PodcastPlayerPage = asBaseComponent<GetPodcastMetadataResponse, PodcastPlayerPageOwnProps>(
    (props: GetPodcastMetadataResponse & PodcastPlayerPageOwnProps & BaseComponentProps) => {
        if (!!props.error) {
            return (
                <Heading3 color={TypographyColor.Warning}>
                    Could not play podcast
                </Heading3>
            )
        }
        const [ audio, setAudio ] = useState<HTMLAudioElement>(new Audio(`/vfile/${props.metadata.audioUrl}`));

        const [ isAudioPlaying, setIsAudioPlaying ] = useState<boolean>(false);
        const [ seekValue, setSeekValue ] = useState<number>(0);
        const handleSeekValueChange = (event: Event, newValue: number | number[]) => {
            const val = newValue as number;
            audio.currentTime = val / seekStepGranularity * audio.duration;
            setSeekValue(val);
        };

        useEffect(() => {
            audio.addEventListener("timeupdate", function() {
                setSeekValue(this.currentTime / audio.duration * seekStepGranularity);
            });
            audio.addEventListener("ended", function() {
                setIsAudioPlaying(false);
            });
        }, []);

        const toggleAudio = () => {
            if (isAudioPlaying) {
                audio.pause();
            } else {
                audio.play();
            }
            setIsAudioPlaying(!isAudioPlaying);
        }

        const classes = styleClasses();
        return (
            <CenteredComponent>
                <DisplayCard>
                    <Grid container>
                        <Grid item xs={12}>
                            <Heading1 color={TypographyColor.Primary}>
                                {props.metadata.episodeTitle}
                            </Heading1>
                            <Heading3>
                                {props.metadata.podcastTitle}
                            </Heading3>
                            <Paragraph>
                                {props.metadata.episodeDescription}
                            </Paragraph>
                        </Grid>
                        <Grid item xs={12} md={10}>
                            <PrimarySlider aria-label="Time" max={seekStepGranularity} value={seekValue} onChange={handleSeekValueChange} />
                        </Grid>
                        <Grid item xs={12} md={2}>
                            <CenteredComponent>
                                <PrimaryButton className={classes.toggleButton} onClick={toggleAudio}>
                                    { isAudioPlaying ? "Pause" : "Play" }
                                </PrimaryButton>
                            </CenteredComponent>
                        </Grid>
                    </Grid>
                </DisplayCard>
            </CenteredComponent>
        );
    },
    (
        ownProps: PodcastPlayerPageOwnProps,
        onSuccess: (resp: GetPodcastMetadataResponse) => void,
        onError: (err: Error) => void,
    ) => getPodcastMetadata({
        userPodcastId: ownProps.match.params.userPodcastID,
    }, onSuccess, onError),
    true,
);

export default PodcastPlayerPage;
