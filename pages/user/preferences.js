const { tastes, goals } = require('../../utils/constants')

Page({
  data: {
    tastes,
    goals,
    taste: '',
    goal: '',
    avoid: [],
    avoidOptions: ['香菜', '葱', '蒜', '辣椒', '海鲜', '乳糖', '花生', '芹菜'],
    familySize: 3
  },

  onLoad() {
    const pref = wx.getStorageSync('preferences')
    if (pref) {
      this.setData({
        taste: pref.taste || '',
        goal: pref.goal || '',
        avoid: pref.avoid || [],
        familySize: pref.familySize || 3
      })
    }
  },

  setTaste(e) { this.setData({ taste: e.currentTarget.dataset.val }) },
  setGoal(e) { this.setData({ goal: e.currentTarget.dataset.val }) },

  toggleAvoid(e) {
    const val = e.currentTarget.dataset.val
    const avoid = this.data.avoid.includes(val)
      ? this.data.avoid.filter(i => i !== val)
      : [...this.data.avoid, val]
    this.setData({ avoid })
  },

  changeFamily(e) {
    const delta = parseInt(e.currentTarget.dataset.delta)
    this.setData({ familySize: Math.max(1, Math.min(20, this.data.familySize + delta)) })
  },

  save() {
    const { taste, goal, avoid, familySize } = this.data
    wx.setStorageSync('preferences', { taste, goal, avoid, familySize })
    wx.showToast({ title: '已保存', icon: 'success' })
    setTimeout(() => wx.navigateBack(), 800)
  }
})
