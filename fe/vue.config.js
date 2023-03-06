module.exports = (options = {}) => ({
  publicPath: './',
  // outputDir: '../backend/src/main/resources/static',
  devServer: {
    host: 'localhost',
    port: 8081,
    disableHostCheck: true,   // That solved it
    proxy: {
      '/api/': {
        target: 'http://10.176.40.46:8080',
        changeOrigin: true,
        pathRewrite: {
          '^/api': ''
        }
      },
      '/cert/': {
        // target: 'http://10.176.40.46/dpki/',
        target: 'http://10.176.40.48/dpki',
        changeOrigin: true,
        pathRewrite: {
          '^/cert': ''
        }
      }
    },
  },
})
