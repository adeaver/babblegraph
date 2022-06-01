import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Autocomplete from '@material-ui/lab/Autocomplete';
import FormControlLabel from '@material-ui/core/FormControlLabel';

import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton } from 'common/components/Button/Button';
import timezones, { TimeZone } from 'common/data/timezone/timezone';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import { Heading3 } from 'common/typography/Heading';
import Form from 'common/components/Form/Form';
import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';

import {
    GetUserNewsletterScheduleResponse,
    getUserNewsletterSchedule,

    UpdateUserNewsletterScheduleResponse,
    updateUserNewsletterSchedule,
} from 'ConsumerWeb/api/user/preferences';

import { WordsmithLanguageCode } from 'common/model/language/language';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

const styleClasses = makeStyles({
    timeOption: {
        padding: '5px',
    },
    confirmationForm: {
        padding: '10px 0',
        width: '100%',
    },
    emailField: {
        width: '100%',
    },
    submitButtonContainer: {
        alignSelf: 'center',
        padding: '5px',
    },
    submitButton: {
        display: 'block',
        margin: 'auto',
    },
});

const daysOfTheWeek = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];

type TimeSelectorAPIProps = GetUserNewsletterScheduleResponse;

type TimeSelectorOwnProps = {
    subscriptionManagementToken: string;
    languageCode: WordsmithLanguageCode;
    emailAddress?: string;
    omitEmailAddress?: boolean;
    postSubmit?: () => void;
}

type TimeZoneOption = {
    label: string,
    id: string,
}

const TimeSelector = asBaseComponent(
    (props: BaseComponentProps & TimeSelectorAPIProps & TimeSelectorOwnProps) => {
        const [ ianaTimezone, setIANATimezone ] = useState<string>(props.schedule.ianaTimezone);
        const handleTimezoneUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedTimezone: TimeZoneOption) => {
            setIANATimezone(selectedTimezone.id);
        }

        const [ period, setPeriod ] = useState<string>(props.schedule.hourIndex >= 12 ? "PM" : "AM");
        const handlePeriodChange = (_: React.ChangeEvent<HTMLSelectElement>, selectedPeriod: string) => {
            setPeriod(selectedPeriod);
        }

        const [ hourIndex, setHourIndex ] = useState<number>(props.schedule.hourIndex);
        const handleHourIndexChange = (_: React.ChangeEvent<HTMLSelectElement>, selectedHour: number) => {
            setHourIndex(
                selectedHour % 12 + (period === "PM" ? 12 : 0)
            );
        };

        const [ quarterHourIndex, setQuarterHourIndex ] = useState<number>(props.schedule.quarterHourIndex);
        const handleMinuteChange = (_: React.ChangeEvent<HTMLSelectElement>, selectedMinute: number) => {
            setQuarterHourIndex(selectedMinute / 15);
        };

        const [ isActiveForDays, setIsActiveForDays ] = useState<Array<boolean>>(props.schedule.isActiveForDays);

        const [ emailAddress, setEmailAddress ] = useState<string>(props.emailAddress);
        const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setEmailAddress((event.target as HTMLInputElement).value);
        };

        const [ isLoading, setIsLoading ] = useState<boolean>(false);

        const handleSubmit = () => {
            setIsLoading(true);
            updateUserNewsletterSchedule({
                subscriptionManagementToken: props.subscriptionManagementToken,
                languageCode: props.languageCode,
                emailAddress: emailAddress,
                schedule: {
                    ianaTimezone: ianaTimezone,
                    hourIndex: hourIndex,
                    quarterHourIndex: quarterHourIndex,
                    isActiveForDays: isActiveForDays,
                },
                numberOfArticlesPerEmail: 12,
            },
            (resp: UpdateUserNewsletterScheduleResponse) => {
                setIsLoading(false);
                !!props.postSubmit && props.postSubmit();
            },
            (err: Error) => {
                setIsLoading(false);
                props.setError(err);
            });
        }

        const timezoneOptions = timezones.map((t: TimeZone) => ({
            label: t.name,
            id: t.tzCode,
        }));
        const classes = styleClasses()
        return (
            <Grid container>
                <Grid item xs={12}>
                    <Heading3 color={TypographyColor.Primary}>
                        Which days would you like to receive a newsletter?
                    </Heading3>
                </Grid>
                {
                    daysOfTheWeek.map((name: String, idx: number) => (
                        <Grid key={`day-of-week-${idx}`} xs={2} md={3}>
                            <FormControlLabel
                                control={
                                    <PrimaryCheckbox
                                        checked={isActiveForDays[idx]}
                                        onChange={() => {
                                            setIsActiveForDays(
                                                isActiveForDays.map((isChecked: boolean, i: number) => idx === i ? !isChecked : isChecked)
                                            )
                                        }}
                                        name={`day-of-week-${idx}`} />
                                }
                                label={name} />
                        </Grid>
                    ))
                }
                <Grid item xs={12}>
                    <Heading3 color={TypographyColor.Primary}>
                        What time do you want to receive your email?
                    </Heading3>
                </Grid>
                <Grid item xs={4} className={classes.timeOption}>
                    <Autocomplete
                        id="hour-selector"
                        onChange={handleHourIndexChange}
                        options={Array(12).fill(0).map((_: number, idx: number) => idx + 1)}
                        value={!(hourIndex % 12) ? 12 : hourIndex % 12}
                        getOptionLabel={(option: number) => `${option}`}
                        getOptionSelected={(option: number) => option - 1 === hourIndex}
                        renderInput={(params) => <PrimaryTextField label="Hour" {...params} />} />
                </Grid>
                <Grid item xs={4} className={classes.timeOption}>
                    <Autocomplete
                        id="minute-selector"
                        onChange={handleMinuteChange}
                        options={Array(4).fill(0).map((_: number, idx: number) => idx * 15)}
                        value={quarterHourIndex * 15}
                        getOptionLabel={(option: number) => `${option}`}
                        getOptionSelected={(option: number) => option === quarterHourIndex * 15}
                        renderInput={(params) => <PrimaryTextField label="Minute" {...params} />} />
                </Grid>
                <Grid item xs={4} className={classes.timeOption}>
                    <Autocomplete
                        id="period-selector"
                        onChange={handlePeriodChange}
                        options={["AM", "PM"]}
                        value={period}
                        getOptionLabel={(option: string) => option}
                        getOptionSelected={(option: string) => option === period}
                        renderInput={(params) => <PrimaryTextField label="AM/PM" {...params} />} />
                </Grid>
                <Grid item xs={12}>
                    <Heading3 color={TypographyColor.Primary}>
                        What timezone are you in?
                    </Heading3>
                    <Paragraph>
                        Try typing the name of a big city near you, like New York, Los Angeles, or Chicago.
                    </Paragraph>
                    <Paragraph size={Size.Small}>
                        We do not sell or give away this information.
                    </Paragraph>
                    <Autocomplete
                        id="timezone-selector"
                        onChange={handleTimezoneUpdate}
                        options={timezoneOptions}
                        value={timezoneOptions.filter((option: TimeZoneOption) => option.id === ianaTimezone)[0]}
                        getOptionLabel={(option: TimeZoneOption) => option.label}
                        getOptionSelected={(option: TimeZoneOption) => option.id === ianaTimezone}
                        renderInput={(params) => <PrimaryTextField label="Select Timezone" {...params} />} />
                </Grid>
                <Grid item xs={12}>
                    <Form
                        className={classes.confirmationForm}
                        handleSubmit={handleSubmit}>
                        <Grid container>
                        {
                            !!props.omitEmailAddress ? (
                                <Grid item xs={4} md={5}>
                                    &nbsp;
                                </Grid>
                            ) : (
                                <Grid item xs={8} md={10}>
                                    <PrimaryTextField
                                        id="email"
                                        className={classes.emailField}
                                        label="Email Address"
                                        variant="outlined"
                                        onChange={handleEmailAddressChange} />
                                </Grid>
                            )
                        }
                            <Grid item xs={4} md={2} className={classes.submitButtonContainer}>
                                <PrimaryButton
                                    type="submit"
                                    className={classes.submitButton}
                                    disabled={!emailAddress && !props.omitEmailAddress}>
                                    Submit
                                </PrimaryButton>
                            </Grid>
                        </Grid>
                    </Form>
                </Grid>
            </Grid>
        );
    },
    (
        ownProps: TimeSelectorOwnProps,
        onSuccess: (resp: TimeSelectorAPIProps) => void,
        onError: (err: Error) => void,
    ) => {
        getUserNewsletterSchedule({
            subscriptionManagementToken: ownProps.subscriptionManagementToken,
            languageCode: ownProps.languageCode,
        },
        (resp: GetUserNewsletterScheduleResponse) => {
            onSuccess({
                ...resp,
            });
        },
        onError);
    },
    false,
);

export default TimeSelector;
