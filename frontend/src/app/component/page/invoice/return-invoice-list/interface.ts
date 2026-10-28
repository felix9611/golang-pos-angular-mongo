export interface InvoiceListRequestDto {
    page: number
    limit: number
    dateRange?: string[]
    invoiceNumber: string
    caseNumber?: string
    locationId?: string[]
}