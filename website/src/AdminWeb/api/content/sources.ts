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

export type SourceSeed = {
    id: string;
    rootId: string;
    url: string;
    isActive: boolean;
}

export type SourceSeedTopicMapping = {
	id: string;
	sourceSeedId: string;
	topicId: string;
	isActive: boolean;
}

export type SourceFilter = {
	id: string;
	rootId: string;
	isActive: boolean;
	useLdJsonValidation: boolean | undefined;
	paywallClasses: string[];
	paywallIds: string[];
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

export type GetAllSourceSeedsForSourceRequest = {
    sourceId: string;
}

export type GetAllSourceSeedsForSourceResponse = {
    sourceSeeds: Array<SourceSeed>;
}

export function getAllSourceSeedsForSource(
    req: GetAllSourceSeedsForSourceRequest,
    onSuccess: (resp: GetAllSourceSeedsForSourceResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetAllSourceSeedsForSourceRequest, GetAllSourceSeedsForSourceResponse>(
        '/ops/api/content/get_all_source_seeds_for_source_1',
        req,
        onSuccess,
        onError,
    );
}

export type AddSourceSeedRequest = {
    sourceId: string;
    url: string;
}

export type AddSourceSeedResponse = {
   id: string;
}

export function addSourceSeed(
    req: AddSourceSeedRequest,
    onSuccess: (resp: AddSourceSeedResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<AddSourceSeedRequest, AddSourceSeedResponse>(
        '/ops/api/content/add_source_seed_for_source_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateSourceSeedRequest = {
    sourceSeedId: string;
    url: string;
    isActive: boolean;
}

export type UpdateSourceSeedResponse = {
   success: boolean;
}

export function updateSourceSeed(
    req: UpdateSourceSeedRequest,
    onSuccess: (resp: UpdateSourceSeedResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateSourceSeedRequest, UpdateSourceSeedResponse>(
        '/ops/api/content/update_source_seed_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetSourceSourceSeedMappingsForSourceRequest = {
    sourceId: string;
}

export type GetSourceSourceSeedMappingsForSourceResponse = {
    sourceSeedMappings: Array<SourceSeedTopicMapping>;
}

export function getSourceSourceSeedMappingsForSource(
    req: GetSourceSourceSeedMappingsForSourceRequest,
    onSuccess: (resp: GetSourceSourceSeedMappingsForSourceResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetSourceSourceSeedMappingsForSourceRequest, GetSourceSourceSeedMappingsForSourceResponse>(
        '/ops/api/content/get_source_seed_mappings_for_source_1',
        req,
        onSuccess,
        onError,
    );
}

export type SourceSeedMappingsUpdate = {
    sourceSeedId: string;
    isActive: boolean;
    topicIds: Array<string>;
}

export type UpsertSourceSeedMappingsRequest = {
    updates: Array<SourceSeedMappingsUpdate>;
}

export type UpsertSourceSeedMappingsResponse = {
    success: boolean;
}

export function upsertSourceSeedMappings(
    req: UpsertSourceSeedMappingsRequest,
    onSuccess: (resp: UpsertSourceSeedMappingsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpsertSourceSeedMappingsRequest, UpsertSourceSeedMappingsResponse>(
        '/ops/api/content/upsert_source_seed_mappings_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetSourceFilterForSourceIDRequest = {
    sourceId: string;
}

export type GetSourceFilterForSourceIDResponse = {
    sourceFilter: SourceFilter | undefined;
}

export function getSourceFilterForSourceID(
    req: GetSourceFilterForSourceIDRequest,
    onSuccess: (resp: GetSourceFilterForSourceIDResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetSourceFilterForSourceIDRequest, GetSourceFilterForSourceIDResponse>(
        '/ops/api/content/get_source_filter_for_source_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpsertSourceFilterForSourceRequest = {
    sourceId: string;
	isActive: boolean;
	useLdJsonValidation: boolean | undefined;
	paywallClasses: string[];
	paywallIds: string[];
}

export type UpsertSourceFilterForSourceResponse = {
	sourceFilterId: string;
}

export function upsertSourceFilterForSource(
    req: UpsertSourceFilterForSourceRequest,
    onSuccess: (resp: UpsertSourceFilterForSourceResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpsertSourceFilterForSourceRequest, UpsertSourceFilterForSourceResponse>(
        '/ops/api/content/upsert_source_filter_for_source_1',
        req,
        onSuccess,
        onError,
    );
}
