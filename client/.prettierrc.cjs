module.exports = {
  "root": true,
  "plugins": [
    import("@ianvs/prettier-plugin-sort-imports"),
    import("prettier-plugin-tailwindcss")
  ],
  "printWidth": 80,
  "tabWidth": 2,
  "singleQuote": false,
  "trailingComma": "all",
  "arrowParens": "always",
  "semi": false,
  "endOfLine": "auto",
  "tailwindFunctions": ["clsx"],
}