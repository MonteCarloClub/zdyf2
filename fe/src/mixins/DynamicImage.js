// 支持动态资源路径
export const DynamicImage = function (dir = "") {
    return {
        methods: {
            imgHome(name) {
                // https://stackoverflow.com/questions/40491506/vue-js-dynamic-images-not-working/47480286
                // expression inside v-bind is executed in runtime, webpack aliases work in compile time.
                // https://github.com/vuejs/vue-loader/issues/896
                return require("../assets/images/" + dir + name);
            },
        }
    }
}