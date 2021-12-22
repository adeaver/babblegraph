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
const minimumNumberOfArticles = 4;
const maximumNumberOfArticles = 12;

const styleClasses = makeStyles({
    preferencesContainer: {
        padding: '10px',
    },
});

type CustomizationByDayToolProps = {
    languageCode: string;
    preferencesByDay: Array<DayPreferences>;

    handleUpdatePreferencesByDay: (d: Array<DayPreferences>) => void;
}

type DayIndexToPreferencesMap = { [dayIndex: number]: DayPreferences }

const CustomizationByDayTool = (props: CustomizationByDayToolProps) => {
    const dayPreferencesToIndexMap = props.preferencesByDay.reduce((agg: DayIndexToPreferencesMap, curr: DayPreferences) => ({
        ...agg,
        [curr.dayIndex]: curr,
    }), {});

    const [ preferencesByDay, setPreferencesByDay ] = useState<DayIndexToPreferencesMap>(dayPreferencesToIndexMap);

    const daysOfTheWeekForLanguageCode = daysOfTheWeekByLanguageCode[props.languageCode];
    const handleUpdateDayPreferences = (d: DayPreferences) => {
        const nextObj = {
            ...preferencesByDay,
            [d.dayIndex]: d,
        };
        setPreferencesByDay(nextObj);
        props.handleUpdatePreferencesByDay(Object.values(nextObj));
    }

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
                                    handleUpdateDayPreferences={handleUpdateDayPreferences}
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
    dayPreferences: DayPreferences;

    handleUpdateDayPreferences: (d: DayPreferences) => void;
}

const DayPreferencesView = (props: DayPreferencesViewProps) => {

    const [ isActive, setIsActive ] = useState<boolean>(props.dayPreferences.isActive);
    const handleUpdateIsActive = () => {
        const nextValue = !isActive;
        props.handleUpdateDayPreferences({
            ...props.dayPreferences,
            isActive: nextValue,
        });
        setIsActive(nextValue);
    };

    const [ numberOfArticles, setNumberOfArticles ] = useState<number>(props.dayPreferences.numberOfArticles);
    const handleUpdateNumberOfArticles = (event: React.ChangeEvent<HTMLInputElement>) => {
        const numberOfArticles = parseInt((event.target as HTMLInputElement).value, 10);
        setNumberOfArticles(numberOfArticles);
        if (numberOfArticles >= minimumNumberOfArticles && numberOfArticles <= maximumNumberOfArticles) {
            props.handleUpdateDayPreferences({
                ...props.dayPreferences,
                numberOfArticles: numberOfArticles,
            });
        }
    }

    const [ contentTopics, setContentTopics ] = useState<string[]>(props.dayPreferences.contentTopics);
    const handleUpdateContentTopics = (remove: boolean) => {
        return (contentTopic: string) => {
            const nextContentTopics = remove ? (
                contentTopics.filter((topic: string) => topic !== contentTopic)
            ) : (
                !contentTopics.some((topic: string) => topic === contentTopic) && contentTopics.concat(contentTopic)
            );
            setContentTopics(nextContentTopics);
            props.handleUpdateDayPreferences({
                ...props.dayPreferences,
                contentTopics: contentTopics,
            });
        }
    }

    return (
        <DisplayCard>
            <Grid container>
                <Grid item xs={1} md={2}>
                    <PrimaryCheckbox
                        checked={isActive}
                        onChange={handleUpdateIsActive}
                        name={`checkbox-${props.dayTitle}`} />

                </Grid>
                <Grid item xs={11} md={10}>
                    <Heading3
                        align={Alignment.Left}
                        color={isActive ? TypographyColor.Primary : TypographyColor.Gray}>
                        {props.dayTitle}
                    </Heading3>
                </Grid>
            </Grid>
            <PrimaryTextField
                id="number-of-articles"
                value={numberOfArticles}
                type="number"
                label="Number of Articles per Email"
                variant="outlined"
                error={numberOfArticles < minimumNumberOfArticles || numberOfArticles > maximumNumberOfArticles}
                helperText={`Must select between ${minimumNumberOfArticles} and ${maximumNumberOfArticles}`}
                onChange={handleUpdateNumberOfArticles} />
            {
                /* TODO: this
                    <ContentTopicView /> */
            }
        </DisplayCard>
    );
}


export default CustomizationByDayTool;
