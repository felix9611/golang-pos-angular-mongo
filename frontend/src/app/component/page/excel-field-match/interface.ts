export interface ExcelFieldListForm {
    _id?: string,
    functionCode: string
    functionName: string
    functionType: string
    fieldLists: ExcelFieldList[]
}

export interface ExcelFieldList {
    dbFieldName: string
    excelFieldName: string
    sort: number
}