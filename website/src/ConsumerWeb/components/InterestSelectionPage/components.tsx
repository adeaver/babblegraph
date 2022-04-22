import React, { useState } from 'react';

import Grid from '@material-ui/core/Grid';
import FormControlLabel from '@material-ui/core/FormControlLabel';

import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';
import { toTitleCase } from 'util/string/StringConvert';

import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';
import { WordsmithLanguageCode } from 'common/model/language/language';

import {
    Topic,
    GetTopicsForLanguageResponse,
    getTopicsForLanguage,
} from 'ConsumerWeb/api/content/content';

type TopicSelectorOwnProps = {
    languageCode: WordsmithLanguageCode;

    handleTopicsChange: (topics: Topic[]) => void;
}

export const TopicSelector = asBaseComponent(
    (props: BaseComponentProps & TopicSelectorOwnProps & GetTopicsForLanguageResponse) => {
        const topicsByID = (props.results || []).reduce((acc: { [topicID: string]: Topic }, next: Topic) => ({
            ...acc,
            [next.topicId]: next,
        }), {});

        const [ checkedTopics, setCheckedTopics ] = useState<Topic[]>([]);
        const handleToggleSelected = (event: React.ChangeEvent<HTMLInputElement>) => {
            const topicID = event.target.value as string;
            let nextTopics = checkedTopics;
            if (event.target.checked) {
                nextTopics = checkedTopics.filter((t: Topic) => t.topicId !== topicID).concat(topicsByID[topicID]);
            } else {
                nextTopics = checkedTopics.filter((t: Topic) => t.topicId !== topicID);
            }
            setCheckedTopics(nextTopics);
            props.handleTopicsChange(nextTopics);
        }

        return (
            <Grid container>
                {
                    (props.results || []).map((t: Topic) => (
                        <Grid item xs={4} md={3}>
                            <FormControlLabel
                                control={
                                    <PrimaryCheckbox
                                        value={t.topicId}
                                        onChange={handleToggleSelected}
                                        checked={checkedTopics.some((t2: Topic) => t.topicId === t2.topicId)}
                                        name={`checkbox-${t.topicId}`} />
                                }
                                label={toTitleCase(t.englishLabel)} />
                        </Grid>
                    ))
                }
            </Grid>
        );
    },
    (
        ownProps: TopicSelectorOwnProps,
        onSuccess: (resp: GetTopicsForLanguageResponse) => void,
        onError: (err: Error) => void,
    ) => getTopicsForLanguage({ languageCode: ownProps.languageCode }, onSuccess, onError),
    false,
);
