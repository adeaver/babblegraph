import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { ClientError } from 'ConsumerWeb/api/clienterror';
import { WordsmithLanguageCode } from 'common/model/language/language';

export type UpsertUserVocabularyRequest = {
    subscriptionManagementToken: string;
    languageCode: WordsmithLanguageCode;
    vocabularyEntry: VocabularyEntry;
}

export type VocabularyEntry = {
    displayText: string;
    definitionId: string | undefined;
    entryType: UserVocabularyType;
    studyNote: string | null;
    isVisible: boolean;
    isActive: boolean;
}

export enum UserVocabularyType {
    Lemma = 'lemma',
    Phrase = 'phrase'
}

export type UpsertUserVocabularyResponse = {
    id: string | undefined;
    // TODO: new error type
    error: ClientError | undefined;
}

export function upsertUserVocabulary(
    req: UpsertUserVocabularyRequest,
    onSuccess: (resp: UpsertUserVocabularyResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpsertUserVocabularyRequest, UpsertUserVocabularyResponse>(
        '/api/user/upsert_user_vocabulary_entry_1',
        req,
        onSuccess,
        onError,
    );
}

