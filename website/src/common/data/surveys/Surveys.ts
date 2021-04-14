import { Survey } from 'common/data/surveys/typedefs';

import LowOpen1Survey from 'common/data/surveys/LowOpen1';

export enum SurveyKey {
    LowOpen1 = 'LowOpen1',
}

const AvailableSurveys: { [key: string]: Survey } = {
    [SurveyKey.LowOpen1]: LowOpen1Survey,
}

export const getSurveyForKey = (key: SurveyKey) => {
    return AvailableSurveys[key];
}
