const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const { ESBuildMinifyPlugin } = require("esbuild-loader");
const path = require('path');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const { optimizeImports } = require("carbon-preprocess-svelte");

const mode = process.env.NODE_ENV || 'development';
const prod = mode === 'production';

module.exports = {
	entry: {
		'bundle': ['./src/main.js']
	},
	resolve: {
		alias: {
			svelte: path.dirname(require.resolve('svelte/package.json'))
		},
		extensions: ['.mjs', '.js', '.svelte'],
		mainFields: ['svelte', 'browser', 'module', 'main']
	},
	output: {
		path: path.join(__dirname, '/dist'),
		filename: '[name].js',
		chunkFilename: '[name].[id].js'
	},
	module: {
		rules: [
			{
				test: /\.svelte$/,
				use: {
					loader: "svelte-loader",
					options: {
						preprocess: [optimizeImports()],
						hotReload: !prod,
						compilerOptions: { dev: !prod },
					},
				},
			},
			{
				test: /\.css$/,
				use: [MiniCssExtractPlugin.loader, "css-loader"],
			},
			{
				test: /node_modules\/svelte\/.*\.mjs$/,
				resolve: { fullySpecified: false },
			},
		]
	},
	mode,
	plugins: [
		new MiniCssExtractPlugin({
			filename: '[name].css'
		}),
		new CopyWebpackPlugin({
			patterns: [
				{ from: 'public' }
			]
		})
	],
	devtool: 'source-map',
	devServer: {
		hot: true,
		proxy: {	
			'/api': 'https://yourinstance.com/api'
		}
	},
	optimization: {
		minimizer: [new ESBuildMinifyPlugin({ target: "es2015" })],
	},
};
