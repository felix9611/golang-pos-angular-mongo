import { Component, OnInit } from '@angular/core'
import { ProductFormDto, StockInOutProductLocationDto, StockMoveProductLocationDto } from '../interface'
import { getApiWithAuth, postApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { CommonModule } from '@angular/common'
import { NzFormModule } from 'ng-zorro-antd/form'
import { MatButtonModule } from '@angular/material/button'
import { MatIconModule } from '@angular/material/icon'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { NzInputModule } from 'ng-zorro-antd/input'
import { FormsModule } from '@angular/forms'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import { NzPaginationModule } from 'ng-zorro-antd/pagination'
import { NzModalModule } from 'ng-zorro-antd/modal'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { NzMessageService } from 'ng-zorro-antd/message'
import { ProductLocationListComponent } from '../product-location/product-location.component'

@Component({
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
        NzInputNumberModule,
        ProductLocationListComponent
    ],
    templateUrl: './stock-out.component.html'
})
export class StockOutComponent implements OnInit {
    constructor(
        private message: NzMessageService,
    ) {}

    ngOnInit(): void {
        this.loadLocationList()
    }

    reloadProductLocation: boolean = true

    stockInForm: StockInOutProductLocationDto  = {
        productCode: '',
        productId: '',
        placeCode: '',
        locationId: '',
        qty: 0,
        totalPrice: 0,
        totalCost: 0
    }

    productForm: ProductFormDto = {
        _id: '',
        productCode: '',
        productName: '',
        itemCode: '',
        brandCode: '',
        brandName: '',
        typeId: '',
        deptId: '',
        vendorId: '',
        unit: '',
        costPrice: 0,
        retailPrice: 0,
        description: '',
        remark: ''
    }

    locationList: any[] = []
    async loadLocationList() {
        const data = await getApiWithAuth('/base/location/getAll')
        this.locationList = data
    }

    async productCodeChanged(event: any) {
        console.log(event)
        if (event) {
            
            const data = await getApiWithAuth(`/product/product-list/code/${event}`)
            this.stockInForm.productId = data._id
            this.productForm = data
        }
    }

    qtyChanged(event: any) {
        this.stockInForm.totalCost = event * this.productForm.costPrice
        this.stockInForm.totalPrice = event * this.productForm.retailPrice
    }

    resetForm() {
        this.stockInForm  = {
            productCode: '',
            productId: '',
            placeCode: '',
            locationId: '',
            qty: 0,
            totalPrice: 0,
            totalCost: 0
        }

        this.productForm = {
            _id: '',
            productCode: '',
            productName: '',
            itemCode: '',
            brandCode: '',
            brandName: '',
            typeId: '',
            deptId: '',
            vendorId: '',
            unit: '',
            costPrice: 0,
            retailPrice: 0,
            description: '',
            remark: ''
        }
    }

    async submitForm() {
        this.reloadProductLocation = false
        const res = await postApiWithAuth('/product/product-list/stock-out', this.stockInForm)
        if (res) {
            this.message.success('Stock out success!')
            this.reloadProductLocation = true
        } else {
            this.message.error('Something error! Please try again!')
            this.reloadProductLocation = true
        }
    }
}
