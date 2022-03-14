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
