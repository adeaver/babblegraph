import { DisplayLanguage } from 'common/model/language/language';

export enum TextBlock {

}

const translations = {

}

export const getTextBlocksForLanguage = (language: DisplayLanguage | undefined) => {
    const displayLanguage = !!language ? language : DisplayLanguage.English;
    return translations[displayLanguage];
}
