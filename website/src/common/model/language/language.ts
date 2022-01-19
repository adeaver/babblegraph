export enum DisplayLanguage {
    Spanish = 'es',
    English = 'en',
}

export enum WordsmithLanguageCode {
    Spanish = 'es',
}

export const getEnglishNameForLanguageCode = (code: WordsmithLanguageCode) => {
    switch (code) {
        case WordsmithLanguageCode.Spanish:
            return "Spanish"
        default:
            throw new Error(`Unrecognized language: ${code}`)
    }
}
