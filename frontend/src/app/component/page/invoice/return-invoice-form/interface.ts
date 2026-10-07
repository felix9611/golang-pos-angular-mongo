export interface CreateReturnInvoiceItemDto {
    returnInvoiceId: string
    productId: string
    qty: number
    price: number
    taxType: string
    taxCode: string
    taxRate: number
    taxAmount: number
    product: any
}

export interface CreateReturnInvoiceDto {
    returnInvoiceNo: string
    returnReason: string
    returnLocationId: string
    returnDatail: string
    returnMethod: string
    processMethod: string
    refund: boolean
    refundAmount:  number
    refundMethod: string
    returnToVendor: boolean
    returnItems: CreateReturnInvoiceItemDto[]
}

export interface UpdateReturnInvoiceDto extends CreateReturnInvoiceDto {
    _id: string
}