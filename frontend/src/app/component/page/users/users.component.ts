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
import { UserForm } from './interface'
import { NzMessageService } from 'ng-zorro-antd/message'
import { NzPaginationModule } from 'ng-zorro-antd/pagination'
import { NzUploadFile, NzUploadModule } from 'ng-zorro-antd/upload'
import { imgToBase64, uploadImgToBase64 } from '../../../../tool/imageUpload'
import { Observable, Observer, Subscription } from 'rxjs'
import { NzSelectModule } from 'ng-zorro-antd/select'
import { findMenuItem } from '../../tool-function'
import { UserStoreService } from '../../../../state/user.service'

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
        NzUploadModule,
        NzSelectModule
    ],
    templateUrl: './users.component.html',
    styleUrl: './users.component.css',
})
export class UsersComponent {
    private rightSubscription: Subscription
    constructor(
        private message: NzMessageService,
        private userStoreService: UserStoreService
    ) {

        this.rightSubscription = this.userStoreService.menuRole$.subscribe((data: any) => {
            const answer = findMenuItem(data, 'User', 'users')
            this.userRightInside = {
                read: answer?.read ?? false,
                write: answer.write ?? false,
                update: answer.update ?? false,
                delete: answer.delete ?? false
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
        delete: false
    }


    searchForm: any = {
        page: 1,
        limit: 10
    }

    editForm: UserForm=  {
        _id: '',
        username: '',
        avatarBase64: '',
        deptId: 0,
        email: '',
        roles: [],
    }

    okText: string = 'Create'

    dataLists: any[] = []
    totals: number = 0
    editFormDialog: boolean = false
    removeDialog: boolean = false
    handleRemoveId: string = ''

    fileList: any[] = []
    deptLists: any[] = []
    roleLists: any[] = []

    department?: any = {}

    updateFile(file: any) {
        let testmsg = file.name.substring(file.name.lastIndexOf('.')+1)
        console.log(testmsg)
        const isJpg = testmsg === 'jpg' || testmsg === 'png' || testmsg === 'JPG' || testmsg === 'PNG'
        const isLt2M = file.size / 1024 / 1024 < 3
        if (!isJpg) {
            this.message.error('Only upload jpg or png file!')
        }

        if (!isLt2M) {
            this.message.error('File size cannot over 3MB!')
        }
        if (isJpg && isLt2M){
            this.fileList.push(file)
        }
        this.imgToBase64()
    }

    imgToBase64() {
        this.fileList.map(async (file: any) => {
            const response: any = await uploadImgToBase64(file.raw)
            console.log(response)
        })
    }


    ngOnInit() {
        this.loadUserLists()
        this.loadDeptLists()
        this.loadRoleLists()
    }

    async loadRoleLists() {
        this.roleLists = await getApiWithAuth('/sys/role/getAll')
    }

    async loadDeptLists() {
        this.deptLists = await getApiWithAuth('/sys/department/getAll')
    }

    async submitForm() {
        const url = this.editForm._id === '' ? '/sys/user/create-user' : `/sys/user/update-uesr`

        const res = await postApiWithAuth(url, this.editForm._id? {
            avatarBase64: this.avatarUrl ? this.avatarUrl : this.editForm.avatarBase64,
            ...this.editForm
        } : {
            key: '9@0UtWV:;}m@HkjHyVU=',
            userData: {
                avatarBase64: this.avatarUrl,
                ...this.editForm
            }
        })

        if (res.msg) {
            this.message.error(res.msg)
        } else if (res.matchedCount === 1 || !res.msg) {
            this.message.success('Save successful!')
            this.closeDialog()
            this.loadUserLists()

            this.editForm = {
                _id: '',
                username: '',
                avatarBase64: '',
                deptId: 0,
                email: '',
                roles: []
            }
        }
    }

    async loadUserLists() {
        const res = await postApiWithAuth('/sys/user/list', this.searchForm)
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
        const url = `/sys/user/invalidate-user/${this.handleRemoveId}`

        const res: any = await getApiWithAuth(url)

        this.message.info(res.msg)
        this.loadUserLists()
        this.closeRemoveDialog()
        
    }


    dateFormat(data: string) {
        return data ? moment(new Date(data)).format('DD-MM-YYYY HH:MM') : null
    }

    async getOneData(id:string) {
        const res = await getApiWithAuth(`/sys/user/one/${id}`)
        this.editForm = res
        this.department = res.department
        this.okText = 'Update'
        this.showDialog()
    }

    // upload

    loading = false
    avatarUrl?: string = ''
    previewImage: string = ''
    
    beforeUpload = (file: NzUploadFile, _fileList: NzUploadFile[]): Observable<boolean> =>
        new Observable((observer: Observer<boolean>) => {
        const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
        if (!isJpgOrPng) {
            this.message.error('You can only upload JPG or PNGfile!')
            observer.complete()
            return
        }
        const isLt2M = file.size! / 1024 / 1024 < 3
        if (!isLt2M) {
            this.message.error('Image must smaller than 3MB!')
            observer.complete()
            return
        }
    
        observer.next(isJpgOrPng && isLt2M)
        observer.complete()
    })
        
    
    async handleChange(info: { file: any }) {
        const response: any = await uploadImgToBase64(info.file.originFileObj)
        this.avatarUrl = response.data
    }

}