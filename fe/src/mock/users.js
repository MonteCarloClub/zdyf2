function mockUsers() {

    function getUser(name) {
        const userStr = localStorage.getItem(name)
        return userStr ? JSON.parse(userStr) : false
    }

    function setUser(user) {
        localStorage.setItem(user.name, JSON.stringify(user))
    }

    function successResp(data) {
        return {
            err: 0,
            data
        }
    }

    function errorResp(msg) {
        return {
            err: 1,
            msg
        }
    }

    return {
        signup(user) {
            if (getUser(user.name)) {
                return errorResp("用户已存在")
            }

            setUser(user)
            return successResp(user)
        },

        login(_user) {
            const { fileName, password } = _user
            const user = getUser(fileName);

            if (!user) {
                return errorResp("用户不存在")
            }

            // 登录成功，返回用户信息
            if (user.password === password) {
                user.password += "-hash"
                return successResp({
                    ...user,
                    "appliedAttrMap": {
                        "someone:family": "[159429..., 572246...]",
                        "someone:friend": "[114119..., 477210...]",
                    },
                    "privacyAttrMap": {},
                    "APKMap": {},
                    "ASKMap": {},
                    "EGGAlpha": "[206605, 320061]",
                    "Alpha": "907358",
                    "GAlpha": "[569334, 105875]",
                    "OPKMap": {},
                    "OSKMap": {},
                })
            }
            // 密码错误
            return errorResp("密码错误")
        },

        logout(token) {
            return successResp({ token })
        }
    }
}

export const localUsers = mockUsers()
