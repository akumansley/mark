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
            loader: 'babel',
            query: { presets: ['es2015', 'react'] }
        }]
    },
    resolve: {
        extensions: ['', '.js', '.json', '.coffee']
    },
    plugins: [new HtmlWebpackPlugin({
      title: "Mark",
      template: './client/index.html',
    })],
};
