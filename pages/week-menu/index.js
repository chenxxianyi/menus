Page({
  data: {
    days: [
      { key: 'mon', label: '周一', date: '5/12' },
      { key: 'tue', label: '周二', date: '5/13' },
      { key: 'wed', label: '周三', date: '5/14' },
      { key: 'thu', label: '周四', date: '5/15' },
      { key: 'fri', label: '周五', date: '5/16' },
      { key: 'sat', label: '周六', date: '5/17' },
      { key: 'sun', label: '周日', date: '5/18' }
    ],
    activeDay: 0,
    currentMeals: {
      '早餐': [
        { id: 1, name: '鸡蛋饼', time: 10 },
        { id: 2, name: '牛奶燕麦', time: 5 }
      ],
      '午餐': [
        { id: 3, name: '番茄炒蛋', time: 15 },
        { id: 4, name: '紫薯米饭', time: 30 }
      ],
      '晚餐': [
        { id: 5, name: '清蒸鱼', time: 20 },
        { id: 6, name: '蒜蓉西兰花', time: 12 }
      ]
    }
  },

  setDay(e) {
    this.setData({ activeDay: e.currentTarget.dataset.idx })
    // TODO: 加载对应天的菜单
  },

  regenerate(e) {
    wx.showToast({ title: '已重新生成', icon: 'success' })
  },

  addAllToShopping() {
    wx.showToast({ title: '已加入购物清单', icon: 'success' })
    setTimeout(() => wx.navigateTo({ url: '/pages/shopping/index' }), 800)
  }
})
