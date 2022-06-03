import React, { useState, useEffect } from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Form from 'common/components/Form/Form';
import { loadCaptchaScript } from 'common/util/grecaptcha/grecaptcha';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import { PrimaryButton } from 'common/components/Button/Button';
import LoadingSpinner from 'common/components/LoadingSpinner/LoadingSpinner';

import { WordsmithLanguageCode } from 'common/model/language/language';
import { UserVocabularyEntry } from 'ConsumerWeb/api/user/userVocabulary';

import WordSearchDisplay from "./WordSearchDisplay";

const styleClasses = makeStyles({
    searchTextInput: {
        width: '100%',
    },
    submitButton: {
        margin: '5px 0',
        width: '100%',
    },
});

type WordSearchFormProps = {
    wordReinforcementToken: string;
    subscriptionManagementToken: string;
    userVocabularyEntries: Array<UserVocabularyEntry>;

    handleAddNewUserVocabularyEntry: (newEntry: UserVocabularyEntry) => void;
}

const WordSearchForm = (props: WordSearchFormProps) => {
    const [ hasLoadedCaptcha, setHasLoadedCaptcha ] = useState<boolean>(false);

    useEffect(() => {
        loadCaptchaScript();
        setHasLoadedCaptcha(true);
    }, []);

    const [ searchFormValue, setSearchFormValue ] = useState<string>(null);
    const handleSearchFormValueChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setSearchFormValue((event.target as HTMLInputElement).value);
    };


    const [ searchTerms, setSearchTerms ] = useState<string[]>(null);
    const handleSubmit = () => {
        const terms = searchFormValue.trim().split(/ +/g);
        setSearchTerms(terms);
    }

    const classes = styleClasses();
    if (!hasLoadedCaptcha) {
        return <LoadingSpinner />;
    }
    return (
        <Form handleSubmit={handleSubmit}>
            <Grid container>
                <Grid item xs={12}>
                    <PrimaryTextField
                        className={classes.searchTextInput}
                        id="search-form-value"
                        label="Search for a word or phrase"
                        defaultValue={searchFormValue}
                        variant="outlined"
                        onChange={handleSearchFormValueChange} />
                </Grid>
                <Grid item xs={12} md={6}>
                    <PrimaryButton
                        className={classes.submitButton}
                        disabled={!searchFormValue}
                        type="submit">
                        Search
                    </PrimaryButton>
                </Grid>
                {
                    !!searchTerms && (
                        <Grid key={`terms-${searchTerms.join('-')}`} item xs={12}>
                            <WordSearchDisplay
                                searchTerms={searchTerms}
                                userVocabularyEntries={props.userVocabularyEntries}
                                subscriptionManagementToken={props.subscriptionManagementToken}
                                wordReinforcementToken={props.wordReinforcementToken}
                                handleAddNewUserVocabularyEntry={props.handleAddNewUserVocabularyEntry} />
                        </Grid>
                    )
                }
            </Grid>
        </Form>
    )
}

export default WordSearchForm;
