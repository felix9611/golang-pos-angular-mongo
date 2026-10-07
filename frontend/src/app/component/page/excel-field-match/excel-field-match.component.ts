import { Component } from '@angular/core'
import { CommonModule } from '@angular/common'
import { RouterLink, RouterLinkActive, RouterOutlet, Router } from '@angular/router'
import { FormsModule } from '@angular/forms'
import { deleteApiWithAuth, getApiWithAuth, postApiWithAuth } from '../../../../tool/httpRequest-auth'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzModalModule, NzModalService } from 'ng-zorro-antd/modal'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzFormModule } from 'ng-zorro-antd/form'
import moment from 'moment'
import { ExcelFieldListForm } from './interface'
import { NzMessageService } from 'ng-zorro-antd/message'
import { NzPaginationModule } from 'ng-zorro-antd/pagination'
import { UserStoreService } from '../../../../state/user.service'
import { findMenuItem } from '../../tool-function'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'

@Component({
    // selector: 'app-footer',
    standalone: true,
    imports: [
        CommonModule, 
        NzFormModule, 
        NzButtonModule, 
        FormsModule, 
        NzModalModule, 
        NzTableModule, 
        NzInputModule, 
        NzPaginationModule,
        NzSelectModule,
        NzInputNumberModule,
    ],
    templateUrl: './excel-field-match.component.html',
    styleUrl: './excel-field-match.component.css',
})
export class ExcelFieldMatchComponent {
    constructor(
        private message: NzMessageService,
        private userStoreService: UserStoreService
    ) {

        this.userStoreService.menuRole$.subscribe((data: any) => {
        const answer = findMenuItem(data, 'Excel Field Match', 'excel-field-matchs')
                            
            this.userRightInside = {
                read: answer.read,
                write: answer.write,
                update: answer.update,
                delete: answer.delete
            }
            console.log(this.userRightInside, 'answer')
        })
    }

    fieldForm: any = {
        dbFieldName: '',
        excelFieldName: '',
        sort: 0
    }

    userRightInside: any = {
        read: false,
        write: false,
        update: false,
        delete: false
    }

    searchForm: any = {
        page: 1,
        limit: 10
    }

    editForm: ExcelFieldListForm = {
        _id: '',
        functionCode: '',
        functionName: '',
        functionType: '',
        fieldLists: []
    }

    okText: string = 'Create'

    dataLists: any[] = []
    totals: number = 0
    editFormDialog: boolean = false
    removeDialog: boolean = false
    handleRemoveId: string = ''

    ngOnInit() {
        this.loadRecordLists()
        this.loadCodeType()
    }

    async submitForm() {
        const url = this.editForm._id === '' ? '/sys/excel-field-match/create' : `/sys/excel-field-match/update`

        const res = await postApiWithAuth(url, {
            
            ...this.editForm,
            ...this.editForm._id ? { _id: this.editForm._id} : {},
        })

        if (res.msg) {
            this.message.error(res.msg)
        } else if (res.matchedCount === 1 || !res.msg) {
            this.message.success('Save successful!')
            this.closeDialog()
            this.loadRecordLists()

            this.editForm = {
                _id: '',
                functionCode: '',
                functionName: '',
                functionType: '',
                fieldLists: []
            }
        }
    }

    async loadRecordLists() {
        const res = await postApiWithAuth('/sys/excel-field-match/list', this.searchForm)
        this.dataLists = res.lists
        this.totals = res.total
    }

    showDialog() {
        this.editFormDialog = true
    }

    closeDialog() {
        this.editFormDialog = false
    }

    handleRomeve(id: string) {
        this.handleRemoveId = id
        this.removeDialog = true
    }

    closeRemoveDialog() {
        this.removeDialog = false
    }

    async handleRemove() {
        const url = `/sys/excel-field-match/remove/${this.handleRemoveId}`

        const res: any = await getApiWithAuth(url)

        this.message.info(res.msg)

        this.closeRemoveDialog()
        this.loadRecordLists()
    }


    dateFormat(data: string) {
        return data ? moment(new Date(data)).format('DD-MM-YYYY HH:MM') : null
    }

    async getOneData(id:string) {
        const res = await getApiWithAuth(`/sys/excel-field-match/one/${id}`)
        this.editForm = res
        this.okText = 'Update'
        this.showDialog()
    }

    codeTypeLists: any[] = []
    async loadCodeType() {
        this.codeTypeLists = await getApiWithAuth('/base/code-type/get-type/ExcelFieldMatch')
    }

    fieldAddRow(): void {
        this.editForm.fieldLists = [...this.editForm.fieldLists, { ...this.fieldForm }]
        this.fieldForm = {
            dbFieldName: '',
            excelFieldName: '',
            sort: 0
        }
    }

    deleteRow(index: number): void {
        this.editForm.fieldLists.splice(index, 1)
        this.editForm.fieldLists = [...this.editForm.fieldLists]
    }

}