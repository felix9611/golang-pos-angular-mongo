export interface ListMemberRequestDto {
    page: number
    limit: number
    name?: string
    contact?: string
    classIds?: string[]
}