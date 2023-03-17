/**
 * @param {string} phone
 * @returns {Boolean}
 */
 export function validPhone(phone:string):boolean {
  const reg = /^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$/
  return reg.test(phone)
}

/**
 * @param {string} password
 * @returns {Boolean}
 */
export function validPassword(password:string):boolean {
  const reg = /^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,16}$/
  return reg.test(password)
}


/** 手机号正则 */
export const mobileReg = /^[1][3,4,5,6,7,8,9][0-9]{9}$/