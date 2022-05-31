import React, { useState } from 'react';

import Grid from '@material-ui/core/Grid';
import FormControlLabel from '@material-ui/core/FormControlLabel';

import { PrimaryButton } from 'common/components/Button/Button';
import { PrimaryCheckbox } from 'common/components/Checkbox/Checkbox';
import { PrimaryTextField } from 'common/components/TextField/TextField';
import Form from 'common/components/Form/Form';

import {
    TopicWithDisplay,

    getActiveTopicsForLanguageCode,
    GetActiveTopicsForLanguageCodeResponse,
} from 'ConsumerWeb/api/content';
import {
    getUserContentTopicsForToken,
    GetUserContentTopicsForTokenResponse,

    updateUserContentTopicsForToken,
    UpdateUserContentTopicsForTokenResponse,
} from 'ConsumerWeb/api/user/content';

import { WordsmithLanguageCode } from 'common/model/language/language';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';
import { toTitleCase } from 'util/string/StringConvert';

type InterestSelectorOwnProps = {
    languageCode: WordsmithLanguageCode;
    subscriptionManagementToken: string;
    emailAddress?: string;
    omitEmailAddress?: boolean;
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

        const [ emailAddress, setEmailAddress ] = useState<string>(props.emailAddress);
        const handleEmailAddressChange = (event: React.ChangeEvent<HTMLInputElement>) => {
            setEmailAddress((event.target as HTMLInputElement).value);
        };

        const [ isLoading, setIsLoading ] = useState<boolean>(false);
        const [ error, setError ] = useState<Error>(null);

        const handleSubmit = () => {
            setIsLoading(true);
            updateUserContentTopicsForToken({
                activeTopicIds: Object.keys(selectedTopicIDs)
                    .filter((key: keyof SelectedTopicIDsMap) => selectedTopicIDs[key])
                    .map((key: keyof SelectedTopicIDsMap) => key.toString()),
                emailAddress: emailAddress,
                subscriptionManagementToken: props.subscriptionManagementToken,
            },
            (resp: UpdateUserContentTopicsForTokenResponse) => {
                setIsLoading(false);
            },
            (err: Error) => {
                setIsLoading(false);
                setError(err);
            });
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
                <Grid item xs={12}>
                    <Form handleSubmit={handleSubmit}>
                        <Grid container>
                        {
                            !!props.omitEmailAddress ? (
                                <Grid item xs={4} md={5}>
                                    &nbsp;
                                </Grid>
                            ) : (
                                <Grid item xs={8} md={10}>
                                    <PrimaryTextField
                                        id="email"
                                        label="Email Address"
                                        variant="outlined"
                                        onChange={handleEmailAddressChange} />
                                </Grid>
                            )
                        }
                            <Grid item xs={4} md={2}>
                                <PrimaryButton
                                    type="submit"
                                    disabled={!emailAddress && !props.omitEmailAddress}>
                                    Submit
                                </PrimaryButton>
                            </Grid>
                        </Grid>
                    </Form>
                </Grid>
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
            getUserContentTopicsForToken({
                subscriptionManagementToken: ownProps.subscriptionManagementToken,
            },
            (resp2: GetUserContentTopicsForTokenResponse) => {
                onSuccess({
                    allTopics: resp.topics || [],
                    selectedTopicIDs: resp2.topics || [],
                });
            },
            onError);
        },
        onError);
    },
    false,
);

export default InterestSelector;
