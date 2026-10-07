import { Component, OnInit } from '@angular/core'
import { CommonModule } from '@angular/common'
import { RouterLink, RouterLinkActive, RouterOutlet, Router } from '@angular/router'
import { FormsModule } from '@angular/forms'
import { deleteApiWithAuth, getApiWithAuth, postApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzModalModule, NzModalService } from 'ng-zorro-antd/modal'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzFormModule } from 'ng-zorro-antd/form'
import moment from 'moment'
// import { RoleForm } from './interface'
import { NzMessageService } from 'ng-zorro-antd/message'
import { NzPaginationModule } from 'ng-zorro-antd/pagination'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { StockTakeForm } from './interface'
import { findMenuItem } from '../../../tool-function'
import { UserStoreService } from '../../../../../state/user.service'
import { Subscription } from 'rxjs'

@Component({
    standalone: true,
    imports: [
        NzSelectModule, 
        NzCheckboxModule, 
        CommonModule, 
        NzFormModule, 
        NzButtonModule, 
        FormsModule, 
        NzModalModule, 
        NzTableModule, 
        NzInputModule, 
        NzPaginationModule, 
        FormsModule
    ],
    templateUrl: './stock-take-list.component.html',
    styleUrl: './stock-take-list.component.css',
})
export class StockTakeListComponent implements OnInit {
    private rightSubscription: Subscription
    constructor(
        private message: NzMessageService,
        private routeTo: Router,
        private userStoreService: UserStoreService
    ) {

        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Stock Take', 'stock-take-list')
            this.userRightInside = {
                read: answer?.read ?? false,
                write: answer.write ?? false,
                update: answer.update ?? false,
                delete: answer.delete ?? false,
                upload: answer.upload ?? false
                 // keep default value
            }
        })
    }

    ngOnDestroy() {
        if (this.userStoreService.menuRole$) {
            this.rightSubscription.unsubscribe()
        }
    }

    ngOnInit(): void {
        this.loadLocationList()
        this.loadStockTakeLists()
    }

    userRightInside: any = {
        read: false,
        write: false,
        update: false,
        delete: false
    }

    dataLists: any[] = []
    totals: number = 0
    handleId: string = ''
    finishDialogVisible: boolean = false
    cancelDialogVisible: boolean = false

    editForm: StockTakeForm = {
        actionName: '',
        actionPlaceId: '',
        remark: '',
    }

    searchForm: any = {
        page: 1,
        limit: 10
    }
    formDialog: boolean = false

    showDialog() {
        this.formDialog = true
    }

    closeDialog() {
        this.formDialog = false
    }

    async loadStockTakeLists() {
        const res = await postApiWithAuth('/product/stock-take/list', this.searchForm)
        this.dataLists = res.lists
        this.totals = res.total
    }

    placeLists: any[] = []
    async loadLocationList() {
        this.placeLists = await getApiWithAuth('/base/location/getAll')
    }

    async submitForm() {
        const res = await postApiWithAuth('/product/stock-take/create-form', this.editForm)
        if (res.msg) {
            this.message.error(res.msg)
        } else {
            this.message.success('Save successful!')
            this.closeDialog()
            this.loadStockTakeLists()

            this.editForm = {
                actionName: '',
                actionPlaceId: '',
                remark: ''
            }
        }
    }

    dateFormat(data: string) {
        return data ? moment(new Date(data)).format('DD-MM-YYYY HH:MM') : null
    }

    getToDetail(id: string) {
        this.routeTo.navigate([`/stock-take-form`], { queryParams: { id } })
    }

    openFinishDialog(id: string) {
        this.handleId = id
        this.finishDialogVisible = true
    }

    closeFinishDialog() {
        this.handleId = ''
        this.finishDialogVisible = false
    }

    openCancelDialog(id: string) {
        this.handleId = id
        this.cancelDialogVisible = true
    }

    closeCancelDialog() {
        this.handleId = ''
        this.cancelDialogVisible = false
    }

    async finshForm() {
        const res = await getApiWithAuth(`/product/stock-take/finish/${this.handleId}`)
        if (res.finished) {
            this.message.success(res.msg)
            this.loadStockTakeLists()
            this.closeFinishDialog()
        } else {
            this.message.error(res.msg)
            this.closeFinishDialog()
        }
    }

    async cancelForm() {
        const res = await getApiWithAuth(`/product/stock-take/void/${this.handleId}`)
        if (res.finished) {
            this.message.success(res.msg)
            this.loadStockTakeLists()
            this.closeFinishDialog()
        } else {
            this.message.error(res.msg)
            this.closeFinishDialog()
        }
    }
}