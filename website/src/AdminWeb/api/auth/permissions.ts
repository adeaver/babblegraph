import { makePostRequestWithStandardEncoding } from 'util/bgfetch/bgfetch';

export enum Permission {
    ManagePermissions = 'manage-permissions',
    ViewUserMetrics = 'view-user-metrics',
	WriteBlog = "write-blog",
	PublishBlog = "publish-blog",
}

export type ManageUserPermissionsRequest = {
	adminId: string;
    updates: Array<PermissionUpdate>;
}

export type PermissionUpdate = {
    permission: Permission;
    isActive: boolean;
}

export type ManageUserPermissionResponse = {
	success: boolean;
}

export function manageUserPermissions(
    req: ManageUserPermissionsRequest,
    onSuccess: (resp: ManageUserPermissionResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<ManageUserPermissionsRequest, ManageUserPermissionResponse>(
        '/ops/api/auth/manage_user_permissions_1',
        req,
        onSuccess,
        onError,
    );
}

export type GetUsersWithPermissionsRequest = {};

export type GetUsersWithPermissionsResponse = {
	users: Array<UserWithPermissions>;
}

export type UserWithPermissions = {
	id: string;
	emailAddress: string;
	permissions:  Array<Permission>;
}

export function getUsersWithPermissions(
    req: GetUsersWithPermissionsRequest,
    onSuccess: (resp: GetUsersWithPermissionsResponse) => void,
    onError: (e: Error) => void,
) {
    makePostRequestWithStandardEncoding<GetUsersWithPermissionsRequest, GetUsersWithPermissionsResponse>(
        '/ops/api/auth/get_users_with_permissions_1',
        req,
        onSuccess,
        onError,
    );
}
