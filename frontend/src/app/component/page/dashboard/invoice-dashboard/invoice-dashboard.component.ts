import { Component, OnInit } from '@angular/core'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'
import { getApiWithAuth, postApiWithAuth } from '../../../../../tool/httpRequest-auth'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzModalModule, NzModalService } from 'ng-zorro-antd/modal'
import { NzInputModule } from 'ng-zorro-antd/input'
import { NzFormModule } from 'ng-zorro-antd/form'
import moment from 'moment'
import { NzMessageService } from 'ng-zorro-antd/message'
import { NzPaginationModule } from 'ng-zorro-antd/pagination'
import { DashboardReqDto, DashboardReqFilterDto } from './interface'
import { CanvasChartComponent } from '../../../components/chart/chart.component'
import { transformData, transformDataNoDate, transformDataPointsOnly, transformDate } from './function' 
import { NzSelectModule } from 'ng-zorro-antd/select'
import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import { UserStoreService } from '../../../../../state/user.service'
import { findMenuItem } from '../../../tool-function'
import { Subscription } from 'rxjs'
import { ListOption } from '../../../components/filter-list/filter-list.component'
import { FilterListApiSmall } from '../../../components/filter-list-small/filter-list.component'

@Component({
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
        CanvasChartComponent,
        NzSelectModule,
        NzDatePickerModule,
        FilterListApiSmall
    ],
    templateUrl: './invoice-dashboard.component.html',
    // styleUrl: './invoice-dashboard.component.css',
})
export class InvoiceDashboardComponent implements OnInit {
    private rightSubscription: Subscription
    private userRightSubscription: Subscription
    constructor(
        private message: NzMessageService,
        private userStoreService: UserStoreService
    ) {
        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'Dashboard', 'dashboard')
            this.userRightInside = {
                read: answer?.read ?? false
                 // keep default value
            }
        })

        this.userRightSubscription = this.userStoreService.user$.subscribe((data: any) => {
            if (!(data.superAdmin || data.manager)) {
                this.globalFilter.deptIds = data.deptIds
                this.globalFilter.placeIds = data.placeIds
            }
        })              
    }

    ngOnDestroy() {
        if (this.userStoreService.menuRole$) {
            this.rightSubscription.unsubscribe()
        }
        if (this.userStoreService.user$) {
            this.userRightSubscription.unsubscribe()
        }
    }


    userRightInside = {
        read: false
    }

    ngOnInit(): void {
        this.runSearch()
        
    }

    async runSearch() {
        await this.getByDeptAndQtys()
        await this.getByDeptAndPrice()
        await this.getByTypeAndPrice()
        await this.getByTypeAndQtys()

        await this.getByLocationAndQtys() 
        await this.getByLocationAndPrice()
        await this.getbyYearMonthAndQtys()
        await this.getbyYearMonthAndPrice() 
    }

    globalFilter: DashboardReqFilterDto = {
        placeIds: [],
        deptIds: [],
        typeIds: []
    }

    byTypeAndQtysLoading: boolean = false
    byTypeAndQtys: any = {}
    async getByTypeAndQtys() {
        this.byTypeAndQtysLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'type',
            valueField: 'qtys'
        }

        const res = await this.runQueryData(dataQuery)
        const qtys = transformDataNoDate(res, 'column', false, 'typeName', 'qtys')
        const qtysDataSet = qtys[0].dataPoints
        
        this.byTypeAndQtys = {
                    data: [
                        {
                            name: 'Total Sales Qtys',
                            dataPoints: qtysDataSet
                        }
                    ],
                    animationEnabled: true,
                    axisY: {
                        title: 'Qtys'
                    },
                    axisX: {
                        title: 'Product Types'
                    },
                    toolTip: {
                        shared: true
            }
        }
        this.byTypeAndQtysLoading = true
    }

    byTypeAndPriceLoading: boolean = false
    byTypeAndPrice: any = {}
    async getByTypeAndPrice() {
        this.byTypeAndPriceLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'type',
            valueField: 'price'
        }

        const res = await this.runQueryData(dataQuery)
        const price = transformDataNoDate(res, 'column', false, 'typeName', 'price')
        const priceDataSet = price[0].dataPoints
        this.byTypeAndPrice = {
            data: [
                {
                    name: 'Total Sales Price',
                    dataPoints: priceDataSet
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Price'
            },
            axisX: {
                title: 'Product Types',
             /*   valueFormatString: "YYYY - MMM",
                xValueType: "dateTime" */
            },
            toolTip: {
                shared: true
            }
        }
        this.byTypeAndPriceLoading = true
    }

    byDeptAndQtysLoading: boolean = false
    byDeptAndQtys: any = {}
    async getByDeptAndQtys() {
        this.byDeptAndQtysLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'dept',
            valueField: 'qtys'
        }

        const res = await this.runQueryData(dataQuery)
        const price = transformDataNoDate(res, 'column', false, 'deptName', 'qtys')
        const priceDataSet = price[0].dataPoints

        this.byDeptAndQtys = {
            data: [
                {
                    name: 'Total Sales Qtys',
                    dataPoints: priceDataSet
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Qtys'
            },
            axisX: {
                title: 'Departments',
             /*   valueFormatString: "YYYY - MMM",
                xValueType: "dateTime" */
            },
            toolTip: {
                shared: true
            }
        }
        this.byDeptAndQtysLoading = true
    }

    byDeptAndPriceLoading: boolean = false
    byDeptAndPrice: any = {}
    async getByDeptAndPrice() {
        this.byDeptAndPriceLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'dept',
            valueField: 'price'
        }

        const res = await this.runQueryData(dataQuery)
        const price = transformDataNoDate(res, 'column', false, 'deptName', 'price')
        const priceDataSet = price[0].dataPoints

        this.byDeptAndPrice = {
            data: [
                {
                    name: 'Total Sales Price',
                    dataPoints: priceDataSet
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Price'
            },
            axisX: {
                title: 'Departments',
             /*   valueFormatString: "YYYY - MMM",
                xValueType: "dateTime" */
            },
            toolTip: {
                shared: true
            }
        }
        this.byDeptAndPriceLoading = true
    }

    byLocationAndQtysLoading: boolean = false
    byLocationAndQtys: any = {}
    async getByLocationAndQtys() {
        this.byLocationAndQtysLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'location',
            valueField: 'qtys'
        }

        const res = await this.runQueryData(dataQuery)
        const qtys = transformDataNoDate(res, 'column', false, 'placeName', 'qtys')
        const qtysDataSet = qtys[0].dataPoints

        this.byLocationAndQtys = {
            data: [
                {
                    name: 'Total Sales Qtys',
                    dataPoints: qtysDataSet
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Qtys'
            },
            axisX: {
                title: 'Locations',
             /*   valueFormatString: "YYYY - MMM",
                xValueType: "dateTime" */
            },
            toolTip: {
                shared: true
            }
        }
        this.byLocationAndQtysLoading = true
    }

    byLocationAndPriceLoading: boolean = false
    byLocationAndPrice: any = {}
    async getByLocationAndPrice() {
        this.byLocationAndPriceLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'location',
            valueField: 'price'
        }

        const res = await this.runQueryData(dataQuery)
        const price = transformDataNoDate(res, 'column', false, 'placeName', 'price')
        const priceDataSet = price[0].dataPoints

        this.byLocationAndPrice = {
            data: [
                {
                    name: 'Total Sales Price',
                    dataPoints: priceDataSet
                }
            ],
            animationEnabled: true,
            axisY: {
                title: 'Price'
            },
            axisX: {
                title: 'Locations'
            },
            toolTip: {
                shared: true
            }
        }
        this.byLocationAndPriceLoading = true
    }

    byYearMonthAndQtysLoading: boolean = false
    byYearMonthAndQtys: any = {}
    async getbyYearMonthAndQtys() {
        this.byYearMonthAndQtysLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'year-month',
            valueField: 'qtys'
        }

        const res = await this.runQueryData(dataQuery)
        const dataSet = transformDate(res, 'qtys', ['year', 'month'])

        this.byYearMonthAndQtys = {
            data: [{
                type: 'spline',
                name: 'Sales',
                xValueType: 'dateTime',
                dataPoints: dataSet
            }],
            animationEnabled: true,
            xValueType: "dateTime",
            axisY: {
                title: 'Qtys'
            },
            axisX: {
                title: 'Year - Month',
                valueFormatString: "YYYY - MMM",
                xValueType: "dateTime"
            },
            toolTip: {
                shared: true
            }
        }
        this.byYearMonthAndQtysLoading = true
    }

    byYearMonthAndPriceLoading: boolean = false
    byYearMonthAndPrice: any = {}
    async getbyYearMonthAndPrice() {
        this.byYearMonthAndPriceLoading = false
        const dataQuery: DashboardReqDto = {
            dataTypeValue: 'year-month',
            valueField: 'price'
        }

        const res = await this.runQueryData(dataQuery)
        const dataSet = transformDate(res, 'price', ['year', 'month'])

        this.byYearMonthAndPrice = {
            data: [{
                type: 'spline',
                name: 'Sales',
                xValueType: 'dateTime',
                dataPoints: dataSet
            }],
            animationEnabled: true,
            xValueType: "dateTime",
            axisY: {
                title: 'Price'
            },
            axisX: {
                title: 'Year - Month',
                valueFormatString: "YYYY - MMM",
                xValueType: "dateTime"
            },
            toolTip: {
                shared: true
            }
        }
        this.byYearMonthAndPriceLoading = true
    }

    async runQueryData(dataQuery: DashboardReqDto) {
        const finalQuery = {
            ...dataQuery,
            filter: this.globalFilter
        }

        return await postApiWithAuth('/query/invoice', finalQuery)
    }
//  Department data
    departmentTotal: number = 0
    departments: any = []
    departmentOptions: ListOption[] = [];
    departmentSearchForm: any = {
        name: '',
        page: 1,
        limit: 10
    }
    loadDeptLists = async (
        url: string,
        body: Record<string, unknown>,
    ): Promise<ListOption[]> => {
        const payload = {
          ...this.departmentSearchForm,
          ...body,
        }
    
        try {
            const res = await postApiWithAuth(url, payload);
        
            this.departments = res ?? [];
            this.departmentTotal = res.length ?? 0;
        
            this.departmentOptions = this.departments.map((a: any) => ({
                value: a.id,
                label: `${a.deptCode} - ${a.deptName}`,
            }))
    
            return this.departmentOptions
        } catch (error) {
            this.message.error('API request failed')
            throw error
        }
    }
//  Product Type data
    productTypeTotal: number = 0
    productTypes: any = []
    productTypeOptions: ListOption[] = [];
    productTypeSearchForm: any = {
        name: '',
        page: 1,
        limit: 10
    }
    loadTypeLists = async (
        url: string,
        body: Record<string, unknown>,
    ): Promise<ListOption[]> => {
        const payload = {
          ...this.productTypeSearchForm,
          ...body,
        }
    
        try {
            const res = await postApiWithAuth(url, payload);
        
            this.productTypes = res ?? [];
            this.productTypeTotal = res.length ?? 0;
        
            this.productTypeOptions = this.productTypes.map((a: any) => ({
                value: a.id,
                label: `${a.typeCode} - ${a.typeName}`,
            }))
    
            return this.productTypeOptions
        } catch (error) {
            this.message.error('API request failed')
            throw error
        }
    }
    // Location data
    placeTotal: number = 0
    places: any = []
    placeOptions: ListOption[] = [];
    placeSearchForm: any = {
        placeCode: '',
        placeName: '',
        page: 1,
        limit: 10
    }
    loadPlaceLists = async (
        url: string,
        body: Record<string, unknown>,
    ): Promise<ListOption[]> => {
        const payload = {
          ...this.placeSearchForm,
          ...body,
        }
    
        try {
            const res = await postApiWithAuth(url, payload);
        
            this.places = res ?? [];
            this.placeTotal = res.length ?? 0;
        
            this.placeOptions = this.places.map((a: any) => ({
                value: a.id,
                label: `${a.placeCode}`,
            }))
    
            return this.placeOptions
        } catch (error) {
            this.message.error('API request failed')
            throw error
        }
    }
}