module.exports = {
    "transpileDependencies": [
        "vuetify"
    ],
    devServer: {
        proxy: {
            '/': {
                target: 'http://localhost:18080/',
                changeOrigin: true,
                secure: false,
                ws: false,
                pathRewrite: {
                    '^/': '',
                }
            }
        },
        disableHostCheck: true
    },
}
