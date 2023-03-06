import router from './router'
import { getters } from './store/store'
import { actions } from './store/actions'
import { singletonPromise } from "./utils/singleton";

// 设置 autoLogin 中的 this 为 actions，因为 actions.silentLogin 中可能会用到它的 this
let autoLogin = singletonPromise(actions.silentLogin).bind(actions)

const loginUrl = '/login'

// 白名单：不需要登录也能访问的页面
const whiteList = [
    '/',
    '/signup',
    loginUrl,
    '/certificates'
]

// 自动登录
autoLogin()

router.beforeEach(async (to, from, next) => {
    
    // 判断是否已经登录
    const isLogin = getters.isLogin()
    if (isLogin) {
        if (to.path === loginUrl) {
            // 不允许重复进入登录页面
            next({ path: '/' })
        } else {
            next()
        }
    }
    // 未登录
    else {
        if (whiteList.indexOf(to.path) !== -1) {
            // 无需登录即可访问
            next()
        } else {
            // 尝试根据之前登录成功的保留的信息自动（静默）登录
            autoLogin()
                .then(() => {
                    next()          // 登录成功正常跳转
                })
                .catch(() => {      // 登录失败就重定向去登录
                    next({
                        path: loginUrl,
                        query: { redirect: to.path }
                    })
                })
        }
    }
})