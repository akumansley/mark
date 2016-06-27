var HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    entry: './client/main.js',
    output: {
        path: 'server/data/static/build',
        filename: 'bundle.js'
    },
    module: {
        loaders: [{
            test: /\.js$/,
            exclude: /node_modules/,
            loaders: ['react-hot', 'babel?presets[]=es2015&presets[]=react'],
        }]
    },
    resolve: {
        extensions: ['', '.js', '.json', '.coffee']
    },
    plugins: [new HtmlWebpackPlugin({
        title: "Mark",
        template: './client/index.html',
    })],
    devServer: {
        proxy: {
            "/api/*": {
                "target": {
                    "host": "localhost",
                    "port": 8081,
                }
            }
        }
    },
};
