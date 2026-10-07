export interface DashboardReqFilterDto {
    typeIds?: string[]
    placeIds?: string[]
    deptIds?: string[]
    dateRange?: string[]
}

export interface DashboardReqDto {
    // dateTypeValue?: 'YearMonth' | 'YearQuarter' | 'none'
    dataTypeValue?: 'dept' | 'type' | 'location' | 'year-month' | 'stockIn' | 'stockOut' | 'stockMove'
}