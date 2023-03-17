/**
 * 获取url参数值
 * @param url url地址
 */
export function useGetQuery(url: string){
  // 通过 ? 分割获取后面的参数字符串
  let urlStr = url.split('?')[1]
  // 创建空对象存储参数
  let obj: any = {};
    // 再通过 & 将每一个参数单独分割出来
  let paramsArr = urlStr.split('&')
  for(let i = 0,len = paramsArr.length;i < len;i++){
        // 再通过 = 将每一个参数分割为 key:value 的形式
    let arr = paramsArr[i].split('=')
    obj[arr[0]] = arr[1];
  }
  return obj
}

/**
 * 深度合并对象
 */
export function mergeObj(from: any, to: any) {
  for (let key in from) {
    if (String(from[key]) === '[object Object]' && to[key]) {
      mergeObj(from[key], to[key])
    } else {
      to[key] = from[key]
    }
  }
  return to
}