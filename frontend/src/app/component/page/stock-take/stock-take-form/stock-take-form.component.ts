import { Component, OnInit } from '@angular/core'
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
import { FormsModule } from '@angular/forms'
import { CommonModule } from '@angular/common'
import { ActivatedRoute, Router } from '@angular/router'
import { StockTakeFormEdit, StockTakeItemDto, StockTakeItemFromDto } from './interface'
import { getApiWithAuth, postApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import { MatButtonModule } from '@angular/material/button'
import { MatIconModule } from '@angular/material/icon'
import { UserStoreService } from '../../../../../state/user.service'
import { findMenuItem } from '../../../tool-function'
import { Subscription } from 'rxjs'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'

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
        FormsModule,
        NzDatePickerModule,
        MatIconModule, 
        MatButtonModule,
        NzInputNumberModule
    ],
    templateUrl: './stock-take-form.component.html',
    styleUrl: './stock-take-form.component.css',
})
export class StockTakeFormComponent implements OnInit {
    private rightSubscription: Subscription
await: any
    constructor(
        private route: ActivatedRoute, 
        private routeTo: Router,
        private message: NzMessageService,
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

    userRightInside: any = {
        read: false,
        write: false,
        update: false,
        delete: false,
        upload: false
    }

    theId: any = ''

    editForm: StockTakeFormEdit = {
        actionName: '',
        actionPlaceId: '',
        remark: '',
        stockTakeItems: [],
        createdTime: '',
        finishTime: '',
        createBy: '',
        _id: '',
        status: 0
      }

    itemForm: StockTakeItemFromDto = {
        stockTakeId: '',
        productId: '',
        productCode: '',
        productName: '',
        placeId: '',
        status: '',
        remark: '',
        orgQty: 0,
        checkQty: 0
    }

    ngOnInit() {
        this.loadLocationList()
        this.route.queryParams.subscribe((x: any) => {
            if (x.id) {
                this.theId = x.id
                this.getOne()
            }
        })
    }

    async getOne() {
        this.editForm = await getApiWithAuth(`/product/stock-take/one/${ this.theId}`)    
    }

    placeLists: any[] = []
    async loadLocationList() {
        this.placeLists = await getApiWithAuth('/base/location/getAll')
    }

    async updateForm() {
        const res = await postApiWithAuth('/product/stock-take/update-form', this.editForm)
        if (res.finished) {
            this.message.success(res.msg)
        } else {
            this.message.error(res.msg)
        }
    }

    async assetCodeChanged(event: any) {
        if (event) {
            const data = await getApiWithAuth(`/product/product-list/code/${event}`)
            this.itemForm.productId = data._id
          //  this.itemForm.productCode = data.productCode
            this.itemForm.productName = data.productName
            this.placeCheckStatus(data._id)
        }
    }


    async placeCheckStatus(productId: string) {
        const res = await postApiWithAuth('/product/product-list/location/check', { productId, locationId: this.editForm.actionPlaceId })
        if (res.status === true) {
            this.itemForm.status = 'Exist'
            this.itemForm.orgQty = res.data.qty
        } else if (res.status === false) {
            this.itemForm.status = 'Not Exist OR Wrong Location'
        }
    }

    async submitItem() {
        const finalData = {
            stockTakeId: this.editForm._id,
            productId: this.itemForm.productId,
            productCode: this.itemForm.productCode,
            placeId: this.editForm.actionPlaceId,
            status: this.itemForm.status,
            remark: this.itemForm.remark,
            qty: this.itemForm.checkQty
        }

        const res = await postApiWithAuth('/product/stock-take/item-submit', finalData)

        if (res._id) {
            this.message.success('Added!')
            this.getOne()
        } else {
            this.message.error('Ooops! something wrong! Please try again!')
        }

    }

    dateFormat(data: string) {
        return data ? moment(new Date(data)).format('DD-MM-YYYY HH:MM') : null
    }

    backToList() {
        this.routeTo.navigate([`/stock-takes`])
    }
}