const path = require('path')

module.exports = {
    entry: path.resolve(__dirname, './src/scripts/index.js'),
    output: {
        path: path.resolve(__dirname, './src/scripts/build'),
        filename: 'build.js'
    },
    watch: true
}