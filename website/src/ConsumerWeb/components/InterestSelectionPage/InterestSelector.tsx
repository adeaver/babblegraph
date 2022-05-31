import React, { useState } from 'react';

import Grid from '@material-ui/core/Grid';

import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';
import FormControlLabel from '@material-ui/core/FormControlLabel';

import {
    TopicWithDisplay,

    getActiveTopicsForLanguageCode,
    GetActiveTopicsForLanguageCodeResponse,
} from 'ConsumerWeb/api/content';

import { WordsmithLanguageCode } from 'common/model/language/language';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';
import { toTitleCase } from 'util/string/StringConvert';

type InterestSelectorOwnProps = {
    languageCode: WordsmithLanguageCode;
}

type InterestSelectorAPIProps = {
    allTopics: Array<TopicWithDisplay>;
    selectedTopicIDs: Array<string>;
}

type SelectedTopicIDsMap = { [topicID: string]: boolean }

const InterestSelector = asBaseComponent(
    (props: InterestSelectorOwnProps & BaseComponentProps & InterestSelectorAPIProps) => {
        const [ selectedTopicIDs, setSelectedTopicIDs ] = useState<SelectedTopicIDsMap>(
            props.selectedTopicIDs.reduce((acc: SelectedTopicIDsMap, next: string) => ({
                ...acc,
                [next]: true,
            }), {})
        );
        const handleTopicIDChange = (topicID: string, isSelected: boolean) => {
            setSelectedTopicIDs({
                ...selectedTopicIDs,
                [topicID]: isSelected,
            })
        }

        return (
            <Grid container>
                {
                    props.allTopics.map((t: TopicWithDisplay) => (
                        <Grid key={`topic-${t.topic.id}`} item xs={4}>
                            <FormControlLabel
                                control={
                                    <PrimaryCheckbox
                                        checked={!!selectedTopicIDs[t.topic.id]}
                                        onChange={() => { handleTopicIDChange(t.topic.id, !selectedTopicIDs[t.topic.id]) }}
                                        name={`topic-${t.topic.id}`} />
                                }
                                label={toTitleCase(t.topic.label.replace("current-events-", " ").replace(/\-/g, " "))} />
                        </Grid>
                    ))
                }
            </Grid>
        );
    },
    (
        ownProps: InterestSelectorOwnProps,
        onSuccess: (resp: InterestSelectorAPIProps) => void,
        onError: (err: Error) => void,
    ) => {
        getActiveTopicsForLanguageCode({
            languageCode: ownProps.languageCode,
        },
        (resp: GetActiveTopicsForLanguageCodeResponse) => {
            // TODO: add user topics
            onSuccess({
                allTopics: resp.topics || [],
                selectedTopicIDs: [],
            });
        },
        onError);
    },
    false,
);

export default InterestSelector;
