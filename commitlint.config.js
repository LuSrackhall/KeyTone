module.exports = {
  extends: ["@commitlint/config-conventional"],
  rules: {
    // 提交时的主内容限制默认为100, 但我是个啰嗦的人, 往往100是不够用的, 因此将其扩展到150。
    "header-max-length": [2, "always", 150],
    // "header-max-length": [1, "always", 100], // 无法同时定义警告和错误<不过可以通过自定义插件来完成这个需求, 但没必要>
    // 这个配置的含义如下：

    // 2: 规则的严重性级别。

    // 0：表示关闭这个规则。
    // 1：表示这个规则的警告级别。
    // 2：表示这个规则的错误级别。
    // 'always': 应用规则的条件。

    // 'always'：表示这个规则应该总是被应用。
    // 'never'：表示这个规则应该永远不会被应用。
    // 150: 规则的具体值，这里是指 commit message header 的最大长度为 150 个字符。

    // 所以，'header-max-length': [2, 'always', 150] 的意思是：

    // 如果 commit message header 的长度超过 150 个字符，Commitlint 会抛出一个错误（严重性级别为 2）。
    // 这个规则总是会被应用（条件为 'always'）。
    // 通过这种方式，你可以详细地控制 Commitlint 的行为，根据你的需求定制不同的规则和严重性级别。
  },
};
