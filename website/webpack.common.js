var path = require('path');

module.exports = {
    entry: {
        consumer: './src/ConsumerWeb/index.tsx',
        ops: './src/AdminWeb/index.tsx',
    },
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: '[name]Bundle.js',
    },
    resolve: {
        extensions: [ '.tsx', '.ts', '.js' ],
        modules: [
            path.resolve(__dirname, 'node_modules'),
            path.resolve(__dirname, './src'),
        ]
    },
    module: {
        rules : [
            {test : /\.(js)$/, use:'babel-loader'},
            {test: /\.tsx?$/, use: 'ts-loader', exclude: /node_modules/},
            {test : /\.s?css$/, use:['style-loader', 'css-loader', 'sass-loader'], exclude: /node_modules/},
            {test: /\.json$/, use: 'json-loader'},
        ],
    },
};
