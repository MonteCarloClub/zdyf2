
export const localFiles = {

    same: function (_) {
        return _
    },

    files: function () {
        return {
            "contents": [
                {
                    "fileName": "test.txt",
                    "policy": "(someone:friend AND someone:family)",
                    "cipher": "xxx",
                    "tags": [
                        "shanghai",
                        "myc",
                        "edu",
                        "test"
                    ],
                    "sharedUser": "someone"
                }
            ],
            "bookmark": "g1AAAA...",
            "pageSize": 10,
            "count": 1
        }
    }
}
