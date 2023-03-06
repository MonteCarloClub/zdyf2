// 提供过滤空值的功能
export const FilterEmpty = {
    methods: {
        filterEmpty(arr) {
            if (arr == undefined) {
                return ""
            }
            if (Array.isArray(arr)) {
                let res = arr.join(" ").split(" ").filter(s => s)
                return res
            }
            return arr
        },
    }
}