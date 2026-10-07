export interface UserInfo {
    _id?: string
    username: string
    accessToken: string
    roles?: number[],
    deptId?: number,
    avatarBase64?: string
    email: string
    lastLogin?: string
    roleLists?: any[]
    loginRecords?: any[]
}