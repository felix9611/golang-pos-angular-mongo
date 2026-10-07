import { Routes } from '@angular/router'
import { AuthGuard } from '../../state/AuthGuard'
import { UserInfoComponent } from './page/userInfo/user-info.component'
import { DepartmentComponent } from './page/department/department.component'
import { RoleComponent } from './page/role/role.component'
import { UsersComponent } from './page/users/users.component'
import { CodeTypeComponent } from './page/code-type/code-type.component'
import { VendorComponent } from './page/vendor/vendor.component'
import { LocationComponent } from './page/location/location.component'
import { ActionRecordComponent } from './page/action-record/action-record.component'
import { TaxInformationComponent } from './page/tax-information/tax-information.component'
import { InvoiceDashboardComponent } from './page/dashboard/invoice-dashboard/invoice-dashboard.component'
import { MenuListComponent } from './page/menu/menu.component'
import { ExcelFieldMatchComponent } from './page/excel-field-match/excel-field-match.component'
import { ProductTypeComponent } from './page/product/product-type/product-type.component'
import { ProductListComponent } from './page/product/product-list/product-list.component'
import { ProductFormComponent } from './page/product/product-form/product-form.component'
import { StockInComponent } from './page/stock/stock-in/stock-in.component'
import { ProductLocationListComponent } from './page/stock/product-location/product-location.component'
import { StockOutComponent } from './page/stock/stock-out/stock-out.component'
import { StockMoveComponent } from './page/stock/stock-move/stock-move.component'
import { InventoryRecordComponent } from './page/stock/inventory-record/inventory-record.component'
import { MemberClassesComponent } from './page/member/member-class/member-class.component'
import { MemberFormComponent } from './page/member/member-form/member-form.component'
import { MemberListComponent } from './page/member/member-list/member-list.component'
import { SalePointComponent } from './page/sale-point/sale-point.component'
import { InvoiceDetailComponent } from './page/invoice/invoice-detail/invoice-detail.component'
import { InvoiceListComponent } from './page/invoice/invoice-list/invoice-list.component'
import { StockDashboardComponent } from './page/dashboard/stock-dashboard/stock-dashboard.component'
import { StockTakeListComponent } from './page/stock-take/stock-take-list/stock-take-list.component'
import { StockTakeFormComponent } from './page/stock-take/stock-take-form/stock-take-form.component'
import { ReturnInvoiceFormComponent } from './page/invoice/return-invoice-form/return-invoice-form.component'
import { ReturnInvoiceListComponent } from './page/invoice/return-invoice-list/return-invoice-list.component'
import { ReturnInvoiceDetailComponent } from './page/invoice/return-invoice-detail/return-invoice-detail.component'

export const pagesRoutes: Routes = [
    {
        path: 'action-record',
        component: ActionRecordComponent,
        canActivate: [AuthGuard]
    },
    
    {
        path: 'stock-dashboard',
        component: StockDashboardComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'invoice-dashboard',
        component: InvoiceDashboardComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'excel-field-matchs',
        component: ExcelFieldMatchComponent,
        canActivate: [AuthGuard]
    },
   
    {
        path: 'code-type',
        component: CodeTypeComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'departments',
        component: DepartmentComponent,
        canActivate: [AuthGuard]
    },
   
    {
        path: 'location',
        component: LocationComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'tax-information',
        component: TaxInformationComponent,
        canActivate: [AuthGuard]
    },
   
    {
        path: 'role',
        component: RoleComponent,
        canActivate: [AuthGuard]
    },
   
    {
        path: 'menu',
        component: MenuListComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'users',
        component: UsersComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'user-info',
        component: UserInfoComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'vendor',
        component: VendorComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'product-type',
        component: ProductTypeComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'product-list',
        component: ProductListComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'product-create',
        component: ProductFormComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'product-detail',
        component: ProductFormComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'product-location-lists',
        component: ProductLocationListComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'stock-in',
        component: StockInComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'stock-move',
        component: StockMoveComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'stock-out',
        component: StockOutComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'inventory-record',
        component: InventoryRecordComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'member-class',
        component: MemberClassesComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'member-list',
        component: MemberListComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'create-member',
        component: MemberFormComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'member-detail',
        component: MemberFormComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'sale-point',
        component: SalePointComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'invoice-list',
        component: InvoiceListComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'invoice-detail',
        component: InvoiceDetailComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'create-return-invoice',
        component: ReturnInvoiceFormComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'return-invoice-list',
        component: ReturnInvoiceListComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'return-invoice-detail',
        component: ReturnInvoiceDetailComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'stock-take-list',
        component: StockTakeListComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'stock-take-form',
        component: StockTakeFormComponent,
        canActivate: [AuthGuard]
    }
]