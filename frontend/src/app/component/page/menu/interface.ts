export interface MenuForm {
    _id?: string,
    mainId: string
    name: string
    icon: string
    path: string
    sort: number
    type: string
    excelFunctionCode: string
    excelFunctionName: string
    menuIds: string[]
}