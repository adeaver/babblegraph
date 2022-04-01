import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import Autocomplete from '@material-ui/lab/Autocomplete';

import { Heading3 } from 'common/typography/Heading';
import Paragraph, { Size } from 'common/typography/Paragraph';
import { Alignment, TypographyColor } from 'common/typography/common';
import timezones, { TimeZone } from 'common/data/timezone/timezone';
import { PrimaryTextField } from 'common/components/TextField/TextField';

const styleClasses = makeStyles({
    timezoneSelector: {
        maxWidth: '100%',
        minWidth: '100%',
        marginTop: '10px',
    },
    timeSelector: {
        maxWidth: '100%',
        minWidth: '100%',
    },
    timeSelectorContainer: {
        padding: '5px',
    },
});

type TimeSelectorProps = {
    initialIANATimezone: string;
    ianaTimezone: string;
    hourIndex: number;
    quarterHourIndex: number;

    handleUpdateIANATimezone: (tz: string) => void;
    handleUpdateHourIndex: (h: number) => void;
    handleUpdateQuarterHourIndex: (q: number) => void;
}

type TimeZoneOption = {
    label: string,
    id: string,
}

const TimeSelector = (props: TimeSelectorProps) => {
    const classes = styleClasses();
    const currentTimezoneOptions = timezones.filter((t: TimeZone) => t.tzCode === props.ianaTimezone);
    const initialTimezoneOptions = timezones.filter((t: TimeZone) => t.tzCode === props.initialIANATimezone);

    const initialTimezone = initialTimezoneOptions.length ? initialTimezoneOptions[0].name : "Unknown";
    const [ currentTimezone, setCurrentTimezone ] = useState<string>(currentTimezoneOptions.length ? currentTimezoneOptions[0].name : props.ianaTimezone.replace("_", " ").split("/")[1]);

    const [ period, setPeriod ] = useState<string>(props.hourIndex >= 12 ? "PM" : "AM");
    const [ displayHour, setDisplayHour ] = useState<number>(
        props.hourIndex === 12 || props.hourIndex === 0 ? 12 : (
            props.hourIndex - (period === "PM" ? 12 : 0)
        )
    );

    const handleTimezoneUpdate = (_: React.ChangeEvent<HTMLSelectElement>, selectedTimezone: TimeZoneOption) => {
        const foundTimezone = timezones.filter((t: TimeZone) => t.tzCode === selectedTimezone.id)[0];
        if (foundTimezone) {
            props.handleUpdateIANATimezone(foundTimezone.tzCode);
            setCurrentTimezone(foundTimezone.name);
        }
    }
    const handleHourUpdate = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const hourInTwelveHourFormat = parseInt(e.target.value, 10);
        setDisplayHour(hourInTwelveHourFormat);
        if (hourInTwelveHourFormat === 12) {
            props.handleUpdateHourIndex(period === "PM" ? 12 : 0);
            return;
        }
        props.handleUpdateHourIndex(hourInTwelveHourFormat + (period === "PM" ? 12 : 0));
    }
    const handleQuarterHourUpdate = (e: React.ChangeEvent<HTMLSelectElement>) => {
        props.handleUpdateQuarterHourIndex(parseInt(e.target.value, 10))
    }
    const handlePeriodUpdate = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const proposedPeriod = e.target.value;
        if (period === proposedPeriod) {
            return;
        }
        const updatedHourIndexForProposedPeriod = proposedPeriod === "PM" ? (
            (props.hourIndex + 12) % 24
        ) : (
            Math.abs(props.hourIndex - 12)
        );
        debugger;
        props.handleUpdateHourIndex(updatedHourIndexForProposedPeriod);
        setPeriod(proposedPeriod);
    }
    return (
        <div>
            <Heading3 color={TypographyColor.Primary}>
                What time do you want to receive your newsletter?
            </Heading3>
            <Grid container>
                <Grid item xs={4} className={classes.timeSelectorContainer}>
                    <FormControl className={classes.timeSelector}>
                        <InputLabel id="hour-selector-label">Select Hour</InputLabel>
                        <Select
                            labelId="hour-selector-label"
                            id="hour-selector"
                            value={displayHour}
                            onChange={handleHourUpdate}>
                            {
                                Array(12).fill(0).map((_, idx: number) => (
                                    <MenuItem key={`hour-selector-${idx}`} value={idx+1}>{idx+1}</MenuItem>
                                ))
                            }
                        </Select>
                    </FormControl>
                </Grid>
                <Grid item xs={4} className={classes.timeSelectorContainer}>
                    <FormControl className={classes.timeSelector}>
                        <InputLabel id="quarter-hour-selector-label">Select Minute</InputLabel>
                        <Select
                            labelId="quarter-hour-selector-label"
                            id="quarter-hour-selector"
                            value={props.quarterHourIndex}
                            onChange={handleQuarterHourUpdate}>
                            {
                                Array(4).fill(0).map((_, idx: number) => (
                                    <MenuItem key={`quarter-hour-selector-${idx}`} value={idx*15}>{idx*15}</MenuItem>
                                ))
                            }
                        </Select>
                    </FormControl>
                </Grid>
                <Grid item xs={4} className={classes.timeSelectorContainer}>
                    <FormControl className={classes.timeSelector}>
                        <InputLabel id="period-selector-label">Select Period (AM/PM)</InputLabel>
                        <Select
                            labelId="period-selector-label"
                            id="period-selector"
                            value={period}
                            onChange={handlePeriodUpdate}>
                            {
                                ["AM", "PM"].map((period: string, idx: number) => (
                                    <MenuItem key={`period-selector-${idx}`} value={period}>{period}</MenuItem>
                                ))
                            }
                        </Select>
                    </FormControl>
                </Grid>
            </Grid>
            <Paragraph size={Size.Small}>
                Your timezone is currently set as {initialTimezone}
            </Paragraph>
            <Autocomplete
                id="timezone-selector"
                onChange={handleTimezoneUpdate}
                options={timezones.map((t: TimeZone) => ({
                    label: t.name,
                    id: t.tzCode,
                }))}
                getOptionLabel={(option: TimeZoneOption) => option.label}
                getOptionSelected={(option: TimeZoneOption) => option.label === currentTimezone}
                renderInput={(params) => <PrimaryTextField label="Select Timezone" {...params} />} />
        </div>
    );
}

export default TimeSelector;
