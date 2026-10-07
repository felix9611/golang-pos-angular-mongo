export interface TaxInfomationForm {
    _id?: string,
    nationCode: string
    nationName: string
    countryCode: string
    countryName: string
    taxType: string
    taxCode: string
    taxName: string
    taxRate: number
    importRate?: number
    remark?: string
}