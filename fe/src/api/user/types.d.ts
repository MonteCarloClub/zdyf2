declare namespace API {

    /** 用户信息参数 */
    type DemoParams = {
        [key: string] : string;
    };

    /** 用户索引参数 */
    type UserParams = {
        uid : string;
    };

    /** 黑名单 */
    type BlacklistResponse = {
        certificates: string[];
    };

    /** 用户移出黑名单 */
    type rmBlacklistResponse = string;

    /** 用户加入黑名单 */
    type addBlacklistResponse = string;
}