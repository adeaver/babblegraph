import React, { useEffect, useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import CircularProgress from '@material-ui/core/CircularProgress';
import VolumeUpIcon from '@material-ui/icons/VolumeUp';
import Autocomplete from '@material-ui/lab/Autocomplete';

import CenteredComponent from 'common/components/CenteredComponent/CenteredComponent';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimarySlider } from 'common/components/Slider/Slider';
import { Heading1, Heading3, Heading4 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import Paragraph, { Size } from 'common/typography/Paragraph';
import Color from 'common/styles/colors';
import { asLeftZeroPaddedString } from 'util/string/NumberString';
import { PrimaryTextField } from 'common/components/TextField/TextField';

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
const volumeGranularity = 100;
const minimumPlaybackRate = 0.5;

const styleClasses = makeStyles({
    toggleButton: {
        width: '100%',
    },
    loadingSpinner: {
        color: Color.Primary,
        display: 'block',
        margin: 'auto',
    },
    volumeContainer: {
        alignItems: 'center',
    },
    volumeIcon: {
        color: Color.Primary,
        display: 'block',
        margin: 'auto',
    },
    podcastImage: {
        width: '100%',
        maxWidth: '540px',
        borderRadius: '10px',
        margin: '25px auto',
        display: 'block',
    },
});

const getTimePartsFromSeconds = (timeInSeconds: number) => {
    const hours = Math.trunc(timeInSeconds / 3600);
    const minutes = Math.trunc(Math.max(timeInSeconds / 60 - hours * 60, 0));
    const seconds = Math.trunc(Math.max(timeInSeconds - (hours * 3600 + minutes * 60), 0));
    return {
        hours: hours,
        minutes: minutes,
        seconds: seconds,
    }
}

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
        const [ audio, setAudio ] = useState<HTMLAudioElement>(new Audio(props.metadata.audioUrl));

        const [ audioCurrentTimeString, setAudioCurrentTimeString ] = useState<string>("00:00:00");
        const [ totalTimeString, setTotalTimeString ] = useState<string>("Loading...");

        const [ isAudioPlaying, setIsAudioPlaying ] = useState<boolean>(false);
        const [ seekValue, setSeekValue ] = useState<number>(0);
        const handleSeekValueChange = (event: Event, newValue: number | number[]) => {
            const val = newValue as number;
            audio.currentTime = val / seekStepGranularity * audio.duration;
            setSeekValue(val);
        };

        const [ volumeValue, setVolumeValue ] = useState<number>(volumeGranularity * audio.volume);
        const handleVolumeValueChange = (event: Event, newValue: number | number[]) => {
            const val = newValue as number;
            audio.volume = val / volumeGranularity;
            setVolumeValue(val);
        };

        const [ playbackSpeed, setPlaybackSpeed ] = useState<number>(1);
        const handlePlaybackSpeedChange = (_: React.ChangeEvent<HTMLSelectElement>, selectedPlaybackSpeed: string) => {
            const val = parseFloat(selectedPlaybackSpeed);
            setPlaybackSpeed(val);
            audio.playbackRate = val;
        }

        const [ isLoading, setIsLoading ] = useState<boolean>(false);
        useEffect(() => {
            audio.addEventListener('loadedmetadata', function() {
                const timeParts = getTimePartsFromSeconds(audio.duration);
                setTotalTimeString(`${asLeftZeroPaddedString(timeParts.hours, 99)}:${asLeftZeroPaddedString(timeParts.minutes, 60)}:${asLeftZeroPaddedString(timeParts.seconds, 60)}`)
            })

            audio.addEventListener("timeupdate", function() {
                setSeekValue(this.currentTime / audio.duration * seekStepGranularity);
                const timeParts = getTimePartsFromSeconds(this.currentTime);
                console.log(timeParts);
                setAudioCurrentTimeString(`${asLeftZeroPaddedString(timeParts.hours, 99)}:${asLeftZeroPaddedString(timeParts.minutes, 60)}:${asLeftZeroPaddedString(timeParts.seconds, 60)}`)
            });
            audio.addEventListener("ended", function() {
                setIsAudioPlaying(false);
            });
            audio.addEventListener("waiting", function() {
                setIsLoading(true);
            });
            audio.addEventListener("playing", function() {
                setIsAudioPlaying(true);
                setIsLoading(false);
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
                        {
                            !!props.metadata.imageUrl && (
                                <Grid item xs={12}>
                                    <img className={classes.podcastImage} src={props.metadata.imageUrl} />
                                </Grid>
                            )
                        }
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
                        <Grid item xs={6}>
                            <Paragraph size={Size.Small} align={Alignment.Left}>
                                {audioCurrentTimeString}
                            </Paragraph>
                        </Grid>
                        <Grid item xs={6}>
                            <Paragraph size={Size.Small} align={Alignment.Right}>
                                {totalTimeString}
                            </Paragraph>
                        </Grid>
                        <Grid item xs={12} md={2}>
                            <CenteredComponent>
                                {
                                    isLoading ? (
                                        <CircularProgress className={classes.loadingSpinner} />
                                    ) : (
                                        <PrimaryButton className={classes.toggleButton} onClick={toggleAudio}>
                                            { isAudioPlaying ? "Pause" : "Play" }
                                        </PrimaryButton>
                                    )
                                }
                            </CenteredComponent>
                        </Grid>
                        <Grid item xs={12}>
                            <CenteredComponent>
                                <Grid className={classes.volumeContainer} container>
                                    <Grid item xs={2} md={2}>
                                        <VolumeUpIcon className={classes.volumeIcon} />
                                    </Grid>
                                    <Grid item xs={10} md={10}>
                                        <PrimarySlider aria-label="Volume" max={volumeGranularity} value={volumeValue} onChange={handleVolumeValueChange} />
                                    </Grid>
                                </Grid>
                            </CenteredComponent>
                        </Grid>
                        <Grid item xs={12}>
                            <CenteredComponent>
                                <Autocomplete
                                    id={`playback-speed-selector`}
                                    onChange={handlePlaybackSpeedChange}
                                    options={Array(20).fill(0).map((_, idx: number) => (idx+1)/10).filter((v: number) => v >= minimumPlaybackRate).map((v: number) => `${v}`)}
                                    value={`${playbackSpeed}`}
                                    renderInput={(params) => <PrimaryTextField label="Playback Speed" {...params} />} />
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
