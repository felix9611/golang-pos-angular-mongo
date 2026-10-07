export interface BudgetForm {
    _id?: string,
    deptId: string
    placeId: string
    budgetNo?: string
    budgetName: string
    year: string
    month: string
    budgetAmount: number
    budgetFrom: string
    budgetTo: string
    budgetStatus: string
    remark?: string
}