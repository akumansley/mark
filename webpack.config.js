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
        },
        {
           test: /\.(png|jpg)$/, 
           loader: 'url-loader?limit=8192'  // inline base64 URLs for <=8k images, direct URLs for the rest
        }
      ]
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
            },
            "/assets/*": {
                "target": {
                    "host": "localhost",
                    "port": 8081,
                }
            }
        }
    },
};
