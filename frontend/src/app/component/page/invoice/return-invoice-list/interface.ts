export interface ReturnInvoiceListRequestDto {
    page: number
    limit: number
    dateRange: string[]
    invoiceNumber: string
    caseNumber: string
    returnLocationIds: string[]
}