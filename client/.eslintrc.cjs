# ESLint configuration for JavaScript/React
# https://eslint.org/docs/user-guide/configuring

module.exports = {
  root: true,
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:react/recommended',
    'plugin:react/jsx-runtime',
    'plugin:react-hooks/recommended',
  ],
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
    ecmaFeatures: {
      jsx: true,
    },
  },
  settings: {
    react: {
      version: 'detect',
    },
  },
  plugins: ['react', 'react-hooks', 'react-refresh'],
  rules: {
    // React rules
    'react/prop-types': 'off', // Turn off if using TypeScript
    'react/react-in-jsx-scope': 'off', // Not needed in React 17+
    'react/jsx-uses-react': 'off',
    'react/jsx-uses-vars': 'error',
    'react/jsx-key': 'error',
    'react/no-array-index-key': 'warn',
    'react/jsx-no-duplicate-props': 'error',
    'react/jsx-no-undef': 'error',
    'react/no-children-prop': 'warn',
    'react/no-unescaped-entities': 'warn',
    
    // React Hooks rules
    'react-hooks/rules-of-hooks': 'error',
    'react-hooks/exhaustive-deps': 'warn',
    
    // React Refresh
    'react-refresh/only-export-components': 'warn',
    
    // General JavaScript rules
    'no-console': ['warn', { allow: ['warn', 'error'] }],
    'no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
    'no-undef': 'error',
    'no-duplicate-imports': 'error',
    'prefer-const': 'warn',
    'no-var': 'error',
    'eqeqeq': ['error', 'always'],
    'curly': ['error', 'all'],
    'brace-style': ['error', '1tbs'],
    'comma-dangle': ['error', 'always-multiline'],
    'quotes': ['error', 'single', { avoidEscape: true }],
    'semi': ['error', 'always'],
    'indent': ['error', 2, { SwitchCase: 1 }],
    'max-len': ['warn', { code: 100, ignoreStrings: true, ignoreTemplateLiterals: true }],
    'arrow-parens': ['error', 'always'],
    'arrow-body-style': ['warn', 'as-needed'],
  },
  ignorePatterns: [
    'dist',
    'build',
    'node_modules',
    '*.config.js',
    'vite.config.js',
  ],
};
