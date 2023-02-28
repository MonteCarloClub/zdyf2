declare namespace API {

    /** 用户信息参数 */
    type UserInfoParams = {
        [key: string] : string;
    };

    /** 用户信息 */
    type UserInfoResponse = {
        priv: string;
        pub: string;
    };
}