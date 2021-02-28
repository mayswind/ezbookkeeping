const GitRevisionPlugin = require('git-revision-webpack-plugin');
const MomentLocalesPlugin = require('moment-locales-webpack-plugin');

const pkgFile = require('./package.json');
const licenseFile = require('./third-patry-licenses.json');

module.exports = {
    pages: {
        desktop: {
            entry: 'src/desktop-main.js',
            template: 'src/public/desktop.html',
            filename: 'desktop.html',
            chunks: ['vendors-common-bundle', 'vendors-desktop-bundle', 'common-bundle', 'desktop']
        },
        mobile: {
            entry: 'src/mobile-main.js',
            template: 'src/public/mobile.html',
            filename: 'mobile.html',
            chunks: ['vendors-common-bundle', 'vendors-mobile-bundle', 'common-bundle', 'mobile']
        }
    },
    publicPath: '',
    productionSourceMap: false,
    configureWebpack: {
        plugins: [
            new MomentLocalesPlugin()
        ]
    },
    chainWebpack: config => {
        config.optimization.splitChunks({
            cacheGroups: {
                'vendors-common-bundle': {
                    name: 'vendors-common-bundle',
                    test: /[\\/]node_modules[\\/]/,
                    chunks: 'initial',
                    priority: 10,
                    minChunks: 2
                },
                'vendors-bundle': {
                    name: (module, chunks) => {
                        const allChunksNames = chunks.map((item) => item.name).join('-');
                        return `vendors-${allChunksNames}-bundle`;
                    },
                    test: /[\\/]node_modules[\\/]/,
                    chunks: 'initial',
                    priority: 5,
                    minChunks: 1
                },
                'common-bundle': {
                    name: 'common-bundle',
                    chunks: 'initial',
                    priority: 1,
                    minChunks: 2
                }
            }
        });

        config.plugin('define').tap(definitions => {
            const gitRevisionPlugin = new GitRevisionPlugin();
            definitions[0]['process.env']['VERSION'] = JSON.stringify(pkgFile.version);
            definitions[0]['process.env']['COMMIT_HASH'] = JSON.stringify(gitRevisionPlugin.commithash());
            definitions[0]['process.env']['BUILD_UNIXTIME'] = JSON.stringify(parseInt((new Date().getTime() / 1000).toString()));
            definitions[0]['process.env']['LICENSES'] = JSON.stringify(licenseFile);

            return definitions;
        });
    },
    devServer: {
        host: '0.0.0.0',
        port: 8081,
        disableHostCheck: true,
        proxy: {
            '/api': {
                target: 'http://127.0.0.1:8080/',
                changeOrigin: true
            }
        }
    }
}
