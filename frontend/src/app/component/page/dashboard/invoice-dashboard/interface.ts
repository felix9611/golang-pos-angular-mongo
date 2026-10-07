export interface DashboardReqFilterDto {
    typeIds?: string[]
    placeIds?: string[]
    deptIds?: string[]
    salesDateRange?: string[]
}

export interface DashboardReqDto {
    // dateTypeValue?: 'YearMonth' | 'YearQuarter' | 'none'
    dataTypeValue?: 'dept' | 'type' | 'location' | 'year-month'
    valueField: 'qtys' | 'price'
}