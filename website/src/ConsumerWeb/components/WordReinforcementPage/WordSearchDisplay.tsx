import React, { useState } from 'react';

import Grid from '@material-ui/core/Grid';

import { withCaptchaToken } from 'common/util/grecaptcha/grecaptcha';

import {
    IDType,
    LanguageLookupID,
    SearchResult,
    SearchTextResult,
    SearchTextResponse,
    searchText
} from 'ConsumerWeb/api/language/search';

import { WordsmithLanguageCode } from 'common/model/language/language';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

type WordSearchDisplayOwnProps = {
    searchTerms: string[];
    wordReinforcementToken: string;
    subscriptionManagementToken: string;
}

type WordSearchDisplayAPIProps = SearchTextResponse;

const WordSearchDisplay = asBaseComponent(
    (props: BaseComponentProps & WordSearchDisplayOwnProps & WordSearchDisplayAPIProps) => {
        return (
            <Grid container>
                Hello
            </Grid>
        );
    },
    (
        ownProps: WordSearchDisplayOwnProps,
        onSuccess: (resp: WordSearchDisplayAPIProps) => void,
        onError: (err: Error) => void,
    ) => {
        withCaptchaToken("searchtext", (token: string) => {
            searchText({
                wordReinforcementToken: ownProps.wordReinforcementToken,
                languageCode: WordsmithLanguageCode.Spanish,
                text: ownProps.searchTerms,
                captchaToken: token,
            },
            (resp: SearchTextResponse) => onSuccess(resp),
            onError)
        });
    },
    false
);

export default WordSearchDisplay;
