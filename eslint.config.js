import js from "@eslint/js";
import globals from "globals";

export default [
  {
    ignores: ["node_modules/", "eslint.config.js"],
  },
  js.configs.recommended,
  {
    files: ["**/*.js"],
    languageOptions: {
      ecmaVersion: 2020,
      sourceType: "script",
      globals: {
        ...globals.browser,
        firebase: "readonly",
      },
    },
    rules: {
      "no-unused-vars": ["error", { argsIgnorePattern: "^_" }],
      "no-var": "off",
      eqeqeq: ["error", "always"],
      curly: ["error", "multi-line"],
    },
  },
];
