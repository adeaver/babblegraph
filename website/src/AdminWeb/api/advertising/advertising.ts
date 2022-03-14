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
