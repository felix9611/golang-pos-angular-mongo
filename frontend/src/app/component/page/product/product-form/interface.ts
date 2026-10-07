export interface ProductFileDto {
    fileName: string
    base64: string
    fileType: string
}

export interface CreateProduct {
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
    taxType: string
    uploaProductFiles?: ProductFileDto[]
    productListFiles?: any[]
}

export interface ProductFormDto extends CreateProduct {
    _id?: string
    productCode?: string
}
