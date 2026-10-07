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
    ],
    templateUrl: './invoice-dashboard.component.html',
    // styleUrl: './invoice-dashboard.component.css',
})
export class InvoiceDashboardComponent implements OnInit {
    private rightSubscription: Subscription
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
                    
    }

    ngOnDestroy() {
        if (this.userStoreService.menuRole$) {
            this.rightSubscription.unsubscribe()
        }
    }


    userRightInside = {
        read: false
    }

    ngOnInit(): void {
        this.runSearch()
   //     this.getByTotalCost()
  //      this.getByTotalCount()
        this.loadTypeList()
        this.loadDeptList()
        this.loadLocationList()
        
    }

    typeLists: any[] = []
    async loadTypeList() {
        this.typeLists = await getApiWithAuth('/product/product-type/getAll')
    }

    deptLists: any[] = []
    async loadDeptList() {
        this.deptLists = await getApiWithAuth('/sys/department/getAll')
    }

    placeLists: any[] = []
    async loadLocationList() {
        this.placeLists = await getApiWithAuth('/base/location/getAll')
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

    globalFilter: DashboardReqFilterDto = {}

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

        return await postApiWithAuth('/invoice/query/data-group-by', finalQuery)
    }

}