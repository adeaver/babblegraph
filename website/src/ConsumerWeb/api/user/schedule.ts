import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export type DayPreferences = {
    isActive: boolean;
    numberOfArticles: number;
    contentTopics: string[];
    dayIndex: number;
}

export type GetUserScheduleRequest = {
    token: string;
    languageCode: string;
}

export type GetUserScheduleResponse = {
    userIanaTimezone: string;
    hourIndex: number;
    quarterHourIndex: number;
    preferencesByDay: Array<DayPreferences>;
}

export function getUserSchedule(
    req: GetUserScheduleRequest,
    onSuccess: (resp: GetUserScheduleResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUserScheduleRequest, GetUserScheduleResponse>(
        '/api/user/get_user_newsletter_schedule_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateUserScheduleRequest = {
    emailAddress: string;
    token: string;
    languageCode: string;
    hourIndex: number;
    quarterHourIndex: number;
    ianaTimezone: string;
}

export type UpdateUserScheduleResponse = {
    success: boolean;
}

export function updateUserSchedule(
    req: UpdateUserScheduleRequest,
    onSuccess: (resp: UpdateUserScheduleResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateUserScheduleRequest, UpdateUserScheduleResponse>(
        '/api/user/update_user_newsletter_schedule_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateUserScheduleWithDayPreferencesRequest = {
    emailAddress: string;
    token: string;
    languageCode: string;
    hourIndex: number;
    quarterHourIndex: number;
    ianaTimezone: string;
    dayPreferences: Array<DayPreferences>;
}

export type UpdateUserScheduleWithDayPreferencesResponse = {
    success: boolean;
}

export function updateUserScheduleWithDayPreferences(
    req: UpdateUserScheduleWithDayPreferencesRequest,
    onSuccess: (resp: UpdateUserScheduleWithDayPreferencesResponse) => void,
    onError: (err: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateUserScheduleWithDayPreferencesRequest, UpdateUserScheduleWithDayPreferencesResponse>(
        '/api/user/update_user_newsletter_schedule_and_day_preferences_1',
        req,
        onSuccess,
        onError,
    );
}
