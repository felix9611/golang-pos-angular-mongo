export interface UserInfo {
    _id?: string
    username: string
    accessToken: string
    roles?: any[],
    deptIds?: any[],
    department?: any[],
    avatarBase64?: string
    email: string
    lastLogin?: string
    roleLists?: any[]
    loginRecords?: any[],
    superAdmin?: boolean
    placeIds?: string[],
    locations?: any[]
}