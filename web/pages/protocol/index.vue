<script lang="ts" setup name="Users">
import type { TabsPaneContext } from "element-plus";
import { ref } from "vue";
import Style from "./style/index.vue";
import Energy from "./energy/index.vue";
import { getProto } from "~/api/energy";
import MarkdownIt from "markdown-it";
import hljs from "highlight.js";
const word = ref([]);
const result = ref();
const getInfo = () => {
  getProto().then((res: any) => {
    word.value = res.filter((item: any) => {
      return item.privacyWrite == history.state.url;
    });
    result.value = md.render(word.value && word.value[0] && word.value[0].privacyWord);
  });
};
getInfo();
let md: any = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true, // Convert '\n' in paragraphs into <br>
  highlight: function (str: any, lang: any) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return (
          `<pre><code class="language-${lang} hljs">` +
          hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
          "</code></pre>"
        );
      } catch (__) {}
    }

    return (
      '<pre><code class="language-none hljs">' +
      md.utils.escapeHtml(str) +
      "</code></pre>"
    );
  },
});
definePageMeta({
  layout: false,
});
</script>
<template>
  <div class="content" v-html="result"></div>
</template>
<style scoped lang="scss"></style>
