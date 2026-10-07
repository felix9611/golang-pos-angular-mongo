export interface InvoiceItemDto {
    productId: string
    productCode: string
    qty: number
    price: number
    discount: number
    discountType: string
    taxType: string
    taxCode: string
    taxRate: number
    taxAmount: number 
}

export interface InvoicePaymentDto {
    method: string
    amount: number
    paymentTime?: string
    balance: number
    findRedemption: number
}

export interface CreateInvoiceForm {
    memberId: string
    locationId: string
    locationCode: string
    totalAmount: number
    discount: number
    discountType: string
    taxTotal: number
    remark: string
    invoiceItems: InvoiceItemDto[]
    invoicePayments: InvoicePaymentDto[]
}

export interface MmeberForm {
    _id: string
    memberCode: string
    name: string
    phone: string
    email: string
}
