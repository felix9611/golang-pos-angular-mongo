export interface StockTakeItemDto {
    _id?: string
    stockTakeId: string
    productId: string
    productCode: string
    placeId: string
    status: string
    checkTime: string
    qty: number
    remark?: string
}

export interface StockTakeItemFromDto {
    stockTakeId: string
    productId: string
    productCode: string
    productName: string
    placeId: string
    status: string
    orgQty?: number
    checkQty?: number
    remark?: string
}

export interface StockTakeFormEdit {
    actionName: string
    actionPlaceId: string
    remark?: string
    stockTakeItems: any
    createdTime: string
    finishTime?: string
    _id: string
    createBy: string
    status: number
}