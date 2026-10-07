export interface CreateProduct {
    
}

export interface ProductFormDto extends CreateProduct {
    _id: string
    productCode: string
    productName: string
    itemCode: string
    brandCode: string
    brandName: string
    typeId: string
    deptId: string
    vendorId: string
    unit: string
    costPrice: number
    retailPrice: number
    description: string
    remark: string
}

export interface StockMoveProductLocationDto {
    productCode: string
    productId: string
    fromPlaceCode: string
    fromLocationId: string
    toPlaceCode: string
    toLocationId: string
    qty: number
    totalPrice: number
    totalCost: number
}

export interface StockInOutProductLocationDto {
    productCode: string
    productId: string
    placeCode: string
    locationId: string
    qty: number
    totalPrice: number
    totalCost: number
}

export interface ListProductLocationtRequestDto {
    page: number
    limit: number
    locationIds?: string[]
    assetCode?: string
}

export interface ListInvRecordDto {
    page: number
    limit: number,
    dateRange?: string[]
}
