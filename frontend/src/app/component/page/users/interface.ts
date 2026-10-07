export interface UserForm {
    _id?: string,
    username: string
    password?: string
    avatarBase64?: string
    deptId?: number
    email: string
    roles?: number[]
}