Page({
  data: {
    dishes: [
      { id: 1, name: '清蒸鲈鱼', time: 20, cal: 180, bg: 'var(--sky-bg)' },
      { id: 2, name: '蒜蓉西兰花', time: 12, cal: 85, bg: 'var(--mint-bg)' },
      { id: 3, name: '紫薯米饭', time: 30, cal: 200, bg: 'var(--lavender-bg)' }
    ],
    ingredients: ['鲈鱼', '西兰花', '紫薯', '大米', '蒜', '葱', '姜', '盐', '生抽']
  },

  addToShopping() {
    wx.showToast({ title: '已加入购物清单', icon: 'success' })
    setTimeout(() => wx.redirectTo({ url: '/pages/shopping/index' }), 1000)
  },

  saveMenu() {
    wx.showToast({ title: '菜单已保存', icon: 'success' })
  }
})
