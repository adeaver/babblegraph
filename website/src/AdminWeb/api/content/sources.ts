import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { WordsmithLanguageCode } from 'common/model/language/language';
import { CountryCode } from 'common/model/geo/geo';

export enum SourceType {
    NewsWebsite = "news-website",
}

export enum IngestStrategy {
    WebsiteHTML1 = "website-html-1",
}

export type Source = {
	id: string;
    title: string;
	url: string;
	type: SourceType;
	country: CountryCode;
	ingestStrategy: IngestStrategy;
	languageCode: WordsmithLanguageCode;
	isActive: boolean;
	monthlyAccessLimit: number | undefined;
    shouldUseUrlAsSeedUrl: boolean;
}

export type GetAllSourcesRequest = {}

export type GetAllSourcesResponse = {
    sources: Array<Source>;
}

export function getAllSources(
    req: GetAllSourcesRequest,
    onSuccess: (resp: GetAllSourcesResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetAllSourcesRequest, GetAllSourcesResponse>(
        '/ops/api/content/get_all_sources_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetSourceByIDRequest = {
   id: string;
}

export type GetSourceByIDResponse = {
    source: Source,
}

export function getSourceByID(
    req: GetSourceByIDRequest,
    onSuccess: (resp: GetSourceByIDResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetSourceByIDRequest, GetSourceByIDResponse>(
        '/ops/api/content/get_source_by_id_1',
        req,
        onSuccess,
        onError,
    );
}

export type AddSourceRequest = {
    title: string | undefined;
	url: string;
	type: SourceType;
	ingestStrategy: IngestStrategy;
	languageCode: WordsmithLanguageCode;
    monthlyAccessLimit: number | undefined;
	country: CountryCode;
    shouldUseUrlAsSeedUrl: boolean;
}

export type AddSourceResponse = {
	id: string;
    title: string;
}

export function addSource(
    req: AddSourceRequest,
    onSuccess: (resp: AddSourceResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<AddSourceRequest, AddSourceResponse>(
        '/ops/api/content/add_source_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateSourceRequest = {
    id: string;
    title: string;
	languageCode: WordsmithLanguageCode;
	url: string;
	type: SourceType;
	ingestStrategy: IngestStrategy;
	isActive: boolean;
    monthlyAccessLimit: number | undefined;
	country: CountryCode;
    shouldUseUrlAsSeedUrl: boolean;
}

export type UpdateSourceResponse = {
    success: boolean;
}

export function updateSource(
    req: UpdateSourceRequest,
    onSuccess: (resp: UpdateSourceResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateSourceRequest, UpdateSourceResponse>(
        '/ops/api/content/update_source_1',
        req,
        onSuccess,
        onError,
    );
}
