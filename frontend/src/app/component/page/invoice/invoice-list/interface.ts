export interface InvoiceListRequestDto {
    page: number
    limit: number
    number?: string
    dateRange?: string[]
}