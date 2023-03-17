// https://nuxt.com/docs/api/configuration/nuxt-config

import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'
import path from 'path'

import { loadEnv } from 'vite'

const envScript = (process.env as any).npm_lifecycle_script.split(' ')
const envName = envScript[envScript.length - 1] // 通过启动命令区分环境
const envData = loadEnv(envName, 'env')

export default defineNuxtConfig({
  // ssr: process.env.NODE_ENV !== "development",
  runtimeConfig: {
    public: envData
  },
  css: ['element-plus/dist/index.css',"~/assets/css/main.scss"],
  vite: {
    envDir: '~/env', // 指定env文件夹
    server: {
      proxy: {
        '/accounts': {
          target: 'http://10.1.0.57:8086',  //这里是接口地址
          changeOrigin: true
        },
        '/v1': {
          target: 'http://10.1.0.212:8085',
          changeOrigin: true,
        }
      }
    },
    
    plugins: [
      createSvgIconsPlugin({
          iconDirs: [path.resolve(process.cwd(), 'assets/svg')]
      })
  ],
  },
})
