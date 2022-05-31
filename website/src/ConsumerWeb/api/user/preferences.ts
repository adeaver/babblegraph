import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { ClientError } from 'ConsumerWeb/api/clienterror';
import { WordsmithLanguageCode } from 'common/model/language/language';

export type Schedule = {
    ianaTimezone: string;
    hourIndex: number;
    quarterHourIndex: number;
    isActiveForDays: Array<boolean>;
}

export type GetUserNewsletterScheduleRequest = {
    subscriptionManagementToken: string;
    languageCode: WordsmithLanguageCode;
}

export type GetUserNewsletterScheduleResponse = {
    schedule: Schedule | undefined;
    numberOfArticlesPerEmail: string;
    error: ClientError | undefined;
}

export function getUserNewsletterSchedule(
    req: GetUserNewsletterScheduleRequest,
    onSuccess: (resp: GetUserNewsletterScheduleResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserNewsletterScheduleRequest, GetUserNewsletterScheduleResponse>(
        '/api/user/get_user_newsletter_schedule_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateUserNewsletterScheduleRequest = {
    subscriptionManagementToken: string;
    languageCode: WordsmithLanguageCode;
    emailAddress: string | undefined;
    schedule: Schedule;
    numberOfArticlesPerEmail: number;
}

export type UpdateUserNewsletterScheduleResponse = {
    error: ClientError | undefined;
    success: boolean;
}

export function updateUserNewsletterSchedule(
    req: UpdateUserNewsletterScheduleRequest,
    onSuccess: (resp: UpdateUserNewsletterScheduleResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateUserNewsletterScheduleRequest, UpdateUserNewsletterScheduleResponse>(
        '/api/user/update_user_newsletter_schedule_1',
        req,
        onSuccess,
        onError,
    );
}
