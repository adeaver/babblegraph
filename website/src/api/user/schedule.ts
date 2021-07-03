import { makePostRequestWithStandardEncoding } from 'api/bgfetch/bgfetch';

export type AddUserNewsletterScheduleRequest = {
    userScheduleDayRequests: Array<ScheduleDayRequest>;
    languageCode: string;
    ianaTimezone: string;
}

export type ScheduleDayRequest = {
    dayOfWeekIndex: number;
    contentTopics: string[];
    numberOfArticles: number;
    isActive: boolean;
}

export type AddUserNewsletterScheduleResponse = {
    success: boolean;
}

export function addUserNewsletterSchedule(
    req: AddUserNewsletterScheduleRequest,
    onSuccess: (resp: AddUserNewsletterScheduleResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<AddUserNewsletterScheduleRequest, AddUserNewsletterScheduleResponse>(
        '/api/user/add_user_schedule_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetUserNewsletterScheduleRequest = {
    ianaTimezone: string;
}

export type GetUserNewsletterScheduleResponse = {
    scheduleByLanguageCode: Array<ScheduleByLanguageCode>;
}

export type ScheduleByLanguageCode = {
    languageCode: string;
    scheduleDays: Array<ScheduleDay>;
}

export type ScheduleDay = {
    dayOfWeekIndex: number;
    hourOfDayIndex: number;
    quarterHourIndex: number;
    contentTopics: string[];
    numberOfArticles: number;
    isActive: boolean;
}

export function getUserNewsletterSchedule(
    req: GetUserNewsletterScheduleRequest,
    onSuccess: (resp:  GetUserNewsletterScheduleResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserNewsletterScheduleRequest, GetUserNewsletterScheduleResponse>(
        '/api/user/get_user_schedule_1',
        req,
        onSuccess,
        onError,
    );
}
