export interface MemberClassFrom {
    classCode: string
    className: string
    remark: string
    _id?: string
}

export interface ListMemberClassRequestDto {
    page: number
    limit: number
    name?: string
}