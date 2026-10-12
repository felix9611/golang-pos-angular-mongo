export interface MemberClassFrom {
    classCode: string
    className: string
    remark: string
    id?: string
}

export interface ListMemberClassRequestDto {
    page: number
    limit: number
    name?: string
}