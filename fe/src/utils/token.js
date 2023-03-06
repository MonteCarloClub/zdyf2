import * as Cookies from "./cookies"

const TOKEN_KEY = 'Admin-Token'

export function getToken() {
    let token = Cookies.get(TOKEN_KEY)
    if (token) {
        let _ = token.split("-")
        if (_.length === 2) {
            return {
                name: _[0],
                password: _[1]
            }
        }
        return false
    }
    return false
}

export function setToken(token) {
    return Cookies.set(TOKEN_KEY, token, { 
        secure: true, 
        "max-age": 3600 // 登录有效期，单位（秒）
    })
}

export function removeToken() {
    return Cookies.remove(TOKEN_KEY)
}