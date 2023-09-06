<template>
  <div class="wrapper">
    <div class="content" v-html="result"></div>
  </div>
</template>

<script setup lang="ts">
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'
import { getProto } from '~/api/user';

const route = useRoute()
const { protocolName } = route.params

const post = ref('')
const result = ref()

const getMarkdown = () => {
  getProto(protocolName as string).then((res: any) => {
    post.value = res.content
    console.log(post.value)
    result.value = md.render(post.value)
  })
}
getMarkdown()

let md: any = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,        // Convert '\n' in paragraphs into <br>
  highlight: function (str:any, lang:any) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre><code class="language-${lang} hljs">` +
               hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
               '</code></pre>';
      } catch (__) {}
    }

    return '<pre><code class="language-none hljs">' + md.utils.escapeHtml(str) + '</code></pre>';
  }
})


definePageMeta({
  layout: false
})
</script>
<style lang="scss">
  .content {
    max-width: 1000px;
    margin: 0 auto
  }
</style>