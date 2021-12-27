import React, { useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import ClearIcon from '@material-ui/icons/Clear';
import Autocomplete from '@material-ui/lab/Autocomplete';

import Color from 'common/styles/colors';
import DisplayCard from 'common/components/DisplayCard/DisplayCard';
import { Heading3, Heading4 } from 'common/typography/Heading';
import { Alignment, TypographyColor } from 'common/typography/common';
import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton } from 'common/components/Button/Button';
import Paragraph, { Size } from 'common/typography/Paragraph';

import {
    DayPreferences,
} from 'ConsumerWeb/api/user/schedule';
import {
    contentTopicDisplayMappings,
    ContentTopicDisplayMapping,
} from 'ConsumerWeb/api/user/contentTopics';

const daysOfTheWeekByLanguageCode: { [languageCode: string]: Array<string> } = {
    "es": [ "Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"],
}
const minimumNumberOfArticles = 4;
const maximumNumberOfArticles = 12;
const maxContentTopicsPerDay = 6;

const styleClasses = makeStyles({
    preferencesContainer: {
        padding: '10px',
    },
    alignedContainer: {
        display: 'flex',
        alignItems: 'center',
    },
    removeContentTopicIcon: {
        color: Color.Warning,
    },
    buttonWithMargin: {
        margin: '10px 0',
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
        <div>
            <Heading3 color={TypographyColor.Primary}>
                Customize your newsletter by day
            </Heading3>
            <Paragraph size={Size.Small}>
                Hit the checkmark on any day you do not wish to receive a newsletter
            </Paragraph>
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
                        return;
                    })
                }
            </Grid>
        </div>
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

    const classes = styleClasses();
    return (
        <DisplayCard>
            <Grid className={classes.alignedContainer} container>
                <Grid item xs={1} md={2}>
                    <PrimaryCheckbox
                        checked={isActive}
                        onChange={handleUpdateIsActive}
                        name={`checkbox-${props.dayTitle}`} />

                </Grid>
                <Grid item xs={11} md={10}>
                    <Heading4
                        align={Alignment.Left}
                        color={isActive ? TypographyColor.Primary : TypographyColor.Gray}>
                        {props.dayTitle}
                    </Heading4>
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
            <ContentTopicView
                dayIndex={props.dayPreferences.dayIndex}
                selectedContentTopics={contentTopics}
                handleRemoveContentTopic={handleUpdateContentTopics(true)}
                handleAddContentTopic={handleUpdateContentTopics(false)} />
        </DisplayCard>
    );
}

type ContentTopicViewProps = {
    dayIndex: number;
    selectedContentTopics: string[];

    handleRemoveContentTopic: (topic: string) => void;
    handleAddContentTopic: (topic: string) => void;
}

const ContentTopicView = (props: ContentTopicViewProps) => {
    const selectedContentTopicMappings: Array<ContentTopicDisplayMapping> = !props.selectedContentTopics ? [] : (
        Object.values(
            props.selectedContentTopics
                .map((selectedContentTopicApiValue: string) => {
                    return contentTopicDisplayMappings.filter((mapping: ContentTopicDisplayMapping) => mapping.apiValue.indexOf(selectedContentTopicApiValue) !== -1)
                })
                .filter((mappings: Array<ContentTopicDisplayMapping>) => mappings.length === 1)
                .reduce((acc: { [key: string]: ContentTopicDisplayMapping }, mappings: ContentTopicDisplayMapping[]) => ({
                    ...acc,
                    [mappings[0].displayText]: mappings[0],
                }), {})
        )
    );

    const [ newOption, setNewOption ] = useState<ContentTopicDisplayMapping | null>(null);
    const availableMappings = contentTopicDisplayMappings
        .filter((mapping: ContentTopicDisplayMapping) => {
            return mapping.apiValue.reduce((acc: number, apiValue: string) => {
                const nextValue = (props.selectedContentTopics || []).indexOf(apiValue);
                return nextValue > acc ? nextValue : acc
            }, -1) === -1
        });
    const handleAddTopicChange = (_: React.ChangeEvent<HTMLSelectElement>, selectedOption: ContentTopicDisplayMapping) => {
        setNewOption(selectedOption);
    }
    const submitNewOption = () => {
        setNewOption(null);
        newOption.apiValue.forEach((apiValue: string) => props.handleAddContentTopic(apiValue));
    }

    const classes = styleClasses();
    return (
        <div>
            <Heading4 color={TypographyColor.Primary}>
                Topics on this day
            </Heading4>
            <Paragraph size={Size.Small}>
                Any topics selected will show up in your newsletter for this day if there is content available for it. You can select up to 6 topics.
                If you select fewer than 4 topics, then random topics from your selected interests will be in the newsletter.
            </Paragraph>
            {
                selectedContentTopicMappings.map((mapping: ContentTopicDisplayMapping, idx: number) => (
                    <Grid container
                        className={classes.alignedContainer}>
                        <Grid item xs={10} md={11}>
                            <Paragraph align={Alignment.Left}>
                                { mapping.displayText }
                            </Paragraph>
                        </Grid>
                        <Grid item xs={2} md={1}>
                            <ClearIcon
                                className={classes.removeContentTopicIcon}
                                onClick={() => {
                                    mapping.apiValue.forEach((apiValue: string) => props.handleRemoveContentTopic(apiValue))
                                }}  />
                        </Grid>
                    </Grid>
                ))
            }
            <Paragraph size={Size.Small}>
                Add up to 6 topics you want to appear on this day
            </Paragraph>
            <Autocomplete
                id={`${props.dayIndex}-topic-selector`}
                onChange={handleAddTopicChange}
                options={availableMappings}
                disabled={(props.selectedContentTopics || []).length >= maxContentTopicsPerDay}
                getOptionLabel={(option: ContentTopicDisplayMapping) => option.displayText}
                renderInput={(params) => <PrimaryTextField label="Add a topic" {...params} />} />
            <PrimaryButton
                className={classes.buttonWithMargin}
                disabled={(props.selectedContentTopics || []).length >= maxContentTopicsPerDay}
                onClick={submitNewOption}>
                Add new topic
            </PrimaryButton>
        </div>
    );
}


export default CustomizationByDayTool;
