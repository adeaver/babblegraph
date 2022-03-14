import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

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
