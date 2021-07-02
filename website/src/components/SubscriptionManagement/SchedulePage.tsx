import React, { useEffect, useState } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import Grid from '@material-ui/core/Grid';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import ClearIcon from '@material-ui/icons/Clear';

import Page from 'common/components/Page/Page';
import { Heading1 } from 'common/typography/Heading';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';

import { ContentHeader } from './common';

import {
    getUserProfile,
    GetUserProfileResponse
} from 'api/useraccounts/useraccounts';
import {
    getUserNewsletterSchedule,
    GetUserNewsletterScheduleResponse,
    ScheduleByLanguageCode,
    ScheduleDay,

    ScheduleDayRequest,
    AddUserNewsletterScheduleResponse,
    addUserNewsletterSchedule,
} from 'api/user/schedule';

type Params = {
    token: string;
}

const defaultNumberOfArticles = 12;
const dayNameForLanguageCode = {
    "es": ["Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sabádo"],
}

type SchedulePageProps = RouteComponentProps<Params>

type ScheduleDaysByLanguageCode = { [languageCode: string]: ScheduleDayByDayIndex }
type ScheduleDayByDayIndex = { [dayIndex: number]: ScheduleDay }

const SchedulePage = (props: SchedulePageProps) => {
    const { token } = props.match.params;

    const [ subscriptionLevel, setSubscriptionLevel ] = useState<string | null>(null);
    const [ isLoadingUserProfile, setIsLoadingUserProfile ] = useState<boolean>(true);

    const [ userScheduleByLanguageCode, setUserScheduleByLanguageCode ] = useState<ScheduleDaysByLanguageCode>({});
    const [ isLoadingSchedule, setIsLoadingSchedule ] = useState<boolean>(true);

    const [ error, setError ] = useState<Error>(null);

    const updateUserScheduleDay = (languageCode: string, dayIndex: number, scheduleDay: ScheduleDay) => {
        setUserScheduleByLanguageCode({
            ...userScheduleByLanguageCode,
            [languageCode]: {
                ...userScheduleByLanguageCode[languageCode],
                [dayIndex]: scheduleDay,
            }
        });
    }

    useEffect(() => {
        getUserProfile({
            subscriptionManagementToken: token,
        },
        (resp: GetUserProfileResponse) => {
            setIsLoadingUserProfile(false);
            if (resp.subscriptionLevel) {
                setSubscriptionLevel(resp.subscriptionLevel);
                getUserNewsletterSchedule({
                    // TODO: replace this
                    ianaTimezone: "US/Eastern",
                },
                (resp: GetUserNewsletterScheduleResponse) => {
                    setIsLoadingSchedule(false);
                    const scheduleByLanguageCodes = resp.scheduleByLanguageCode || [];
                    setUserScheduleByLanguageCode(
                        scheduleByLanguageCodes
                            .reduce((acc: ScheduleDaysByLanguageCode, byLanguageCode: ScheduleByLanguageCode) => ({
                                ...acc,
                                [byLanguageCode.languageCode]:  byLanguageCode.scheduleDays.reduce((acc: ScheduleDayByDayIndex, scheduleDay: ScheduleDay) => ({
                                    ...acc,
                                    [scheduleDay.dayOfWeekIndex]: scheduleDay,
                                }), {})
                            }), userScheduleByLanguageCode)
                    );
                },
                (e: Error) => {
                    setIsLoadingSchedule(false);
                    setError(e);
                });
            } else {
                setIsLoadingSchedule(false);
            }
        },
        (e: Error) => {
            setIsLoadingUserProfile(false);
            setIsLoadingSchedule(false);
            setError(e);
        });
    }, []);

    const isLoading = isLoadingUserProfile && isLoadingSchedule;
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <ContentHeader
                            title="Newsletter Schedule and Preferences"
                            token={token} />
                        {
                            isLoading ? (
                                <LoadingSpinner />
                            ) : (
                                <SchedulePreferencesView
                                    subscriptionLevel={subscriptionLevel}
                                    userScheduleByLanguageCode={userScheduleByLanguageCode}
                                    updateUserScheduleDay={updateUserScheduleDay} />
                            )
                        }
                    </DisplayCard>
                </Grid>
            </Grid>
        </Page>
    );
}

type SchedulePreferencesViewProps = {
    subscriptionLevel: string | null;
    userScheduleByLanguageCode: ScheduleDaysByLanguageCode;

    updateUserScheduleDay: (languageCode: string, dayIndex: number, scheduleDay: ScheduleDay) => void;
}

const SchedulePreferencesView = (props: SchedulePreferencesViewProps) => {
    return (
        <div>
            {
                !props.subscriptionLevel && (
                    <Paragraph color={TypographyColor.Warning}>
                        This feature is only available for premium accounts. You do not have a premium account, so you will not be able to use this.
                    </Paragraph>
                )
            }
            <SchedulePreferencesForm
                languageCode="es"
                schedule={props.userScheduleByLanguageCode["es"]}
                updateUserScheduleDay={props.updateUserScheduleDay} />
        </div>
    );
}

type SchedulePreferencesFormProps = {
    languageCode: string;
    schedule: ScheduleDayByDayIndex,

    updateUserScheduleDay: (languageCode: string, dayIndex: number, scheduleDay: ScheduleDay) => void;
}

const SchedulePreferencesForm = (props: SchedulePreferencesFormProps) => {
    const dayIndices = [...Array(7).keys()];
    return (
        <Grid container>
            {
                dayIndices.map((val: number) => {
                    return (
                        <ScheduleSelector
                            key={`${props.languageCode}-${val}-selector`}
                            languageCode={props.languageCode}
                            dayIndex={val}
                            scheduleForDay={props.schedule ? props.schedule[val] : undefined}
                            updateUserScheduleDay={props.updateUserScheduleDay} />
                    )
                })
            }
        </Grid>
    );
}

type ScheduleDaySelectorProps = {
    languageCode: string;
    dayIndex: number;
    scheduleForDay: ScheduleDay | undefined;

    updateUserScheduleDay: (languageCode: string, dayIndex: number, scheduleDay: ScheduleDay) => void;
}

const ScheduleSelector = (props: ScheduleDaySelectorProps) => {
    const scheduleDay  = props.scheduleForDay ? props.scheduleForDay : {
        dayOfWeekIndex: props.dayIndex,
        hourOfDayIndex: 0,
        quarterHourIndex: 0,
        contentTopics: ["art", "something"],
        numberOfArticles: defaultNumberOfArticles,
        isActive: true,
    }

    const handleUpdateNumberOfArticles = (event: React.ChangeEvent<HTMLInputElement>) => {
        props.updateUserScheduleDay(props.languageCode, props.dayIndex, {
            ...scheduleDay,
            numberOfArticles: parseInt((event.target as HTMLInputElement).value, 10),
        });
    }
    const handleUpdateIsActive = () => {
        props.updateUserScheduleDay(props.languageCode, props.dayIndex, {
            ...scheduleDay,
            isActive: !scheduleDay.isActive,
        });
    }

    const handleRemoveContentTopic = (contentTopic: string) => {
        props.updateUserScheduleDay(props.languageCode, props.dayIndex, {
            ...scheduleDay,
            contentTopics: scheduleDay.contentTopics.filter((ct: string) => ct !== contentTopic),
        });
    }
    const handleAddContentTopic = (contentTopic: string) => {
        props.updateUserScheduleDay(props.languageCode, props.dayIndex, {
            ...scheduleDay,
            contentTopics: scheduleDay.contentTopics.concat(contentTopic),
        });
    }
    const currentContentTopics = scheduleDay.contentTopics.map((contentTopic: string, idx: number) => (
        <div key={`${props.languageCode}-${props.dayIndex}-${idx}-ct`}>
            <Paragraph>
                {contentTopic}
            </Paragraph>
            <ClearIcon onClick={handleRemoveContentTopic(contentTopic)} />
        </div>
    ));
    console.log(currentContentTopics);
    return (
        <Grid item xs={6} md={3}>
            <DisplayCard>
                <Paragraph color={scheduleDay.isActive ? TypographyColor.Primary : TypographyColor.Gray}>
                    { dayNameForLanguageCode[props.languageCode][props.dayIndex] }
                </Paragraph>
                <FormControlLabel
                    control={
                        <PrimaryCheckbox
                            checked={scheduleDay.isActive}
                            onChange={handleUpdateIsActive}
                            name={"checkbox-is-active"} />
                    }
                    label="Receive an email this day?" />
                {
                    scheduleDay.isActive && (
                        <div>
                            <PrimaryTextField
                                id="number_of_articles"
                                defaultValue={scheduleDay.numberOfArticles}
                                label="Number of Articles"
                                variant="outlined"
                                type="number"
                                onChange={handleUpdateNumberOfArticles} />
                            { currentContentTopics }
                        </div>
                    )
                }
            </DisplayCard>
        </Grid>
    )
}

export default SchedulePage;
