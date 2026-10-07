import { Component, OnChanges } from '@angular/core'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { getApiWithAuth, postApiWithAuth } from '../../../../tool/httpRequest-auth'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzModalModule, NzModalService } from 'ng-zorro-antd/modal'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzFormModule } from 'ng-zorro-antd/form'
import moment from 'moment'
import { TaxInfomationForm } from './interface'
import { NzMessageService } from 'ng-zorro-antd/message'
import { NzPaginationModule } from 'ng-zorro-antd/pagination'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { findMenuItem } from '../../tool-function'
import { UserStoreService } from '../../../../state/user.service'
import { Subscription } from 'rxjs'
import { downloadTempExcelFile, formatJson, readExcelFile } from '../../../../tool/excel-helper'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { UploadDialogComponent } from '../../components/upload-dialog-component/upload-dialog-component.component'
import { DownloadExcelTemplateComponent } from '../../components/download-template-component/download-template-component.component'
import { DownloadExcelDataComponent } from '../../components/download-excel-component/download-excel-data-component.component'

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
        NzDatePickerModule,
        NzInputNumberModule, 
        NzUploadModule,
        DownloadExcelTemplateComponent,
        UploadDialogComponent,
        DownloadExcelDataComponent
    ],
    templateUrl: './tax-information.component.html',
    styleUrl: './tax-information.component.css',
})
export class TaxInformationComponent {
    private rightSubscription: Subscription
    constructor(
        private message: NzMessageService,
        private userStoreService: UserStoreService
    ) {

        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Tax Information', 'tax-information')
            this.userRightInside = {
                read: answer?.read ?? false,
                write: answer.write ?? false,
                update: answer.update ?? false,
                delete: answer.delete ?? false,
                upload: answer.upload ?? false
                 // keep default value
            }
            this.excelFileSetting.code = answer?.excelFunctionCode ?? ''
            this.preLoadExcelSetting()
        })
        
    }

    ngOnDestroy() {
        if (this.userStoreService.menuRole$) {
            this.rightSubscription.unsubscribe()
        }
    }

    userRightInside: any = {
        read: false,
        write: false,
        update: false,
        delete: false
    }

    searchForm: any = {
        date: [],
        page: 1,
        limit: 10
    }

    editForm: TaxInfomationForm =  {
        _id: '',
        nationCode: '',
        nationName: '',
        countryCode: '',
        countryName: '',
        taxType: '',
        taxCode: '',
        taxName: '',
        taxRate: 0,
        importRate: 0,
        remark:  ''
    }

    okText: string = 'Create'

    dataLists: any[] = []
    totals: number = 0
    editFormDialog: boolean = false
    removeDialog: boolean = false
    handleRemoveId: string = ''

    fileList: any[] = []
    deptLists: any[] = []
    placeLists: any[] = []

    department?: any = {}


    ngOnInit() {
        this.loadTaxInfoLists()
        this.loadDeptLists()
        this.loadPlaceLists()
        
    }

    excelFileSetting: any = {
        code: ''
    }
    dbFieldList: string[] = []
    excelFieldList: string[] = []
    async preLoadExcelSetting() {
        const res = await getApiWithAuth(`/sys/excel-field-match/code/${this.excelFileSetting.code}`)
        this.dbFieldList = res.fieldLists.map((item: any) => item.dbFieldName)
        this.excelFieldList = res.fieldLists.map((item: any) => item.excelFieldName)
    }

    async loadPlaceLists() {
        this.placeLists = await getApiWithAuth('/base/location/getAll')
    }

    async loadDeptLists() {
        this.deptLists = await getApiWithAuth('/base/budget/getBudgetSummary')
    }

    async submitForm() {
        const url = this.editForm._id === '' ? '/base/tax-information/create' : `/base/tax-information/update`

        console.log(this.editForm)
        const res = await postApiWithAuth(url, {
            ...this.editForm,
            taxRate: this.editForm.taxRate / 100,
            importRate: this.editForm.importRate ? this.editForm.importRate / 100 : 0
        })

        if (res.msg) {
            this.message.error(res.msg)
        } else if (res.matchedCount === 1 || !res.msg) {
            this.message.success('Save successful!')
            this.closeDialog()
            this.loadTaxInfoLists()

            this.editForm = {
                _id: '',
                nationCode: '',
                nationName: '',
                countryCode: '',
                countryName: '',
                taxType: '',
                taxCode: '',
                taxName: '',
                taxRate: 0,
                importRate: 0,
                remark:  ''
            }
        }
    }

    async loadTaxInfoLists() {
        const res = await postApiWithAuth('/base/tax-information/list', this.searchForm)
        this.dataLists = res.lists
        this.totals = res.total
    }

    showDialog() {
        this.editFormDialog = true
    }

    closeDialog() {
        this.editForm = {
            _id: '',
            nationCode: '',
            nationName: '',
            countryCode: '',
            countryName: '',
            taxType: '',
            taxCode: '',
            taxName: '',
            taxRate: 0,
            importRate: 0,
            remark:  ''
        }
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
        const url = `/base/tax-information/remove/${this.handleRemoveId}`

        const res: any = await getApiWithAuth(url)

        this.message.info(res.msg)
        this.loadTaxInfoLists()
        this.closeRemoveDialog()
        
    }


    dateFormat(data: string) {
        return data ? moment(new Date(data)).format('DD-MM-YYYY HH:MM') : null
    }

    async getOneData(id:string) {
        const res = await getApiWithAuth(`/base/tax-information/one/${id}`)
        this.editForm = {
            ...res,
            taxRate: res.taxRate * 100,
            importRate: res.importRate * 100
        }
        this.department = res.department
        this.okText = 'Update'
        this.showDialog()
    }
}