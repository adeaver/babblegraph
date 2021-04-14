import { Survey } from 'common/data/surveys/typedefs';

import LowOpen1Survey from 'common/data/surveys/LowOpen1';
import HighOpen1Survey from 'common/data/surveys/HighOpen1';
import Unsubscribe1Survey from 'common/data/surveys/Unsubscribe1';

export enum SurveyKey {
    LowOpen1 = 'LowOpen1',
    HighOpen1 = 'HighOpen1',
    Unsubscribe1 = 'Unsubscribe1',
}

const AvailableSurveys: { [key: string]: Survey } = {
    [SurveyKey.LowOpen1]: LowOpen1Survey,
    [SurveyKey.HighOpen1]: HighOpen1Survey,
    [SurveyKey.Unsubscribe1]: Unsubscribe1Survey,
}

export const getSurveyForKey = (key: SurveyKey) => {
    return AvailableSurveys[key];
}
