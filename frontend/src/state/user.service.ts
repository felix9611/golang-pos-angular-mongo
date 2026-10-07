import { inject, Injectable } from '@angular/core'
import { BehaviorSubject } from 'rxjs'
import { UserInfo } from './interface'
import { Router } from '@angular/router'
import { LocalStorageService } from './localStorage.service'
import { getApiWithAuth, postApiWithAuth } from '../tool/httpRequest-auth'

@Injectable({
  providedIn: 'root'
})
export class UserStoreService {
    constructor(
        private localStorageService: LocalStorageService
    ) {}

    router = inject(Router)
    private accessToken: string = ''
    private initialState: UserInfo = { 
        _id: '',
        username: '',
        accessToken: '',
        avatarBase64: '',
        deptId: 0,
        roles: [],
        roleLists: [],
        email: '', 
        loginRecords: []
    }
    private userMenu: any[] = []
    private userMenuRole: any[] = []
    private tokenSubject = new BehaviorSubject<string>(this.accessToken)
    private userSubject = new BehaviorSubject<UserInfo>(this.initialState)
    private menuSubject = new BehaviorSubject<any[]>(this.userMenu)
    private menuRoleSubject = new BehaviorSubject<any[]>(this.userMenu)
    user$ = this.userSubject.asObservable()
    token$ = this.tokenSubject.asObservable()
    menu$ = this.menuSubject.asObservable()
    menuRole$ = this.menuRoleSubject.asObservable()

    get user(): UserInfo {
        return this.userSubject.value
    }

    get token(): string {
        return this.tokenSubject.value
    }

    async isAuthenticated(): Promise<boolean> {
        const token = this.getAccessToken()
        const check = await this.verifyToken()
        if (check) {
            return !!token // Returns true if token exists
        } else {
            this.router.navigate(['/login'])
            return false
        }
    }

    setAccessToken(token: string): void {
        this.tokenSubject.next(`Bearer ${token}`)
        this.localStorageService.setItem('accessToken', `Bearer ${token}`)
    }

    getAccessToken() {
        return this.localStorageService.getItem('accessToken')
    }

    setUser(info: UserInfo): void {
        this.userSubject.next(info)
    }

    setMenu(menu: any[]): void {
        this.menuSubject.next(menu)
    }

    setMenuRole(menuRole: any[]): void {
        this.menuRoleSubject.next(menuRole)
    }
    

    logout(): void {
        this.userSubject.next({ 
            _id: '',
            username: '',
            accessToken: '',
            avatarBase64: '',
            deptId: 0,
            roles: [],
            roleLists: [],
            email: '',
            loginRecords: []

        })
        this.menuSubject.next([])
        this.menuRoleSubject.next([])
        localStorage.clear()
        this.tokenSubject.next('')
        this.toLoginPage()
    }

    toLoginPage() {
        this.router.navigate(['/login'])
    }

    async loadUserInfo() {
        const data = await getApiWithAuth('/auth/user-profile')
        this.setUser(data)
    }

    async verifyToken(): Promise<boolean> {
        const result: { status: boolean } = await getApiWithAuth('/auth/verify-token')
        if (result.status) {
            return true
        } else {
            return false
        }
    }

    async loadMenus() {
        let menuIds: any = []
        this.user$.subscribe(user => {
            user.roleLists?.forEach(async (role: any) => {
                
                menuIds.push(...role.menuIds)
                const data = await postApiWithAuth('/sys/menu/user/tree-menu', { ids: menuIds })
                this.setMenu(data)
            })
        })
        
    }

    async loadMenuRoles() {
        this.user$.subscribe(async user => {
            const data = await postApiWithAuth('/sys/role/list-permission', { roleIds: user.roles })
            this.setMenuRole(data)
        })
    }

}