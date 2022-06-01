import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export enum OnboardingStatus {
	NotStarted = 'not-started',
	Settings = 'settings',
	Schedule = 'schedule',
	InterestSelection = 'interest-selection',
	Vocabulary = 'vocabulary',
	Finished = 'finished',
}

export type GetOnboardingStatusForUserRequest = {
    onboardingToken: string;
}

export type GetOnboardingStatusForUserResponse = {
    emailAddress: string | undefined;
    onboardingStatus: OnboardingStatus;
    subscriptionManagementToken: string | undefined;
}

export function getOnboardingStatusForUser(
    req: GetOnboardingStatusForUserRequest,
    onSuccess: (resp: GetOnboardingStatusForUserResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetOnboardingStatusForUserRequest, GetOnboardingStatusForUserResponse>(
        '/api/onboarding/get_onboarding_status_for_user_1',
        req,
        onSuccess,
        onError,
    );
}
