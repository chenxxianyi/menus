const login = () => {
  return new Promise((resolve, reject) => {
    wx.login({
      success: res => {
        if (res.code) {
          const { post } = require('./request')
          post('/auth/wechat-login', { code: res.code }).then(data => {
            wx.setStorageSync('token', data.token)
            wx.setStorageSync('userInfo', data.user)
            resolve(data)
          }).catch(reject)
        } else {
          reject(new Error('wx.login failed'))
        }
      },
      fail: reject
    })
  })
}

const getToken = () => wx.getStorageSync('token') || ''
const getUserInfo = () => wx.getStorageSync('userInfo') || null
const isLoggedIn = () => !!getToken()

const requireLogin = (page) => {
  if (!isLoggedIn()) {
    wx.navigateTo({ url: '/pages/user/index' })
    return false
  }
  return true
}

module.exports = { login, getToken, getUserInfo, isLoggedIn, requireLogin }
