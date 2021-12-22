import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { Heading3 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';
import { PrimaryTextField } from 'common/components/TextField/TextField';

import {
    DayPreferences,
} from 'ConsumerWeb/api/user/schedule';

const daysOfTheWeekByLanguageCode: { [languageCode: string]: Array<string> } = {
    "es": [ "Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"],
}

const styleClasses = makeStyles({
    preferencesContainer: {
        padding: '10px',
    },
});

type CustomizationByDayToolProps = {
    languageCode: string;
    preferencesByDay: Array<DayPreferences>;
}

type DayIndexToPreferencesMap = { [dayIndex: number]: DayPreferences }

const CustomizationByDayTool = (props: CustomizationByDayToolProps) => {
    const dayPreferencesToIndexMap = props.preferencesByDay.reduce((agg: DayIndexToPreferencesMap, curr: DayPreferences) => ({
        ...agg,
        [curr.dayIndex]: curr,
    }), {});

    const [ preferencesByDay, setPreferencesByDay ] = useState<DayIndexToPreferencesMap>(dayPreferencesToIndexMap);

    const daysOfTheWeekForLanguageCode = daysOfTheWeekByLanguageCode[props.languageCode];

    const classes = styleClasses();
    return (
        <Grid container>
            {
                Array(7).fill(0).map((_, idx: number) => {
                    const dayPreferencesForIdx = dayPreferencesToIndexMap[idx];
                    if (!!dayPreferencesForIdx) {
                        return (
                            <Grid className={classes.preferencesContainer} item xs={12} md={6}>
                                <DayPreferencesView key={`day-preferences-view-${idx}`}
                                    dayTitle={daysOfTheWeekForLanguageCode[idx]}
                                    dayPreferences={dayPreferencesForIdx} />
                            </Grid>
                        );
                    }
                })
            }
        </Grid>
    );
}

type DayPreferencesViewProps = {
    dayTitle: string;
    dayPreferences: DayPreferences,
}

const DayPreferencesView = (props: DayPreferencesViewProps) => {
    return (
        <DisplayCard>
            <Grid container>
                <Grid item xs={1} md={2}>
                    <PrimaryCheckbox
                        checked={props.dayPreferences.isActive}
                        onChange={}
                        name={`checkbox-${props.dayTitle}`} />

                </Grid>
                <Grid item xs={11} md={10}>
                    <Heading3
                        align={Alignment.Left}
                        color={props.dayPreferences.isActive ? TypographyColor.Primary : TypographyColor.Gray}>
                        {props.dayTitle}
                    </Heading3>
                </Grid>
            </Grid>
            <PrimaryTextField
                id="email"
                value={props.dayPreferences.numberOfArticles}
                type="number"
                label="Number of Articles per Email (4 to 12)"
                variant="outlined"
                onChange={handleEmailAddressChange} />
        </DisplayCard>
    );
}

export default CustomizationByDayTool;
