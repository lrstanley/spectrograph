/** @type {import('vls').VeturConfig} */
module.exports = {
    settings: {
        // "vetur.useWorkspaceDependencies": true,
        // "vetur.experimental.templateInterpolationService": true
        "vetur.format.defaultFormatterOptions": {
            "prettier": {
                "printWidth": 160,
                "tabWidth": 4,
                "trailingComma": "es5",
                "arrowParens": "always",
                "quoteProps": "as-needed",
                "singleQuote": false,
                "useTabs": false,
                "semi": false,
                "bracketSpacing": true
            }
        },
        "vetur.format.options.tabSize": 4,
        "vetur.validation.templateProps": true
    },
    projects: [
        {
            root: './cmd/frontend/',
            package: 'package.json',
            globalComponents: [
                './src/components/**/*.vue'
            ]
        }
    ]
}
