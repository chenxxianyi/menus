App({
  onLaunch() {
    this.login()
  },

  login() {
    wx.login({
      success: res => {
        if (res.code) {
          // TODO: 发送 code 到后端换取 token
          console.log('登录 code:', res.code)
        }
      }
    })
  },

  globalData: {
    userInfo: null,
    token: null,
    baseUrl: 'http://localhost:3000/api'
  }
})
