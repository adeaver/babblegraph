import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';
import { WordsmithLanguageCode } from 'common/model/language/language';

export type Vendor = {
    id: string;
    isActive: boolean;
    name: string;
    websiteUrl: string;
}

export type GetAllVendorsRequest = {}

export type GetAllVendorsResponse = {
    vendors: Array<Vendor>;
}

export function getAllVendors(
    req: GetAllVendorsRequest,
    onSuccess: (resp: GetAllVendorsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetAllVendorsRequest, GetAllVendorsResponse>(
        '/ops/api/advertising/get_all_vendors_1',
        req,
        onSuccess,
        onError,
    );
}

export type InsertVendorRequest = {
    name: string;
    websiteUrl: string;
}

export type InsertVendorResponse = {
    id: string;
}

export function insertVendor(
    req: InsertVendorRequest,
    onSuccess: (resp: InsertVendorResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<InsertVendorRequest, InsertVendorResponse>(
        '/ops/api/advertising/insert_vendor_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateVendorRequest = {
    id: string;
    name: string;
    websiteUrl: string;
    isActive: boolean;
}

export type UpdateVendorResponse = {
    success: boolean;
}

export function updateVendor(
    req: UpdateVendorRequest,
    onSuccess: (resp: UpdateVendorResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateVendorRequest, UpdateVendorResponse>(
        '/ops/api/advertising/update_vendor_1',
        req,
        onSuccess,
        onError,
    );
}

export enum AdvertisementSourceType {
    affiliate = 'affiliate',
}

export type AdvertisementSource = {
    id: string;
    name: string;
    url: string;
    type: AdvertisementSourceType;
    isActive: boolean;
}

export type GetAllSourcesRequest = {}

export type GetAllSourcesResponse = {
    sources: Array<AdvertisementSource>;
}

export function getAllSources(
    req: GetAllSourcesRequest,
    onSuccess: (resp: GetAllSourcesResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetAllSourcesRequest, GetAllSourcesResponse>(
        '/ops/api/advertising/get_all_sources_1',
        req,
        onSuccess,
        onError,
    );
}

export type InsertSourceRequest = {
    name: string;
    websiteUrl: string;
    sourceType: AdvertisementSourceType,
}

export type InsertSourceResponse = {
    id: string;
}

export function insertSource(
    req: InsertSourceRequest,
    onSuccess: (resp: InsertSourceResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<InsertSourceRequest, InsertSourceResponse>(
        '/ops/api/advertising/insert_source_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateSourceRequest = {
    id: string;
    name: string;
    websiteUrl: string;
    sourceType: AdvertisementSourceType,
    isActive: boolean;
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
        '/ops/api/advertising/update_source_1',
        req,
        onSuccess,
        onError,
    );
}

export type Campaign = {
    id: string;
	vendorId: string;
	sourceId: string;
	url: string;
	isActive: boolean;
	name: string;
	shouldApplyToAllUsers: boolean;
	expiresAt: Date | undefined;
    rolloutPercentage: number;
}

export type GetCampaignRequest = {
    id: string;
}

export type GetCampaignResponse = {
    campaign: Campaign;
}

export function getCampaign(
    req: GetCampaignRequest,
    onSuccess: (resp: GetCampaignResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetCampaignRequest, GetCampaignResponse>(
        '/ops/api/advertising/get_campaign_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetAllCampaignsRequest = {}

export type GetAllCampaignsResponse = {
    campaigns: Array<Campaign>;
}

export function getAllCampaigns(
    req: GetAllCampaignsRequest,
    onSuccess: (resp: GetAllCampaignsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetAllCampaignsRequest, GetAllCampaignsResponse>(
        '/ops/api/advertising/get_all_campaigns_1',
        req,
        onSuccess,
        onError,
    );
}

export type InsertCampaignRequest = {
    vendorId: string;
    sourceId: string;
    url: string;
    name: string,
    shouldApplyToAllUsers: boolean;
}

export type InsertCampaignResponse = {
    id: string;
}

export function insertCampaign(
    req: InsertCampaignRequest,
    onSuccess: (resp: InsertCampaignResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<InsertCampaignRequest, InsertCampaignResponse>(
        '/ops/api/advertising/insert_campaign_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateCampaignRequest = {
    campaignId: string;
    vendorId: string;
    sourceId: string;
    url: string;
    name: string,
    isActive: boolean;
    shouldApplyToAllUsers: boolean;
    rolloutPercentage: number;
}

export type UpdateCampaignResponse = {
    success: boolean;
}

export function updateCampaign(
    req: UpdateCampaignRequest,
    onSuccess: (resp: UpdateCampaignResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateCampaignRequest, UpdateCampaignResponse>(
        '/ops/api/advertising/update_campaign_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetCampaignTopicMappingsRequest = {
    campaignId: string;
}

export type GetCampaignTopicMappingsResponse = {
    topicIds: string[];
}

export function getCampaignTopicMappings(
    req: GetCampaignTopicMappingsRequest,
    onSuccess: (resp: GetCampaignTopicMappingsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetCampaignTopicMappingsRequest, GetCampaignTopicMappingsResponse>(
        '/ops/api/advertising/get_campaign_topic_mappings_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateCampaignTopicMappingsRequest = {
    campaignId: string;
    activeTopicMappings: string[];
}

export type UpdateCampaignTopicMappingsResponse = {
    success: boolean;
}

export function updateCampaignTopicMappings(
    req: UpdateCampaignTopicMappingsRequest,
    onSuccess: (resp: UpdateCampaignTopicMappingsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateCampaignTopicMappingsRequest, UpdateCampaignTopicMappingsResponse>(
        '/ops/api/advertising/update_campaign_topic_mappings_1',
        req,
        onSuccess,
        onError,
    );
}

export type Advertisement = {
	id: string;
	languageCode: WordsmithLanguageCode,
	campaignId: string;
    title: string;
	imageUrl: string;
	description: string;
	isActive: boolean;
}

export type CampaignMetrics = {
    lastRefreshedAt: string;
    advertisementMetrics: Array<AdvertisementMetrics>;
}

export type AdvertisementMetrics = {
    advertisementId: string;
    numberOfSends: AdvertisementMetricsNumber;
    numberOfOpenedEmails: AdvertisementMetricsNumber;
    numberOfClicks: AdvertisementMetricsNumber;
}

export type AdvertisementMetricsNumber = {
    total: number;
    unique: number;
}

export type GetAllAdvertisementsForCampaignRequest = {
    campaignId: string;
}

export type GetAllAdvertisementsForCampaignResponse = {
    advertisements: Array<Advertisement>;
    campaignMetrics: CampaignMetrics;
}

export function getAllAdvertisementsForCampaign(
    req: GetAllAdvertisementsForCampaignRequest,
    onSuccess: (resp: GetAllAdvertisementsForCampaignResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetAllAdvertisementsForCampaignRequest, GetAllAdvertisementsForCampaignResponse>(
        '/ops/api/advertising/get_all_advertisements_for_campaign_1',
        req,
        onSuccess,
        onError,
    );
}

export type InsertAdvertisementRequest = {
    languageCode: WordsmithLanguageCode,
    campaignId: string;
    title: string;
    description: string;
    imageUrl: string;
}

export type InsertAdvertisementResponse = {
    id: string;
}

export function insertAdvertisement(
    req: InsertAdvertisementRequest,
    onSuccess: (resp: InsertAdvertisementResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<InsertAdvertisementRequest, InsertAdvertisementResponse>(
        '/ops/api/advertising/insert_advertisement_1',
        req,
        onSuccess,
        onError,
    );
}

export type UpdateAdvertisementRequest = {
    id: string;
    languageCode: WordsmithLanguageCode,
    title: string;
    description: string;
    imageUrl: string;
    isActive: boolean;
}

export type UpdateAdvertisementResponse = {
    success: boolean;
}

export function updateAdvertisement(
    req: UpdateAdvertisementRequest,
    onSuccess: (resp: UpdateAdvertisementResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<UpdateAdvertisementRequest, UpdateAdvertisementResponse>(
        '/ops/api/advertising/update_advertisement_1',
        req,
        onSuccess,
        onError,
    );
}
