import React, { useState } from 'react';

import {
    DayPreferences,
} from 'ConsumerWeb/api/user/schedule';

type CustomizationByDayToolProps = {
    preferencesByDay: Array<DayPreferences>;
}

type DayIndexToPreferencesMap = { [dayIndex: number]: DayPreferences }

const CustomizationByDayTool = (props: CustomizationByDayToolProps) => {
    const dayPreferencesToIndexMap = props.preferencesByDay.reduce((agg: DayIndexToPreferencesMap, curr: DayPreferences) => ({
        ...agg,
        [curr.dayIndex]: curr,
    }), {});

    console.log(dayPreferencesToIndexMap);

    const [ preferencesByDay, setPreferencesByDay ] = useState<DayIndexToPreferencesMap>(dayPreferencesToIndexMap);

    return (
        <div>
            {
                Array(7).fill(0).map((_, idx: number) => {
                    const dayPreferencesForIdx = dayPreferencesToIndexMap[idx];
                    if (!!dayPreferencesForIdx) {
                        return <DayPreferencesView key={`day-preferences-view-${idx}`} {...dayPreferencesForIdx} />
                    }
                })
            }
        </div>
    );
}

const DayPreferencesView = (props: DayPreferences) => {
    return (
        <h3>{props.dayIndex}</h3>
    );
}

export default CustomizationByDayTool;
