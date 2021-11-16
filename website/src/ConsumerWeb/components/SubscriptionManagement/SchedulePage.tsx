import React, { useEffect, useState } from 'react';
import { RouteComponentProps, useHistory } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import ClearIcon from '@material-ui/icons/Clear';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import Divider from '@material-ui/core/Divider';
import Snackbar from '@material-ui/core/Snackbar';

import Alert from 'common/components/Alert/Alert';
import Color from 'common/styles/colors';
import Page from 'common/components/Page/Page';
import { Heading1 } from 'common/typography/Heading';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';
import timezones, { TimeZone } from 'common/data/timezone/timezone.ts';

import { ContentHeader } from './common';

import {
    ContentTopicDisplayMapping,
    contentTopicDisplayMappings,
} from 'ConsumerWeb/api/user/contentTopics';
import {
    getUserProfile,
    GetUserProfileResponse
} from 'ConsumerWeb/api/useraccounts/useraccounts';
import {
    getUserNewsletterSchedule,
    GetUserNewsletterScheduleResponse,
    ScheduleByLanguageCode,
    ScheduleDay,

    ScheduleDayRequest,
    AddUserNewsletterScheduleResponse,
    addUserNewsletterSchedule,
} from 'ConsumerWeb/api/user/schedule';

type Params = {
    token: string;
}

const styleClasses = makeStyles({
    daySelectorContainer: {
        minHeight: '100%',
        padding: '10px',
    },
    timezoneSelector: {
        maxWidth: '100%',
    },
    numberOfArticlesSelector: {
        margin: '5px',
    },
    headerSubContent: {
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
    },
    savePreferencesButton: {
        margin: '10px 0',
    },
    /* Pretty sure this is about as hacky as it gets */
    removeContentTopicButtonContainer: {
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'flex-start',
        justifyContent: 'center',
    },
    removeContentTopicIcon: {
        color: Color.Warning,
    },
    contentTopicSelect: {
        minWidth: '100%',
        marginBottom: '10px',
    },
});

const makeDefaultScheduleDay = (dayIndex: number) => ({
    dayOfWeekIndex: dayIndex,
    hourOfDayIndex: 0,
    quarterHourIndex: 0,
    contentTopics: [],
    numberOfArticles: defaultNumberOfArticles,
    isActive: true,
});

const defaultNumberOfArticles = 12;
const minimumNumberOfArticles = 4;
const maximumNumberOfArticles = 12;

const maxContentTopicsPerDay = 6;

const dayNameForLanguageCode = {
    "es": ["Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sabádo"],
}

type ContentTopicDisplayNamesByAPIValue = { [apiValue: string]: string }
const contentTopicNamesForDisplayValue = contentTopicDisplayMappings.reduce((acc: ContentTopicDisplayNamesByAPIValue, contentTopicDisplayMapping: ContentTopicDisplayMapping) => ({
    ...acc,
    ...contentTopicDisplayMapping.apiValue.reduce((acc: ContentTopicDisplayNamesByAPIValue, apiValue: string) => ({
        ...acc,
        [apiValue]: contentTopicDisplayMapping.displayText,
    }), {}),
}), {});

type SchedulePageProps = RouteComponentProps<Params>

type ScheduleDaysByLanguageCode = { [languageCode: string]: ScheduleDayByDayIndex }
type ScheduleDayByDayIndex = { [dayIndex: number]: ScheduleDay }

const SchedulePage = (props: SchedulePageProps) => {
    const { token } = props.match.params;

    const [ subscriptionLevel, setSubscriptionLevel ] = useState<string | null>(null);
    const [ isLoadingUserProfile, setIsLoadingUserProfile ] = useState<boolean>(true);

    const [ userScheduleByLanguageCode, setUserScheduleByLanguageCode ] = useState<ScheduleDaysByLanguageCode>({});
    const [ isLoadingSchedule, setIsLoadingSchedule ] = useState<boolean>(true);

    const [ wasUpdateNewsletterScheduleSuccessful, setWasUpdateNewsletterScheduleSuccessful ] = useState<boolean | null>(null);
    const [ isLoadingNewsletterUpdate, setIsLoadingNewsletterUpdate ] = useState<boolean>(false);

    const [ ianaTimezone, setIANATimezone ] = useState<string>(Intl.DateTimeFormat().resolvedOptions().timeZone || "US/Eastern");

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
    const handleSubmit = (languageCode: string) => {
        setIsLoadingNewsletterUpdate(true);
        addUserNewsletterSchedule({
            ianaTimezone: ianaTimezone,
            languageCode: languageCode,
            userScheduleDayRequests: [...Array(7).keys()].map((dayIndex: number) => {
                let scheduleDay = makeDefaultScheduleDay(dayIndex);
                if (userScheduleByLanguageCode[languageCode] && userScheduleByLanguageCode[languageCode][dayIndex]) {
                    scheduleDay = userScheduleByLanguageCode[languageCode][dayIndex];
                }
                return {
                    dayOfWeekIndex: dayIndex,
                    contentTopics: scheduleDay.contentTopics,
                    numberOfArticles: scheduleDay.numberOfArticles,
                    isActive: scheduleDay.isActive,
                }
            }),
        },
        (resp: AddUserNewsletterScheduleResponse) => {
            setIsLoadingNewsletterUpdate(false);
            setWasUpdateNewsletterScheduleSuccessful(resp.success);
        },
        (e: Error) => {
            setIsLoadingNewsletterUpdate(false);
            setError(e);
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
                ianaTimezone: ianaTimezone,
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

    const closeSnackbar = () => {
        setWasUpdateNewsletterScheduleSuccessful(null);
        setError(null);
    }

    const isLoading = isLoadingUserProfile || isLoadingSchedule || isLoadingNewsletterUpdate;
    return (
        <Page>
            <Grid container>
                <Grid item xs={false} md={3}>
                    &nbsp;
                </Grid>
                <Grid item xs={12} md={6}>
                    <DisplayCard>
                        <ContentHeader
                            title="Newsletter Schedule and Customization"
                            token={token} />
                        {
                            isLoading ? (
                                <LoadingSpinner />
                            ) : (
                                <SchedulePreferencesView
                                    ianaTimezone={ianaTimezone}
                                    subscriptionLevel={subscriptionLevel}
                                    userScheduleByLanguageCode={userScheduleByLanguageCode}
                                    handleSubmit={handleSubmit}
                                    updateIANATimezone={setIANATimezone}
                                    updateUserScheduleDay={updateUserScheduleDay} />
                            )
                        }
                        <Snackbar open={(!wasUpdateNewsletterScheduleSuccessful && wasUpdateNewsletterScheduleSuccessful != null) || !!error} autoHideDuration={6000} onClose={closeSnackbar}>
                            <Alert severity="error">Something went wrong processing your request.</Alert>
                        </Snackbar>
                        <Snackbar open={wasUpdateNewsletterScheduleSuccessful} autoHideDuration={6000} onClose={closeSnackbar}>
                            <Alert severity="success">Successfully updated your email preferences.</Alert>
                        </Snackbar>
                    </DisplayCard>
                </Grid>
            </Grid>
        </Page>
    );
}

type SchedulePreferencesViewProps = {
    ianaTimezone: string;
    subscriptionLevel: string | null;
    userScheduleByLanguageCode: ScheduleDaysByLanguageCode;

    handleSubmit: (languageCode: string) => void;
    updateIANATimezone: (newIANATimezone: string) => void;
    updateUserScheduleDay: (languageCode: string, dayIndex: number, scheduleDay: ScheduleDay) => void;
}

const SchedulePreferencesView = (props: SchedulePreferencesViewProps) => {
    const classes = styleClasses();
    const currentTimezone = timezones.filter((t: TimeZone) => t.tzCode === props.ianaTimezone)[0].name || props.ianaTimezone.replace("_", " ").split("/")[1];
    return (
        <div>
            {
                !props.subscriptionLevel && (
                    <Paragraph color={TypographyColor.Warning}>
                        This feature is only available for premium accounts. You do not have a premium account, so you will not be able to use this.
                    </Paragraph>
                )
            }
            <div className={classes.headerSubContent}>
                <Paragraph>
                    You can update the schedule on which you receive your newsletter, as well as customizing the content you receive in each newsletter here. When you’re done updating your preferences, click the button below to save them.
                </Paragraph>
                <Paragraph size={Size.Small}>
                    Your timezone is currently set as {currentTimezone}
                </Paragraph>
                <FormControl className={classes.timezoneSelector}>
                    <InputLabel id="timezone-selector-label">Change timezone</InputLabel>
                    <Select
                        labelId="timezone-selector-label"
                        id="timezone-selector"
                        value={props.ianaTimezone}
                        onChange={(e) => { props.updateIANATimezone(e.target.value) }}>
                        {
                            timezones.map((t: TimeZone, idx: number) => (
                                <MenuItem key={`timezone-selector-${idx}`} value={t.tzCode}>{t.name}</MenuItem>
                            ))
                        }
                    </Select>
                </FormControl>
                <PrimaryButton className={classes.savePreferencesButton} onClick={() => props.handleSubmit("es")}>
                    Save your preferences
                </PrimaryButton>
            </div>
            <Divider />
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
    const classes = styleClasses();
    const scheduleDay  = props.scheduleForDay ? props.scheduleForDay : makeDefaultScheduleDay(props.dayIndex);

    const [ contentTopicToAdd, setContentTopicToAdd ] = useState<string>("");

    const handleUpdateNumberOfArticles = (event: React.ChangeEvent<HTMLInputElement>) => {
        const numberOfArticles = parseInt((event.target as HTMLInputElement).value, 10);
        if (numberOfArticles >= minimumNumberOfArticles && numberOfArticles <= maximumNumberOfArticles) {
            props.updateUserScheduleDay(props.languageCode, props.dayIndex, {
                ...scheduleDay,
                numberOfArticles: numberOfArticles,
            });
        }
    }
    const handleUpdateIsActive = () => {
        props.updateUserScheduleDay(props.languageCode, props.dayIndex, {
            ...scheduleDay,
            isActive: !scheduleDay.isActive,
        });
    }

    const contentTopics = (scheduleDay.contentTopics || [])
    const handleRemoveContentTopic = (contentTopic: string) => {
        props.updateUserScheduleDay(props.languageCode, props.dayIndex, {
            ...scheduleDay,
            contentTopics: contentTopics.filter((ct: string) => ct !== contentTopic),
        });
    }
    const handleAddContentTopic = () => {
        if (!contentTopics.some((c: string) => c === contentTopicToAdd)) {
            props.updateUserScheduleDay(props.languageCode, props.dayIndex, {
                ...scheduleDay,
                contentTopics: contentTopics.concat(contentTopicToAdd),
            });
        }
    }
    const currentContentTopics = contentTopics.map((contentTopic: string, idx: number) => (
        <Grid container key={`${props.languageCode}-${props.dayIndex}-${idx}-ct`}>
            <Grid item xs={9}>
                <Paragraph>
                    {contentTopicNamesForDisplayValue[contentTopic]}
                </Paragraph>
            </Grid>
            <Grid item className={classes.removeContentTopicButtonContainer} xs={3}>
                <ClearIcon
                    className={classes.removeContentTopicIcon}
                    onClick={() => handleRemoveContentTopic(contentTopic)}  />
            </Grid>
        </Grid>
    ));
    return (
        <Grid item className={classes.daySelectorContainer} xs={12} md={6} xl={4}>
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
                                className={classes.numberOfArticlesSelector}
                                defaultValue={scheduleDay.numberOfArticles}
                                label="Number of Articles"
                                variant="outlined"
                                type="number"
                                error={scheduleDay.numberOfArticles < minimumNumberOfArticles || scheduleDay.numberOfArticles > maximumNumberOfArticles}
                                helperText={`Must select between ${minimumNumberOfArticles} and ${maximumNumberOfArticles}`}
                                onChange={handleUpdateNumberOfArticles} />
                            <Paragraph>
                                Topics for this day
                            </Paragraph>
                            <Paragraph size={Size.Small}>
                                Any topics selected will show up in your newsletter for this day if there is content available for it. You can select up to 6 topics.
                                If you select fewer than 4 topics, then random topics from your selected interests will be in the email.
                            </Paragraph>
                            { currentContentTopics }
                            <FormControl className={classes.contentTopicSelect}>
                                <InputLabel id={`${props.languageCode}-${props.dayIndex}-cts-selector-label`}>Add topic</InputLabel>
                                <Select
                                    labelId={`${props.languageCode}-${props.dayIndex}-cts-selector-label`}
                                    id={`${props.languageCode}-${props.dayIndex}-cts-selector`}
                                    value={contentTopicToAdd}
                                    disabled={contentTopics.length >= maxContentTopicsPerDay}
                                    helperText="Select up to 6 topics"
                                    onChange={(e) => { setContentTopicToAdd(e.target.value) }}>
                                    {
                                        contentTopicDisplayMappings
                                            .map((displayMapping: ContentTopicDisplayMapping) => (
                                                <MenuItem key={`${props.languageCode}-${props.dayIndex}-${displayMapping.apiValue[0]}-item`} value={displayMapping.apiValue[0]}>{displayMapping.displayText}</MenuItem>
                                            ))
                                    }
                                </Select>
                            </FormControl>
                            <PrimaryButton onClick={handleAddContentTopic} disabled={!contentTopicToAdd}>
                                Add Topic
                            </PrimaryButton>
                        </div>
                    )
                }
            </DisplayCard>
        </Grid>
    )
}

export default SchedulePage;
