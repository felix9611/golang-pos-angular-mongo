export interface UserInfo {
    _id?: string
    username: string
    accessToken: string
    roles?: any[]
    deptIds?: any[]
    departments?: any[]
    avatarBase64?: string
    email: string
    lastLogin?: string
    roleLists?: any[]
    loginRecords?: any[]
    superAdmin?: boolean
    manager?: boolean
    placeIds?: string[]
    locations?: any[]
}