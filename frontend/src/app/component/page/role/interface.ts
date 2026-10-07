export interface RoleForm {
    _id?: string,
    name: string
    code: string
    remark: string
    menuIds: string[]
    read: boolean
    write: boolean
    delete: boolean
    update: boolean
    upload: boolean
}