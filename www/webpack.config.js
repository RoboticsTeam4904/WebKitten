module.exports = {
    entry: "./index.js",
    output: {
        path: __dirname + '/assets',
        filename: 'bundle.js',
        publicPath: '/assets'
    },
    devtool: 'source-map',
    module: {
    loaders: [
      {
        test: /\.js$/,
        exclude: /(node_modules|bower_components)/,
        loader: 'babel', // 'babel-loader' is also a valid name to reference
        query: {
          presets: ['es2015', 'react'],
          sourceMaps: true
        }
      },
      {
        test: /\.css$/,
        loader: 'style-loader!css-loader'
      },
    ]
  }
}