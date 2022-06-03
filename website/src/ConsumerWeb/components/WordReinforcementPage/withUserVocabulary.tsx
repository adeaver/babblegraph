import React, { useState } from 'react';

import {
    UserVocabularyEntry,

    GetUserVocabularyResponse,
    getUserVocabulary,
} from 'ConsumerWeb/api/user/userVocabulary';

import { WordsmithLanguageCode } from 'common/model/language/language';
import {
    asBaseComponent,
    BaseComponentProps,
} from 'common/base/BaseComponent';

type WithUserVocabularyProps = {
    subscriptionManagementToken: string;
}

type UserVocabularyComponentAPIProps = {
    userVocabularyEntries: Array<UserVocabularyEntry>;
}

export type InjectedUserVocabularyComponentProps = {
    userVocabularyEntries: Array<UserVocabularyEntry>;

    handleAddNewVocabularyEntry: (newEntry: UserVocabularyEntry) => void;
} & BaseComponentProps;

export function withUserVocabulary<P extends WithUserVocabularyProps>(WrappedComponent: React.ComponentType<P & InjectedUserVocabularyComponentProps>) {
    return asBaseComponent(
        (props: P & UserVocabularyComponentAPIProps & BaseComponentProps) => {
            const [ userVocabularyEntries, setUserVocabularyEntries ] = useState<Array<UserVocabularyEntry>>(props.userVocabularyEntries);
            const handleAddNewVocabularyEntry = (newEntry: UserVocabularyEntry) => {
                setUserVocabularyEntries(userVocabularyEntries.concat(newEntry));
            }

            const wrappedProps = {
                ...props,
                "userVocabularyEntries": userVocabularyEntries,
            }
            console.log("wrapped", wrappedProps);
            return (
                <WrappedComponent
                    {...wrappedProps }
                    handleAddNewVocabularyEntry={handleAddNewVocabularyEntry} />
            )
        },
        (
            ownProps: P,
            onSuccess: (resp: UserVocabularyComponentAPIProps) => void,
            onError: (err: Error) => void
        ) => {
            getUserVocabulary({
                subscriptionManagementToken: ownProps.subscriptionManagementToken,
                languageCode: WordsmithLanguageCode.Spanish,
            },
            (resp: GetUserVocabularyResponse) => {
                onSuccess({
                    userVocabularyEntries: resp.entries
                });
            },
            onError);
        },
        false,
    );
}
