Page({
  data: {
    userInfo: {},
    isLoggedIn: false
  },

  onShow() {
    const userInfo = wx.getStorageSync('userInfo')
    this.setData({
      userInfo: userInfo || {},
      isLoggedIn: !!userInfo
    })
  },

  goPage(e) {
    wx.navigateTo({ url: e.currentTarget.dataset.url })
  },

  doLogin() {
    wx.getUserProfile({
      desc: '用于完善用户资料',
      success: (res) => {
        const userInfo = {
          name: res.userInfo.nickName,
          avatar: res.userInfo.avatarUrl,
          phone: ''
        }
        wx.setStorageSync('userInfo', userInfo)
        wx.setStorageSync('token', 'mock-token')
        this.setData({ userInfo, isLoggedIn: true })
        wx.showToast({ title: '登录成功', icon: 'success' })
      },
      fail: () => {
        // Mock login fallback
        const userInfo = { name: '美食家', phone: '138****1234' }
        wx.setStorageSync('userInfo', userInfo)
        wx.setStorageSync('token', 'mock-token')
        this.setData({ userInfo, isLoggedIn: true })
        wx.showToast({ title: '登录成功', icon: 'success' })
      }
    })
  }
})
