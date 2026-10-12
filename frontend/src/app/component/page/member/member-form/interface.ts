export interface MemberSpecialDay {
    id?: string
    memberId?: string
    name: string
    date: string
    remark: string
}

export interface MemberForm {
    name: string
    address: string
    phone: string
    email: string
    fax: string
    classId: string
    remark: string
    memberSpecialDays: MemberSpecialDay[]
    id?: string
    memberCode?: string
}